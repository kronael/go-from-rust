//go:build ignore

// filtering slices: new-slice loop, in-place s[:0] trick, slices.DeleteFunc.
// Rust analogue: iter().filter().collect() always allocates; s[:0] reuses backing storage.
// Gotcha: slices.DeleteFunc predicate is "drop if true" — opposite of filter's "keep if true".
// Takeaway: s[:0] is zero-allocation filter-in-place; use DeleteFunc for clean one-liners.

package main

import (
	"fmt"
	"slices"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6}
	keep := func(n int) bool { return n%2 == 0 }

	// 1. new slice — Rust: nums.iter().filter(|&&n| n%2==0).collect()
	var out []int
	for _, n := range nums {
		if keep(n) {
			out = append(out, n)
		}
	}
	fmt.Println("new slice:", out)

	// 2. in-place: s[:0] resets len to 0 but keeps the backing array; zero allocation
	src := []int{1, 2, 3, 4, 5, 6}
	filtered := src[:0]
	for _, n := range src {
		if keep(n) {
			filtered = append(filtered, n)
		}
	}
	fmt.Println("in place:", filtered)

	// 3. slices.DeleteFunc: predicate means DROP (not keep) — easy to get backwards
	xs := []int{1, 2, 3, 4, 5, 6}
	xs = slices.DeleteFunc(xs, func(n int) bool { return n%2 == 0 })
	fmt.Println("DeleteFunc (drops evens):", xs)
}
