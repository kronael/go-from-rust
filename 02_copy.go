//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	source := []int{1, 2, 3}

	// Clone: new backing array, shallow element copy.
	clone := slices.Clone(source)
	source[0] = 9
	fmt.Println("Clone creates:", source, clone)

	// copy fills min(len) elements, returns that count.
	destination := make([]int, 2)
	fmt.Println("copy fills:",
		copy(destination, source), destination)

	// copy also handles overlapping slices.
	overlap := []int{1, 2, 3, 4}
	copy(overlap[1:], overlap)
	fmt.Println("copy overlap:", overlap)
}
