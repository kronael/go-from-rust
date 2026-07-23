//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	big := make([]int, 1000)

	// A subslice keeps the whole backing array reachable.
	front := big[:3]
	end := big[997:]
	fmt.Printf("front: len=%d cap=%d\n",
		len(front), cap(front))
	fmt.Printf("end: len=%d cap=%d\n", len(end), cap(end))

	// Clip caps capacity without copying; still aliases big.
	front[0] = 1
	clipped := slices.Clip(front)
	clipped[0] = 2
	fmt.Println("Clip aliases:", big[0] == 2,
		"cap:", cap(clipped))

	// Clone copies survivors into independent storage.
	clone := slices.Clone(clipped)

	// Nil every alias to free the large array.
	big, front, end, clipped = nil, nil, nil, nil
	fmt.Println("clone survives:", clone)
}
