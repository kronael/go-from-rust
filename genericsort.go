//go:build ignore

// genericsort.go — a single generic Sort via an F-bounded constraint.
// Rust analogue: impl Ord on a type, then slice.sort() — natural order defined once.
// Key takeaway: Comparable[T any] { Compare(T) int } is Go's closest equivalent to
// Rust's Ord trait; the [T Comparable[T]] bound lets T refer to itself (F-bounded).

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

// Define the natural order ONCE, as a method on the type.
func (p Person) Compare(o Person) int {
	return cmp.Or(
		cmp.Compare(p.Age, o.Age),
		cmp.Compare(p.Name, o.Name),
	)
}

// Generic constraint: any type that knows how to compare itself.
type Comparable[T any] interface {
	Compare(T) int
}

// One generic Sort, reusable for EVERY type with a Compare method.
func Sort[T Comparable[T]](s []T) {
	slices.SortFunc(s, func(a, b T) int { return a.Compare(b) })
}

// Works for a totally different type too — just give it Compare.
type Product struct {
	SKU   string
	Price int
}

func (p Product) Compare(o Product) int { return cmp.Compare(p.Price, o.Price) }

func main() {
	people := []Person{{"Cara", 30}, {"Ana", 25}, {"Bob", 25}}
	Sort(people) // no function passed
	fmt.Println("people:", people)

	products := []Product{{"B", 30}, {"A", 10}, {"C", 20}}
	Sort(products) // same Sort, no function passed
	fmt.Println("products:", products)
}
