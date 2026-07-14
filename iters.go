//go:build ignore

// iters.go — range-over-function iterators introduced in Go 1.23.
// Rust analogue: impl Iterator; but Go uses a push model — the iterator calls yield,
// the consumer does not call next(). yield returning false signals an early stop (break).
// Key takeaway: iter.Seq[V] is just func(yield func(V) bool); no trait, no state struct.

package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func main() {
	// === built-in ranges (always existed) ===
	for i, v := range []string{"a", "b"} { // slice: index, value
		fmt.Printf("slice %d=%s\n", i, v)
	}
	for i := range 3 { // integer range (Go 1.22)
		fmt.Println("int range", i)
	}

	// === range-over-function: consume a custom iterator (Go 1.23) ===
	for v := range count(3) {
		fmt.Println("count", v)
	}
	for i, name := range enumerate([]string{"x", "y"}) {
		fmt.Printf("enum %d=%s\n", i, name)
	}

	// break works — the iterator stops cleanly
	for v := range count(100) {
		if v == 2 {
			break
		}
		fmt.Println("until break", v)
	}

	// === stdlib functions that RETURN iterators ===
	m := map[string]int{"go": 1, "rust": 2}
	keys := slices.Sorted(maps.Keys(m)) // iterator -> sorted slice
	fmt.Println("sorted keys:", keys)
}

// single-value iterator: iter.Seq[V] == func(yield func(V) bool)
func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) { // yield returns false when the consumer breaks
				return
			}
		}
	}
}

// two-value iterator: iter.Seq2[K,V] == func(yield func(K, V) bool)
func enumerate[T any](s []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}
