//go:build ignore

package main

import "fmt"

// Constraint interface: the type set ~int | ~float64.
type Number interface {
	~int | ~float64
}

// sum works for any Number; += compiles for the whole set.
func sum[T Number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

// comparable: the built-in constraint for == and !=.
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
	// Type arguments inferred; no turbofish.
	fmt.Println("ints:", sum([]int{1, 2, 3}))
	fmt.Println("floats:", sum([]float64{1.5, 2.5}))

	// A named type satisfies Number via ~float64.
	fmt.Println("named:", sum([]Celsius{20, 1.5}))

	// comparable admits strings; Number excludes them.
	fmt.Println("index:", index([]string{"a", "b", "c"}, "b"))
}
