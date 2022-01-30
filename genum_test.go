package genum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newHolders() map[any]any {
	return make(map[any]any)
}

func TestNewHolder(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		type p int

		holders := newHolders()
		h1 := newHolder[p](holders)
		require.NotNil(t, h1)
		require.Len(t, holders, 1)

		var zero p
		require.Equal(t, h1, holders[zero])
	})

	t.Run("HolderForInt", func(t *testing.T) {
		holders := newHolders()
		require.Panics(t, func() {
			newHolder[int](holders)
		})
	})

	t.Run("HolderForDefaultType", func(t *testing.T) {
		holders := newHolders()
		require.Panics(t, func() {
			newHolder(holders)
		})
	})

	t.Run("TwoHoldersForSameType", func(t *testing.T) {
		type p int
		holders := newHolders()
		newHolder[p](holders)
		require.Panics(t, func() {
			newHolder[p](holders)
		})
	})

	t.Run("HolderForPublicType", func(t *testing.T) {
		type P int
		holders := newHolders()
		require.Panics(t, func() {
			newHolder[P](holders)
		})
	})
}

func TestEnumHolder_All(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		type p int
		type enumType = EnumValue[p]
		holders := newHolders()
		holder := newHolder[p](holders)
		v2 := holder.New(2, "2")
		v1 := holder.New(1, "1")
		require.Equal(t, []enumType{v1, v2}, holder.All())
	})
	t.Run("empty", func(t *testing.T) {
		type p int
		holders := newHolders()
		holder := newHolder[p](holders)
		require.Len(t, holder.All(), 0)
	})
}

func TestEnumHolder_FromInt(t *testing.T) {
	type p int
	holder := NewHolder[p]()
	vMinus := holder.New(-1, "minus")
	v0 := holder.New(0, "zero")
	v1 := holder.New(1, "one")

	t.Run("ok", func(t *testing.T) {
		parsed, err := holder.FromInt(-1)
		require.NoError(t, err)
		require.Equal(t, parsed, vMinus)

		parsed, err = holder.FromInt(0)
		require.NoError(t, err)
		require.Equal(t, parsed, v0)

		parsed, err = holder.FromInt(1)
		require.NoError(t, err)
		require.Equal(t, parsed, v1)
	})

	t.Run("err", func(t *testing.T) {
		parsed, err := holder.FromInt(2)
		require.Error(t, err)
		require.Equal(t, v0, parsed)
	})
}

func TestEnumHolder_FromString(t *testing.T) {
	type p int
	holder := NewHolder[p]()
	vMinus := holder.New(-1, "minus")
	v0 := holder.New(0, "zero")
	v1 := holder.New(1, "one")

	t.Run("ok", func(t *testing.T) {
		parsed, err := holder.FromString("minus")
		require.NoError(t, err)
		require.Equal(t, vMinus, parsed)

		parsed, err = holder.FromString("zero")
		require.NoError(t, err)
		require.Equal(t, v0, parsed)

		parsed, err = holder.FromString("one")
		require.NoError(t, err)
		require.Equal(t, v1, parsed)
	})

	t.Run("error", func(t *testing.T) {
		parsed, err := holder.FromString("asd")
		require.Error(t, err)
		require.Equal(t, v0, parsed)
	})
}

func TestEnumHolder_New(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		type p int
		h := NewHolder[p]()
		h.New(0, "zero")
		h.New(1, "one")
		require.Equal(t, map[string]int{"zero": 0, "one": 1}, h.stringToInt)
		require.Equal(t, map[int]string{0: "zero", 1: "one"}, h.intToString)
	})

	t.Run("DoubleInt", func(t *testing.T) {
		type p int
		h := NewHolder[p]()
		h.New(0, "zero")
		require.Panics(t, func() {
			h.New(0, "one")
		})
	})

	t.Run("DoubleString", func(t *testing.T) {
		type p int
		h := NewHolder[p]()
		h.New(0, "zero")
		require.Panics(t, func() {
			h.New(1, "zero")
		})
	})
}

func TestEnumValue_Int(t *testing.T) {
	type p int

	holder := NewHolder[p]()
	v0 := holder.New(0, "zero")
	v5 := holder.New(5, "five")
	v2 := holder.New(2, "two")

	require.Equal(t, 0, v0.Int())
	require.Equal(t, 2, v2.Int())
	require.Equal(t, 5, v5.Int())
}

func TestEnumValue_String(t *testing.T) {
	type p int

	holder := NewHolder[p]()
	v0 := holder.New(0, "zero")
	v5 := holder.New(5, "five")
	v2 := holder.New(2, "two")

	require.Equal(t, "zero", v0.String())
	require.Equal(t, "two", v2.String())
	require.Equal(t, "five", v5.String())
}

func TestSetPanicOnUnexistedValues(t *testing.T) {
	// check default value
	require.False(t, panicOnEnexistedValues)

	SetPanicOnUnexistedValues(true)
	require.True(t, panicOnEnexistedValues)
	SetPanicOnUnexistedValues(false)
	require.False(t, panicOnEnexistedValues)
}
