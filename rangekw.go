//go:build ignore

// rangekw.go — the `range` keyword over slice, map, string, int, and func.
// Rust uses .iter()/.iter_mut(); Go uses a single `range` keyword for all.
// Yield type varies by operand; two key gotchas: map order is random,
// and ranging over a string yields (byte-index, rune) — NOT bytes.
// Takeaway: always range by index when mutating slice elements.

package main

import "fmt"

func main() {
	s := []string{"a", "b", "c"}
	for i, v := range s { // slice -> (index, value)
		fmt.Printf("slice %d=%s\n", i, v)
	}

	m := map[string]int{"x": 1, "y": 2}
	for k, v := range m { // map -> (key, value), order is RANDOM every run
		fmt.Printf("map %s=%d\n", k, v)
	}

	for i, r := range "héllo" { // string -> (byte index, rune), not bytes
		fmt.Printf("str byte%d=%c (%d)\n", i, r, r)
	}

	for i := range 3 { // integer (Go 1.22+) -> 0..n-1, like Rust 0..3
		fmt.Println("int", i)
	}

	for range s { // discard both values — loop len(s) times
	}

	// GOTCHA: v is a copy — mutating it does NOT touch the slice.
	nums := []int{1, 2, 3}
	for _, v := range nums {
		v *= 10 // dead write; nums unchanged
	}
	fmt.Println("after copy-modify:", nums)

	for i := range nums {
		nums[i] *= 10 // index into the slice to mutate
	}
	fmt.Println("after index-modify:", nums)
}
