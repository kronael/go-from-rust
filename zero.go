//go:build ignore

// slices.Delete zeroes freed tail slots so pointer elements can be GC'd.
// Rust analogue: Vec::remove shifts elements and drops the removed value immediately.
// Gotcha: manual reslice (s = s[:n-1]) does NOT zero — the backing array retains the pointer,
// preventing GC. slices.Delete handles this correctly; hand-rolled reslices do not.

package main

import (
	"fmt"
	"slices"
)

func main() {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}

	s = slices.Delete(s, 2, 3) // len now 2; zeroes slot[2] in the backing array

	full := s[:cap(s)]
	fmt.Printf("len=%d cap=%d, slot[2] after Delete = %v\n", len(s), cap(s), full[2])

	// manual reslice: faster, but leaves the pointer in the backing array — GC leak
	t := []*int{&a, &b, &c}
	t = t[:len(t)-1]
	tfull := t[:cap(t)]
	fmt.Printf("manual reslice, slot[2] = %v (still points at value)\n", tfull[2])
}
