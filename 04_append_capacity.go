//go:build ignore

package main

import "fmt"

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

	// Rust prevents a live slice borrow across a Vec mutation that may move it.
	// Go does not specify append's capacity growth factor; never depend on it.
}
