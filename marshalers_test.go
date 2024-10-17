package genum

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaryMarshaler(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		type p BaseType
		type valueType = EnumValue[p]
		_, holder := NewHolders[p]()

		zero := holder.New(0, "zero")

		otherVal := 123
		other := holder.New(otherVal, "other")

		// zero
		var readValue valueType
		dataBytes, err := zero.MarshalBinary()
		require.NoError(t, err)
		err = readValue.UnmarshalBinary(dataBytes)
		require.NoError(t, err)
		require.Equal(t, zero, readValue)

		// other
		dataBytes, err = other.MarshalBinary()
		require.NoError(t, err)
		err = readValue.UnmarshalBinary(dataBytes)
		require.NoError(t, err)
		require.Equal(t, other, readValue)
	})
	t.Run("UnmarshalError", func(t *testing.T) {
		type p BaseType
		type varType = EnumValue[p]
		_, holder := NewHolders[p]()
		_ = holder.New(2, "two")

		var val varType
		t.Run("Empty", func(t *testing.T) {
			require.Error(t, val.UnmarshalBinary(nil))
		})
		t.Run("Unexpected", func(t *testing.T) {
			buf := make([]byte, 1)
			binary.PutVarint(buf, 3)
			require.Error(t, val.UnmarshalBinary(buf))
		})
	})
}

func TestTextMarshaler(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		type p BaseType
		type valueType = EnumValue[p]
		_, holder := NewHolders[p]()

		zero := holder.New(0, "zero")

		otherVal := 123
		other := holder.New(otherVal, "other")

		// zero
		var readValue valueType
		dataBytes, err := zero.MarshalText()
		require.NoError(t, err)
		require.EqualValues(t, []byte("zero"), dataBytes)

		err = readValue.UnmarshalText(dataBytes)
		require.NoError(t, err)
		require.Equal(t, zero, readValue)

		// other
		dataBytes, err = other.MarshalText()
		require.NoError(t, err)
		require.EqualValues(t, []byte("other"), dataBytes)

		err = readValue.UnmarshalText(dataBytes)
		require.NoError(t, err)
		require.Equal(t, other, readValue)
	})
	t.Run("UnmarshalError", func(t *testing.T) {
		type p BaseType
		type varType = EnumValue[p]
		_, holder := NewHolders[p]()
		_ = holder.New(2, "two")

		var val varType
		t.Run("Empty", func(t *testing.T) {
			require.Error(t, val.UnmarshalText(nil))
		})
		t.Run("Unexpected", func(t *testing.T) {
			buf := make([]byte, 1)
			binary.PutVarint(buf, 3)
			require.Error(t, val.UnmarshalText(buf))
		})
	})
}

func FuzzEnumValueMarshalBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, name string, intVal int) {
		type p BaseType
		type _ = EnumValue[p]

		_, holder := NewHolders[p]()
		defer deleteHolder[p]()

		val := holder.New(intVal, name)
		data, err := val.MarshalBinary()
		require.NoError(t, err)
		readVar, readLen := binary.Varint(data)
		require.Equal(t, int64(intVal), readVar)
		require.Equal(t, len(data), readLen)
	})
}

func FuzzEnumValueUnmarshalBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, valExist bool, data []byte) {
		if len(data) > 11 {
			t.Skip("datalen")
		}

		rawVal, readLen := binary.Varint(data)
		badDataLen := readLen != len(data)
		t.Log("readval:", rawVal, readLen, badDataLen)

		if rawVal > math.MaxInt {
			t.Skip()
		}

		type p BaseType
		type ValueType = EnumValue[p]

		_, holder := NewHolders[p]()
		defer func() {
			deleteHolder[p]()
		}()

		if valExist {
			_ = holder.New(int(rawVal), "val")
		}

		var val ValueType
		err := val.UnmarshalBinary(data)
		if !valExist || badDataLen || readLen <= 0 {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, rawVal, int64(val.val))
		}
	})
}

func FuzzEnumValueMarshalUnmarshalBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, rawVal int64) {
		if rawVal > math.MaxInt {
			t.Skip()
		}

		type p BaseType
		type valType = EnumValue[p]

		_, holder := NewHolders[p]()
		defer func() {
			deleteHolder[p]()
		}()
		val := holder.New(int(rawVal), "val")
		data, err := val.MarshalBinary()
		require.NoError(t, err)

		var val2 valType
		err = val2.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Equal(t, rawVal, int64(val2.Int()))
	})
}
