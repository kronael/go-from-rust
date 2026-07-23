//go:build ignore

package main

import "fmt"

// A constraint is an interface used as a type bound. Number's type set is
// any type whose underlying type is int or float64; ~ admits named types
// like `type Celsius float64`. Rust writes it as a bound: <T: Add + Copy>.
type Number interface {
	~int | ~float64
}

// sum accepts any type in Number; += compiles because the set allows it.
// The Rust analogue is a bounded fn sum<T: Add>(xs).
func sum[T Number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

// comparable is the built-in constraint for == and !=, what map keys and
// this lookup need. Rust's nearest bound is <T: PartialEq>.
func index[T comparable](values []T, target T) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}

type Celsius float64

func main() {
	// Type arguments are inferred from the values; no turbofish needed.
	fmt.Println("ints:", sum([]int{1, 2, 3}))
	fmt.Println("floats:", sum([]float64{1.5, 2.5}))

	// A named type satisfies Number through ~float64; no separate impl needed.
	fmt.Println("named:", sum([]Celsius{20, 1.5}))

	// comparable admits strings: == is defined, but < is not, so sum would
	// not compile for them. The constraint states which operators are legal.
	fmt.Println("index:", index([]string{"a", "b", "c"}, "b"))
}
