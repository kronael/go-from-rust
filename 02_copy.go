//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	source := []int{1, 2, 3}

	// Like Rust's to_vec, Clone makes an independent shallow copy.
	clone := slices.Clone(source)
	source[0] = 9
	fmt.Println("clone:", source, clone)

	// copy reuses a destination and copies the smaller of the two lengths.
	destination := make([]int, 2)
	fmt.Println("copy:", copy(destination, source), destination)

	// ... expands a slice into append's variadic arguments. This also copies,
	// but Clone states the intent directly.
	appendCopy := append([]int(nil), source...)
	fmt.Println("append copy:", appendCopy)

	// copy supports overlap, like Rust's copy_within.
	overlap := []int{1, 2, 3, 4}
	copy(overlap[1:], overlap)
	fmt.Println("copy overlap:", overlap)
}
