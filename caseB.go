//go:build ignore

// caseB.go — when *v IS valid: a slice of pointers.
// (The real fix for caseA.go is to drop the *, since v was already int.
// This file instead shows the other resolution: make the elements *int.)
// Now v is *int and *v is valid. Watch out: if any pointer is nil, *v panics
// at runtime — Go won't catch it at compile time (unlike Rust's Option).

package main

import "fmt"

func main() {
	a, b, c, d := 3, -1, 5, -1
	loads := []*int{&a, &b, &c, &d}
	vis := []int{}
	for i, v := range loads {
		if *v == -1 { // v is *int — valid; would panic if v == nil
			vis = append(vis, i)
		}
	}
	fmt.Println("indices with -1:", vis)
}
