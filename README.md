Enum emulation based on go generics, introduced from go 1.18.
The package allow compile time gurantee about enum-ed var can't contains value outside from enum definition.
No allocations after initialization.

Roadmap:
- genum linter for check exhaustive switch 

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
import "github.com/rekby/genum/example_favnumber"

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

## Benchmarks
```
go test -test.bench=. -test.benchmem ./benchmarks
goos: darwin
goarch: arm64
pkg: github.com/rekby/genum/benchmarks
BenchmarkSwitchIntEnum-10               335044386                3.416 ns/op           0 B/op          0 allocs/op
BenchmarkSwitchGEnum-10                 192531540                6.227 ns/op           0 B/op          0 allocs/op
BenchmarkStringIntWithStringer-10       460761220                2.640 ns/op           0 B/op          0 allocs/op
BenchmarkStringGEnum-10                 81390418                13.84 ns/op            0 B/op          0 allocs/op
BenchmarkToStringGenumWithHolder-10     333978289                3.567 ns/op           0 B/op          0 allocs/op
BenchmarkFromIntInt-10                  1000000000               0.3111 ns/op          0 B/op          0 allocs/op
BenchmarkFromIntGEnum-10                428099434                2.828 ns/op           0 B/op          0 allocs/op
BenchmarkFromStringGEnum-10             425914281                2.766 ns/op           0 B/op          0 allocs/op
```