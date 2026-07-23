//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// append reuses spare capacity, so aliases stay linked.
	reused := make([]int, 2, 4)
	copy(reused, []int{1, 2})
	alias := reused
	reused = append(reused, 3)
	reused[0] = 9
	fmt.Println("reused:", alias, reused)

	// Capping the view forces the next append to detach.
	backing := []int{1, 2, 3, 4}
	limited := backing[:2:2]
	detached := append(limited, 9)
	detached[0] = 7
	fmt.Println("capacity limited:", backing, detached)

	// Grow reserves room for 3 more; length unchanged.
	reserved := []int{1, 2}
	reserved = slices.Grow(reserved, 3)
	fmt.Println("Grow keeps length:", len(reserved),
		"room:", cap(reserved)-len(reserved) >= 3)
}
