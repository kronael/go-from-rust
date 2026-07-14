//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	big := make([]int, 1000)

	// Front and end subslices both keep the backing allocation reachable.
	end := big[:3]
	front := big[997:]
	fmt.Printf("end: len=%d cap=%d\n", len(end), cap(end))
	fmt.Printf("front: len=%d cap=%d\n", len(front), cap(front))

	// Like Rust to_vec, Clone copies into independent storage.
	end[0] = 1
	clone := slices.Clone(end)
	clone[0] = 9
	fmt.Println("original:", end)
	fmt.Println("clone:", clone)

	// The large allocation is collectible only after all its slices are
	// unreachable.
	fmt.Println("independent:", big[0] != clone[0])
}
