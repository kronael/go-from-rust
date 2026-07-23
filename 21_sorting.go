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
	numbers := []int{4, 1, 3}
	slices.Sort(numbers)
	fmt.Println("natural order:", numbers)

	people := []Person{
		{"Cara", 30},
		{"Ana", 25},
		{"Bob", 25},
		{"Dan", 30},
	}

	// SortFunc supplies custom ordering, like Rust's sort_unstable_by.
	slices.SortFunc(people, func(a, b Person) int {
		// cmp.Or selects the first nonzero comparison: age, then name.
		return cmp.Or(
			cmp.Compare(a.Age, b.Age),
			cmp.Compare(a.Name, b.Name),
		)
	})

	fmt.Println("custom order:", people)

	// Sort and SortFunc are unstable; use SortStableFunc when ties must
	// retain order.
}
