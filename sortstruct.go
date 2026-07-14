//go:build ignore

// sortstruct.go — sorting structs with slices.SortFunc + cmp helpers.
// Rust analogue: slice::sort_by / sort_by_key with a closure returning Ordering.
// Key takeaway: cmp.Or chains comparison keys (short-circuits on first non-zero);
// descending order is just swapping a/b in cmp.Compare — no Reverse wrapper needed.

package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	people := []Person{
		{"Cara", 30},
		{"Ana", 25},
		{"Bob", 25},
		{"Dan", 30},
	}

	// 1. sort by a single field (ascending)
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})
	fmt.Println("by age:        ", people)

	// 2. multi-key: age asc, then name asc as tiebreaker
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Or(
			cmp.Compare(a.Age, b.Age),
			cmp.Compare(a.Name, b.Name),
		)
	})
	fmt.Println("age then name: ", people)

	// 3. descending: swap the arguments
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Compare(b.Age, a.Age)
	})
	fmt.Println("by age desc:   ", people)

	// 4. stable sort: keep equal elements in their original order
	slices.SortStableFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})
	fmt.Println("stable by age: ", people)
}
