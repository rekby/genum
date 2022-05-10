package genumenum

import "github.com/rekby/genum"

// Declare private type guarantee about the enum values (but zero value) can be created only from the package
//nolint:unused
type p genum.BaseType

type Variant = genum.EnumValue[p] // define as alias - for save methods

var (
	Holder, holderPrivate = genum.NewHolders[p]()

	A = holderPrivate.New(0, "A")
	B = holderPrivate.New(1, "B")
	C = holderPrivate.New(2, "C") // It is ok for sparse range with spaces in enum values
)

func FromIntUnsafe(val int) Variant {
	return holderPrivate.UnsafeFromInt(val)
}
