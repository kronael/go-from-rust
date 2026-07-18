//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	source := []int{1, 2, 3}

	// Clone creates and returns full-length independent slice storage,
	// like Rust's to_vec.
	clone := slices.Clone(source)
	source[0] = 9
	fmt.Println("Clone creates:", source, clone)

	// copy fills destination storage you provide. It does not grow destination;
	// it copies min(len(destination), len(source)) and returns that count.
	destination := make([]int, 2)
	fmt.Println("copy fills:", copy(destination, source), destination)

	// Both operations copy element values, not objects those values point to.

	// ... expands a slice into append's variadic arguments. This also copies,
	// but Clone states the intent directly.
	appendCopy := append([]int(nil), source...)
	fmt.Println("append copy:", appendCopy)

	// copy supports overlap, like Rust's copy_within.
	overlap := []int{1, 2, 3, 4}
	copy(overlap[1:], overlap)
	fmt.Println("copy overlap:", overlap)
}
