//go:build ignore

package main

import (
	"cmp"
	"fmt"
	"slices"
)

func sortOrdered[T cmp.Ordered](values []T) {
	slices.Sort(values)
}

type Person struct {
	Name string
	Age  int
}

func main() {
	// cmp.Ordered covers primitive types ordered by <, <=, >=, and >.
	numbers := []int{4, 1, 3}
	names := []string{"rust", "go", "c"}
	sortOrdered(numbers)
	sortOrdered(names)
	fmt.Println(numbers)
	fmt.Println(names)

	// Structs have no operator ordering; pass one to SortFunc.
	people := []Person{{"Cara", 30}, {"Ana", 25}, {"Bob", 28}}
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})
	fmt.Println(people)
}
