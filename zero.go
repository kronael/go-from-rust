//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	a, b, c := 10, 20, 30

	// Rust Vec removal/drop makes ownership explicit; Go's GC follows references.
	// Delete clears pointer slots, so the backing array no longer retains them.
	deleted := []*int{&a, &b, &c}
	deleted = slices.Delete(deleted, 2, 3)
	deletedFull := deleted[:cap(deleted)]
	fmt.Println("Delete values:", *deleted[0], *deleted[1])
	fmt.Println("Delete tail nil:", deletedFull[2] == nil)

	// Plain reslicing leaves the pointer reachable through the backing array.
	resliced := []*int{&a, &b, &c}
	resliced = resliced[:2]
	reslicedFull := resliced[:cap(resliced)]
	fmt.Println("reslice tail nil:", reslicedFull[2] == nil)
	fmt.Println("reslice tail value:", *reslicedFull[2])
}
