//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// Go has no pop function: guard, read the end, zero it, then reslice.
	// clear writes 0 here; for pointer elements it writes nil.
	stack := []int{1, 2, 3, 4}
	if len(stack) > 0 {
		index := len(stack) - 1
		last := stack[index]
		clear(stack[index:])
		stack = stack[:index]
		fmt.Println("pop back:", last, stack)
	}

	// slices.Delete preserves order, like Rust Vec::remove.
	ordered := []int{10, 20, 30, 40}
	i := 1
	removed := ordered[i]
	ordered = slices.Delete(ordered, i, i+1)
	fmt.Println("delete:", removed, ordered)

	// slices.Insert preserves order, like Rust Vec::insert.
	ordered = slices.Insert(ordered, i, 25)
	fmt.Println("insert:", ordered)

	// Go has no swap-delete function; move the last value and zero its old slot.
	// Zeroing keeps removed pointer-like values from being retained by the array.
	unordered := []int{10, 20, 30, 40}
	i = 1
	removed = unordered[i]
	unordered[i] = unordered[len(unordered)-1]
	clear(unordered[len(unordered)-1:])
	unordered = unordered[:len(unordered)-1]
	fmt.Println("swap-delete:", removed, unordered)
}
