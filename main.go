package main

import "fmt"

func main() {
	// Rust arrays copy when their elements are Copy; Go arrays are values.
	array := [3]int{1, 2, 3}
	arrayCopy := array
	arrayCopy[0] = 9
	fmt.Println("array copy:", array, arrayCopy)

	// Go slice assignment copies a header, so both slices share elements.
	slice := []int{1, 2, 3}
	sliceAlias := slice
	sliceAlias[0] = 9
	fmt.Println("slice alias:", slice, sliceAlias)
}
