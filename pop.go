//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// Rust Vec::pop returns Option; Go must guard before indexing.
	stack := []int{1, 2, 3, 4}
	if len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Println("pop back:", last, stack)
	}

	// Rust Vec::remove and slices.Delete preserve order in O(n-i).
	ordered := []int{10, 20, 30, 40}
	i := 1
	removed := ordered[i]
	ordered = slices.Delete(ordered, i, i+1)
	fmt.Println("ordered remove:", removed, ordered)

	// Rust Vec::swap_remove and this swap-delete take O(1).
	unordered := []int{10, 20, 30, 40}
	i = 1
	removed = unordered[i]
	unordered[i] = unordered[len(unordered)-1]
	unordered = unordered[:len(unordered)-1]
	fmt.Println("swap-delete:", removed, unordered)

	// Front reslicing is O(1), but the result keeps the allocation reachable.
	queue := []int{1, 2, 3, 4}
	first := queue[0]
	queue = queue[1:]
	fmt.Println("pop front:", first, queue)
}
