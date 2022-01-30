Enum emulation based on go generics, introduced from go 1.18.

Usage example:

Declare enum values
```go
package favnumber

import "github.com/rekby/genum"

// Declare private type guarantee about the enum values (but zero value) can be created only from the package
type p genum.BaseType

type FavoriteNumber = genum.EnumValue[p] // define as alias - for save methods

var (
	Holder = genum.NewHolder[p]()

	// Zero is default value of FavoriteNumber because internal value is 0
	// Enum value with 0 internal state - only value, which can be created outsize of the package
	// you SHOULD declare zero value for prevent usage undefined enum value
	Zero = Holder.New(0, "zero")
	One  = Holder.New(1, "one")
	Five = Holder.New(5, "five") // It is ok for sparse range with spaces in enum values
)
```

Usage enum values:
```go
    package example

    import "fmt"
    import "github.com/rekby/genum/example/favnumber"

    func ExampleUsage() {
		one := favnumber.One
		five, _ := favnumber.Holder.FromString("five")

		var min favnumber.FavoriteNumber
		if one.Int() < five.Int() {
			min = one
		} else {
			min = five
		}
		fmt.Printf("min: %v\n", min)

		isFirst := false
		parsed, _ := favnumber.Holder.FromInt(1)
		switch parsed {
		case favnumber.One:
			isFirst = true
		default:
			isFirst = false
		}

		fmt.Printf("isFirst: %v", isFirst)
		// Output:
		// min: one
		// isFirst: true
	}
```