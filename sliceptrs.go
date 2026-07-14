//go:build ignore

package main

import "fmt"

func main() {
	a, b, c := 3, -1, 5
	loads := []*int{&a, &b, nil, &c}
	indices := []int{}

	for i, load := range loads {
		// Dereferencing nil panics; Rust Option<&T> is normally matched first.
		if load == nil {
			fmt.Println("nil at:", i)
			continue
		}

		// Elements of []*int require dereferencing to read the int.
		if *load == -1 {
			indices = append(indices, i)
		}
	}
	fmt.Println("indices with -1:", indices)
}
