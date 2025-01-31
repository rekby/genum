[![Go Reference](https://pkg.go.dev/badge/github.com/rekby/genum.svg)](https://pkg.go.dev/github.com/rekby/genum)
[![Coverage Status](https://coveralls.io/repos/github/rekby/genum/badge.svg?branch=master)](https://coveralls.io/github/rekby/genum?branch=master)

Enum emulation based on go generics, introduced from go 1.18.

In comparison with int-based enums:
1. Has comparable performance (and workaround for same performance when need).
2. Has method for parse value from string
3. Compile-time gurantee about no unexpected values
4. Without code generation

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
	Holder, holderPrivate = genum.NewHolders[p]()

	// Zero is default value of FavoriteNumber because internal value is 0
	// Enum value with 0 internal state - only value, which can be created outsize of the package
	// you SHOULD declare zero value for prevent usage undefined enum value
	Zero = holderPrivate.New(0, "zero")
	One  = holderPrivate.New(1, "one")
	Five = holderPrivate.New(5, "five") // It is ok for sparse range with spaces in enum values
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

## Performance
GEnum values have comparable perfomance with enums, based on raw int values.

### Switch
genum has same performance in comparison operations (==, switch) as int-based enum: genum internally stored as one int value.
on compare stage types already checked by compiler and runtime has no additional instructions.

### .String()
genum has two way of convert enum-values to string: direct .String() method and call ValueToString on holder object.
Direct call is 10x slower, then .String() of method, generated by stringer: genum internal stored int value only and it need 
additional search of values dictionary by type object.
If .String() must be faster in hot path - it has workaround: call ValueToString on holder object has only one map search by int key.
It performance similar to generated .String() method.

### Create enum value from int
Base variant of genum has 10x slower then int-based enum. Int based enums has no any checks while create new value and
new value can contains any int-value, not only enumerated as constants. Genum check about raw int value defined as
valid enum value.

Usually create value from int need during deserialize input values. In code base values can created by copy of predefined 
values and will same fast at int based.

If you need high-performance unchecked enum-values creations - private holder has unsafe method for create enum values.
It create enum value with raw value without any check (as int-based) with same performance. Your code must gurantee about
input values are valid.

## Benchmarks
```
go test -test.bench=. -test.benchmem ./benchmarks
goos: darwin
goarch: arm64
pkg: github.com/rekby/genum/benchmarks
BenchmarkSwitchIntEnum-10              	312518718	         3.735  ns/op	       0 B/op	       0 allocs/op
BenchmarkSwitchGEnum-10                	349887081	         3.420  ns/op	       0 B/op	       0 allocs/op
BenchmarkStringIntWithStringer-10      	434715668	         2.649  ns/op	       0 B/op	       0 allocs/op
BenchmarkStringGEnum-10                	82807876	        13.93   ns/op	       0 B/op	       0 allocs/op
BenchmarkToStringGenumWithHolder-10    	337688142	         3.563  ns/op	       0 B/op	       0 allocs/op
BenchmarkFromIntInt-10                 	1000000000	         0.3118 ns/op	       0 B/op	       0 allocs/op
BenchmarkFromIntGEnum-10               	426436011	         2.814  ns/op	       0 B/op	       0 allocs/op
BenchmarkFromIntGEnumUnsafe-10         	1000000000	         0.3117 ns/op	       0 B/op	       0 allocs/op
BenchmarkFromStringGEnum-10            	435763668	         2.753  ns/op	       0 B/op	       0 allocs/op
```