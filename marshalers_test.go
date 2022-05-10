package genum

import (
	"encoding/binary"
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

func FuzzEnumValueMarshalBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, name string, intVal int) {
		type p BaseType
		type valType = EnumValue[p]

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
