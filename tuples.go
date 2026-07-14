//go:build ignore

// tuples.go — Go has no tuple type.
// Rust: (i32, String) stored, destructured, returned. Go: use a named struct
// for mixed types, [N]T for same-type fixed-size groups.
// Bonus: structs and arrays are comparable, so they work as map keys; slices do not.
// Takeaway: name your structs — anonymous struct literals are a one-off escape hatch.

package main

import "fmt"

// Idiomatic: a named struct is Go's tuple.
type Pair struct {
	Name string
	Age  int
}

// Generic pair — possible, but NOT idiomatic (prefer meaningful field names).
type Tuple[A, B any] struct {
	First  A
	Second B
}

func main() {
	// mixed-type groups -> slice of named structs
	people := []Pair{{"Ana", 25}, {"Bob", 30}}
	fmt.Println("structs:", people)

	// anonymous struct — fine for one-off grouping, but don't reuse
	pts := []struct{ X, Y int }{{1, 2}, {3, 4}}
	fmt.Println("anon:   ", pts)

	// same-type pairs -> [N]T fixed array
	edges := [][2]int{{1, 2}, {3, 4}}
	fmt.Println("arrays: ", edges)

	// [N]T and structs are comparable -> valid map keys; []T is not
	visited := map[[2]int]bool{}
	visited[[2]int{1, 2}] = true
	fmt.Println("map key:", visited[[2]int{1, 2}])

	// generic Tuple (works, rarely seen in the wild)
	t := Tuple[string, int]{"x", 9}
	fmt.Println("generic:", t.First, t.Second)
}
