package example

import (
	"fmt"

	"github.com/rekby/genum/example/favnumber"
)

func ExampleDefaultValue() {
	var test favnumber.FavoriteNumber
	fmt.Println(test.Int())
	fmt.Println(test.String())
	// Output:
	// 0
	// zero
}

func ExampleFromInt() {
	val, _ := favnumber.Holder.FromInt(1)
	fmt.Println(val.Int())
	fmt.Println(val.String())
	// Output:
	// 1
	// one
}

func ExampleFromString() {
	val, _ := favnumber.Holder.FromString("five")
	fmt.Println(val.Int())
	fmt.Println(val.String())

	// Output:
	// 5
	// five
}

func ExampleSwitch() {
	one := favnumber.One
	switch one {
	case favnumber.Zero:
		fmt.Println("Found zero")
	case favnumber.One:
		fmt.Println("Found one")
	case favnumber.Five:
		fmt.Println("Found five")
	default:
		fmt.Println("No found")
	}

	// Output:
	// Found one
}

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
