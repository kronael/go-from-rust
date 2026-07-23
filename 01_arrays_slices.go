//go:build !web

package main

import "fmt"

func main() {
	// Rust arrays copy when their elements are Copy; Go arrays are values.
	array := [3]int{1, 2, 3}
	arrayCopy := array
	arrayCopy[0] = 9
	fmt.Println("array copy:", array, arrayCopy)

	// A slice value describes a segment of an underlying array. Think of
	// it as an array pointer plus length and capacity: assignment copies
	// that descriptor, not the elements, so both slices still reach the
	// same array.
	slice := []int{1, 2, 3}
	sliceAlias := slice
	sliceAlias[0] = 9
	fmt.Println("slice alias:", slice, sliceAlias)
}
