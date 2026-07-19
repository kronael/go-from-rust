//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	// Spare capacity lets append reuse storage, so existing slices still alias.
	reused := make([]int, 2, 4)
	copy(reused, []int{1, 2})
	alias := reused
	reused = append(reused, 3)
	reused[0] = 9
	fmt.Println("reused:", alias, reused)

	// A full slice expression can cap a view and force the next append to detach.
	backing := []int{1, 2, 3, 4}
	limited := backing[:2:2]
	detached := append(limited, 9)
	detached[0] = 7
	fmt.Println("capacity limited:", backing, detached)

	// Grow ensures room for three more elements without changing the length.
	// It may allocate, so keep the returned slice just as you do with append.
	reserved := []int{1, 2}
	reserved = slices.Grow(reserved, 3)
	fmt.Println("Grow keeps length:", len(reserved), "room:", cap(reserved)-len(reserved) >= 3)

	// Rust prevents a live slice borrow across a Vec mutation that may move it.
	// Go does not specify append's capacity growth factor; never depend on it.
}
