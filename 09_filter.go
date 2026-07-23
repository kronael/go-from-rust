//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// Append-filter builds a new slice (this allocates).
	nums := []int{1, 2, 3, 4, 5, 6}
	var evens []int
	for _, n := range nums {
		if n%2 == 0 {
			evens = append(evens, n)
		}
	}
	fmt.Println("new slice:", evens)

	// DeleteFunc removes matches in place.
	nums = []int{1, 2, 3, 4, 5, 6}
	evens = slices.DeleteFunc(nums,
		func(n int) bool { return n%2 != 0 })
	fmt.Println("DeleteFunc:", evens)
}
