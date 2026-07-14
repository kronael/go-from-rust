//go:build ignore

// popping elements: end/front/index ordered/swap-remove.
// Rust analogue: Vec::pop() -> Option<T>; Go has no built-in — you reslice manually.
// Gotcha: q = q[1:] (pop front) advances the pointer, pinning the backing array forever.
// Takeaway: swap-remove is O(1) and idiomatic when order doesn't matter.

package main

import (
	"fmt"
	"slices"
)

func main() {
	// pop end (stack / LIFO) — O(1); Rust: vec.pop()
	s := []int{1, 2, 3, 4}
	x := s[len(s)-1]
	s = s[:len(s)-1]
	fmt.Printf("pop end: got %d, left %v\n", x, s)

	// pop front (queue / FIFO) — O(1) but leaks: the front of the array is never freed
	q := []int{1, 2, 3, 4}
	f := q[0]
	q = q[1:]
	fmt.Printf("pop front: got %d, left %v\n", f, q)

	// pop at index, preserve order — O(n); Rust: vec.remove(i)
	a := []int{10, 20, 30, 40}
	i := 1
	v := a[i]
	a = slices.Delete(a, i, i+1)
	fmt.Printf("pop idx %d: got %d, left %v\n", i, v, a)

	// swap-remove: move last element into the hole — O(1); Rust: vec.swap_remove(i)
	b := []int{10, 20, 30, 40}
	j := 1
	w := b[j]
	b[j] = b[len(b)-1]
	b = b[:len(b)-1]
	fmt.Printf("swap-remove idx %d: got %d, left %v\n", j, w, b)

	// guard empty: Go panics on out-of-bounds, no Option return
	var empty []int
	if len(empty) > 0 {
		_ = empty[len(empty)-1]
	} else {
		fmt.Println("empty: nothing to pop")
	}
}
