//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// Rust: nums.iter().copied().filter(|n| n % 2 == 0).collect().
	// Collecting into Vec allocates; this loop likewise builds a new slice.
	nums := []int{1, 2, 3, 4, 5, 6}
	var evens []int
	for _, n := range nums {
		if n%2 == 0 {
			evens = append(evens, n)
		}
	}
	fmt.Println("new slice:", evens)

	// DeleteFunc removes elements for which its predicate returns true, in place.
	nums = []int{1, 2, 3, 4, 5, 6}
	evens = slices.DeleteFunc(nums, func(n int) bool { return n%2 != 0 })
	fmt.Println("DeleteFunc:", evens)
}
