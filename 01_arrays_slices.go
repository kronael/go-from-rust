//go:build !web

package main

import "fmt"

func main() {
	// Arrays are values; assignment copies elements.
	array := [3]int{1, 2, 3}
	arrayCopy := array
	arrayCopy[0] = 9
	fmt.Println("array copy:", array, arrayCopy)

	// Slice assignment copies the header, not elements.
	slice := []int{1, 2, 3}
	sliceAlias := slice
	sliceAlias[0] = 9
	fmt.Println("slice alias:", slice, sliceAlias)
}
