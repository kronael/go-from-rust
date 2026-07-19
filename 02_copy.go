//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	source := []int{1, 2, 3}

	// Clone copies the elements into a new backing array. Unlike Rust's to_vec,
	// it does not call Clone on each element: pointer-like elements still share.
	clone := slices.Clone(source)
	source[0] = 9
	fmt.Println("Clone creates:", source, clone)

	// make creates two zero-valued destination elements. copy fills the shorter
	// slice length, returns that count, and never grows the destination.
	destination := make([]int, 2)
	fmt.Println("copy fills:", copy(destination, source), destination)

	// Both operations are shallow. copy also supports overlapping slices, like
	// Rust's copy_within.
	overlap := []int{1, 2, 3, 4}
	copy(overlap[1:], overlap)
	fmt.Println("copy overlap:", overlap)
}
