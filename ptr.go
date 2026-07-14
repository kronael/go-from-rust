//go:build ignore

// indexing *[]int vs *[N]int: slice pointers require explicit deref; array pointers auto-deref.
// Rust analogue: &mut Vec<T> and &mut [T; N] both index directly — Go's *[]T does not.
// Takeaway: (*p)[i] for *[]T; ap[i] works for *[N]T because the spec auto-derefs array pointers.

package main

import "fmt"

func main() {
	s := []int{10, 20, 30}
	p := &s // *[]int — does NOT auto-deref on indexing

	x := (*p)[0]
	fmt.Println("read (*p)[0]:", x)

	(*p)[1] = 99
	fmt.Println("after (*p)[1]=99:", s)

	fmt.Println("len(*p):", len(*p))

	for i, v := range *p {
		fmt.Printf("  [%d]=%d\n", i, v)
	}

	// *[N]int DOES auto-deref — special rule for array pointers only
	arr := [3]int{1, 2, 3}
	ap := &arr // *[3]int
	ap[0] = 7  // shorthand for (*ap)[0]; not allowed for *[]int
	fmt.Println("array ptr direct index:", arr)
}
