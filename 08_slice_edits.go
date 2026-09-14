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

	// DeleteFunc removes matches in place, reusing nums.
	nums = []int{1, 2, 3, 4, 5, 6}
	evens = slices.DeleteFunc(nums,
		func(n int) bool { return n%2 != 0 })
	fmt.Println("DeleteFunc:", evens)

	// No pop: guard, read the end, clear it, then reslice.
	stack := []int{1, 2, 3, 4}
	if len(stack) > 0 {
		index := len(stack) - 1
		last := stack[index]
		clear(stack[index:])
		stack = stack[:index]
		fmt.Println("pop back:", last, stack)
	}

	// slices.Delete removes and preserves order.
	ordered := []int{10, 20, 30, 40}
	i := 1
	removed := ordered[i]
	ordered = slices.Delete(ordered, i, i+1)
	fmt.Println("delete:", removed, ordered)

	// slices.Insert adds and preserves order.
	ordered = slices.Insert(ordered, i, 25)
	fmt.Println("insert:", ordered)

	// Swap-delete: move the last value in, then zero it.
	unordered := []int{10, 20, 30, 40}
	i = 1
	removed = unordered[i]
	unordered[i] = unordered[len(unordered)-1]
	clear(unordered[len(unordered)-1:])
	unordered = unordered[:len(unordered)-1]
	fmt.Println("swap-delete:", removed, unordered)
}
