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

	// Clip reduces capacity without copying, so it still retains and aliases big.
	end[0] = 1
	clipped := slices.Clip(end)
	clipped[0] = 2
	fmt.Println("Clip aliases:", big[0] == 2, "cap:", cap(clipped))

	// Like Rust to_vec, Clone copies values into independent storage.
	clone := slices.Clone(clipped)
	clone[0] = 9
	fmt.Println("Clone independent:", big[0], clone[0])

	// clear removes stored references but keeps an allocation. Set every slice
	// alias to nil, or clone the values that must survive, to release the array.
	big, end, front, clipped = nil, nil, nil, nil
	fmt.Println("clone survives:", clone)
}
