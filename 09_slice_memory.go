//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	a, b, c := 10, 20, 30

	// Reslicing shortens the length only. The tail stays
	// reachable through the backing array, so it is not
	// dropped the way Rust's truncate drops elements.
	resliced := []*int{&a, &b, &c}
	resliced = resliced[:2]
	full := resliced[:cap(resliced)]
	fmt.Println("tail still set:", full[2] != nil, *full[2])

	// clear writes the zero value, nil here, so the
	// reference goes away before the reslice.
	cleared := []*int{&a, &b, &c}
	clear(cleared[2:])
	cleared = cleared[:2]
	fmt.Println("cleared tail nil:",
		cleared[:cap(cleared)][2] == nil)

	// slices.Delete does that clearing for you.
	deleted := []*int{&a, &b, &c}
	deleted = slices.Delete(deleted, 2, 3)
	fmt.Println("Delete tail nil:",
		deleted[:cap(deleted)][2] == nil)

	// Same rule at scale: a 3-element view keeps the
	// whole 1000-element array reachable.
	big := make([]int, 1000)
	front := big[:3]
	fmt.Printf("front: len=%d cap=%d\n",
		len(front), cap(front))

	// Clip caps capacity without copying, so it still
	// aliases big.
	front[0] = 1
	clipped := slices.Clip(front)
	clipped[0] = 2
	fmt.Println("Clip aliases:", big[0] == 2)
	fmt.Println("clipped cap:", cap(clipped))

	// Clone copies into independent storage, so dropping
	// every alias lets the large array be collected.
	clone := slices.Clone(clipped)
	big, front, clipped = nil, nil, nil
	fmt.Println("clone survives:", clone)
}
