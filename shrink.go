//go:build ignore

package main

import "fmt"

// shrinking slices and memory: reslicing never frees the backing array.
// Rust analogue: Vec::truncate is like shrink-end; Vec has no O(1) shrink-from-start,
// but Go's s = s[k:] does exactly that (it just advances the pointer).
// Gotcha: big[997:] still pins the entire 1000-element array via the internal pointer.
// Takeaway: to actually free memory, clone into a right-sized slice.
func main() {
	big := make([]int, 1000) // backing array of 1000

	// shrink end: len drops, backing array unchanged
	end := big[:3]
	fmt.Printf("end   big[:3]   -> len=%d cap=%d  (whole 1000-array still held)\n", len(end), cap(end))

	// shrink start: pointer advances, cap shrinks — but the array is still pinned
	front := big[997:]
	fmt.Printf("front big[997:] -> len=%d cap=%d  (still pins the 1000-array via ptr)\n", len(front), cap(front))

	// clone to actually release: new allocation, old array becomes garbage
	freed := append([]int(nil), big[:3]...) // or slices.Clone(big[:3])
	fmt.Printf("clone           -> len=%d cap=%d  (big array now garbage)\n", len(freed), cap(freed))

	// for pointer slices: zero the tail before truncating so GC can reclaim
	ps := []*int{new(int), new(int), new(int)}
	clear(ps[1:]) // Go 1.21+: zero — Rust drops automatically when element is removed
	ps = ps[:1]
	fmt.Printf("clear+trunc     -> len=%d cap=%d  (tail pointers released)\n", len(ps), cap(ps))
}
