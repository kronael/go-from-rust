//go:build ignore

package main

import (
	"fmt"
	"iter"
	"slices"
)

func count(n int) iter.Seq[int] {
	// Seq[int] is func(yield func(int) bool), not a stateful Iterator object.
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	// range calls the sequence; break makes yield return false so count stops.
	for value := range count(100) {
		fmt.Println("count:", value)
		if value == 2 {
			break
		}
	}

	// slices.All returns the standard index-value Seq2.
	for index, name := range slices.All([]string{"go", "rust"}) {
		fmt.Printf("name %d: %s\n", index, name)
	}

	// Collect consumes a Seq into a slice; iterator adapters compose as functions.
	fmt.Println("collected:", slices.Collect(count(3)))
}
