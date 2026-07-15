//go:build ignore

package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	// Go iterators push values to yield; break makes yield return false.
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

	ages := map[string]int{"go": 15, "rust": 10}
	keys := slices.Sorted(maps.Keys(ages))
	fmt.Println("sorted keys:", keys)
}
