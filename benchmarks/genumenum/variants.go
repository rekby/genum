package genumenum

import "github.com/rekby/genum"

// Declare private type guarantee about the enum values (but zero value) can be created only from the package
type p genum.BaseType

type Variant = genum.EnumValue[p] // define as alias - for save methods

var (
	Holder = genum.NewHolder[p]()

	A = Holder.New(0, "A")
	B = Holder.New(1, "B")
	C = Holder.New(2, "C") // It is ok for sparse range with spaces in enum values
)
