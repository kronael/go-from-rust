//go:build ignore

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

	// SortFunc is unstable, like Rust's sort_unstable_by.
	// Its comparator must define a strict weak ordering.
	slices.SortFunc(people, func(a, b Person) int {
		// Or returns the first nonzero result, but evaluates every argument.
		return cmp.Or(
			cmp.Compare(a.Age, b.Age),
			cmp.Compare(a.Name, b.Name),
		)
	})

	fmt.Println(people)
}
