//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	big := make([]int, 1000)

	// Front and end subslices both keep the backing allocation
	// reachable.
	front := big[:3]
	end := big[997:]
	fmt.Printf("front: len=%d cap=%d\n",
		len(front), cap(front))
	fmt.Printf("end: len=%d cap=%d\n", len(end), cap(end))

	// Clip reduces capacity without copying, so it still
	// retains and aliases big.
	front[0] = 1
	clipped := slices.Clip(front)
	clipped[0] = 2
	fmt.Println("Clip aliases:", big[0] == 2,
		"cap:", cap(clipped))

	// Like Rust to_vec, Clone copies surviving values away
	// from the large array.
	clone := slices.Clone(clipped)

	// Set every alias to nil, or clone the values that must
	// survive, to make the large array eligible for garbage
	// collection. Reclamation timing is not fixed.
	big, front, end, clipped = nil, nil, nil, nil
	fmt.Println("clone survives:", clone)
}
