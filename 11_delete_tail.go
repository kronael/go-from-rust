//go:build ignore

package main

import (
	"fmt"
	"slices"
)

func main() {
	a, b, c := 10, 20, 30

	// Delete clears the emptied pointer slots.
	deleted := []*int{&a, &b, &c}
	deleted = slices.Delete(deleted, 2, 3)
	deletedFull := deleted[:cap(deleted)]
	fmt.Println("Delete values:", *deleted[0], *deleted[1])
	fmt.Println("Delete tail nil:", deletedFull[2] == nil)

	// Reslicing leaves the tail reachable via the array.
	resliced := []*int{&a, &b, &c}
	resliced = resliced[:2]
	reslicedFull := resliced[:cap(resliced)]
	fmt.Println("reslice tail nil:", reslicedFull[2] == nil)
	fmt.Println("reslice tail value:", *reslicedFull[2])

	// clear the tail before reslicing to drop references.
	cleared := []*int{&a, &b, &c}
	clear(cleared[2:])
	cleared = cleared[:2]
	clearedFull := cleared[:cap(cleared)]
	fmt.Println("clear then reslice:", clearedFull[2] == nil)
}
