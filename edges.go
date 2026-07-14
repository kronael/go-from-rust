//go:build ignore

package main

import "fmt"

// appending while ranging: range captures the slice header once at loop start, so it terminates.
// Rust analogue: for x in &v appends nothing to v; the borrow checker prohibits it at compile time.
// Gotcha: C-style for i < len(s) re-evaluates len each iteration and loops forever.
func main() {
	// range version: safe, terminates
	edges := [][2]int{{1, 2}, {3, 4}}
	for _, v := range edges {
		edges = append(edges, [2]int{v[1], v[0]})
	}
	fmt.Println("range result:", edges, "(len", len(edges), ")")

	// C-style len()-in-condition version: would never end
	e2 := [][2]int{{1, 2}, {3, 4}}
	n := 0
	for i := 0; i < len(e2); i++ {
		e2 = append(e2, [2]int{e2[i][1], e2[i][0]})
		if n++; n > 6 {
			fmt.Println("C-style: still growing after", n, "iters -> INFINITE")
			break
		}
	}
}
