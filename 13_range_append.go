//go:build ignore

package main

import "fmt"

func main() {
	edges := [][2]int{{1, 2}, {3, 4}}

	// range fixes its count before the loop starts.
	for _, edge := range edges {
		edges = append(edges, [2]int{edge[1], edge[0]})
	}

	fmt.Println("final:", edges)
}
