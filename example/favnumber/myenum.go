package favnumber

import "github.com/rekby/genum"

// Declare private type guarantee about the enum values (but zero value) can be created only from the package
//nolint:unused
type p genum.BaseType

type FavoriteNumber = genum.EnumValue[p] // define as alias - for save methods

var (
	Holder, holderPrivate = genum.NewHolders[p]()

	// Zero is default value of FavoriteNumber because internal value is 0
	// Enum value with 0 internal state - only value, which can be created outsize of the package
	// you SHOULD declare zero value for prevent usage undefined enum value
	Zero = holderPrivate.New(0, "zero")
	One  = holderPrivate.New(1, "one")
	Five = holderPrivate.New(5, "five") // It is ok for sparse range with spaces in enum values
)
