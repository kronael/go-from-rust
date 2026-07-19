//go:build ignore

package main

import "fmt"

func main() {
	a, b, c := 3, -1, 5
	loads := []*int{&a, &b, nil, &c}
	indices := []int{}

	for i, load := range loads {
		// []*int is roughly Vec<Option<&i32>>: nil means None. Check before
		// *load because dereferencing nil panics.
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
