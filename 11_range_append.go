//go:build ignore

package main

import "fmt"

func main() {
	edges := [][2]int{{1, 2}, {3, 4}}

	// range fixes the original iteration length even when append grows the slice.
	// Rust rejects push during iteration through a shared borrow.
	for _, edge := range edges {
		edges = append(edges, [2]int{edge[1], edge[0]})
	}

	fmt.Println("final:", edges)
}
