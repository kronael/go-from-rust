//go:build ignore

// caseA.go — DELIBERATELY does not compile.
// Demonstrates: "cannot indirect v (variable of type int)" — you cannot
// dereference a plain int. Rust makes this impossible by type; Go catches it
// at compile time.
// Compare with caseB.go where the slice holds *int and *v is valid.

package main

import "fmt"

func main() {
	loads := []int{3, -1, 5, -1}
	vis := []int{}
	for i, v := range loads {
		if *v == -1 { // v is int, not *int — compile error: cannot indirect
			vis = append(vis, i)
		}
	}
	fmt.Println(vis)
}
