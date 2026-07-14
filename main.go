// arrays vs slices, value vs reference semantics, sorting.
// Rust analogue: [T; N] is Go's array; Vec<T> is Go's slice.
// Key contrast: slice sub-slicing aliases the backing array (no clone) —
// mutating b = a[1:3] also mutates a. Rust would require explicit &mut or clone.
// Takeaway: slices are (ptr, len, cap) headers; assignment copies the header, not data.

package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
)

type person struct {
	name string
	age  int
}

func main() {
	// array: fixed size, part of the type — Rust: let arr: [i32; 3]
	var arr [3]int
	arr[0] = 10
	fmt.Printf("array %v len=%d (type [3]int)\n", arr, len(arr))

	s := []int{5, 2, 8, 1, 9, 3}

	// slices package (Go 1.21+): preferred
	slices.Sort(s)
	fmt.Println("sorted asc:", s)

	slices.SortFunc(s, func(a, b int) int { return cmp.Compare(b, a) })
	fmt.Println("sorted desc:", s)

	people := []person{{"Cara", 30}, {"Ana", 25}, {"Bob", 25}}
	slices.SortFunc(people, func(a, b person) int {
		return cmp.Or(cmp.Compare(a.age, b.age), cmp.Compare(a.name, b.name))
	})
	fmt.Println("by age then name:", people)

	// sort package still works for common types
	xs := []int{3, 1, 2}
	sort.Ints(xs)
	fmt.Println("sort.Ints:", xs)

	idx, found := slices.BinarySearch([]int{1, 3, 5, 7}, 5)
	fmt.Printf("binary search 5 -> idx=%d found=%v\n", idx, found)

	// aliasing gotcha: b shares the backing array of a — Rust borrow checker prevents this
	a := []int{1, 2, 3, 4}
	b := a[1:3]
	b[0] = 99
	fmt.Println("aliasing: a =", a, " b =", b)
}
