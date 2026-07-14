//go:build ignore

package main

import "fmt"

func main() {
	src := []int{1, 2, 3}

	// Rust copy_from_slice requires equal lengths.
	// Go copy uses the smaller length.
	dst := make([]int, 3)
	fmt.Println("copy equal:", copy(dst, src), dst)

	small := make([]int, 2)
	fmt.Println("copy smaller:", copy(small, src), small)

	// copy uses len, not cap, so this destination receives nothing.
	empty := make([]int, 0, 10)
	fmt.Println("copy empty:", copy(empty, src), empty)

	// Go variadics have no direct Rust parameter syntax; Rust takes a slice.
	fmt.Println("variadic:", sum(1, 2, 3))
	nums := []int{4, 5, 6}
	fmt.Println("spread slice:", sum(nums...))
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
