//go:build ignore

package main

import "fmt"

func main() {
	a, b, c := 3, -1, 5
	loads := []*int{&a, &b, nil, &c}
	indices := []int{}

	for i, load := range loads {
		// nil means None; deref of nil panics, so check first.
		if load == nil {
			fmt.Println("nil at:", i)
			continue
		}

		if *load == -1 {
			indices = append(indices, i)
		}
	}
	fmt.Println("indices with -1:", indices)
}
