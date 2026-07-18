//go:build ignore

package main

import "fmt"

func main() {
	// A []T argument normally suffices: it still refers to the same elements.
	s := []int{10, 20, 30}
	setFirst(s)
	fmt.Println("slice argument:", s)

	// *[]T is uncommon. Use it when replacing the caller's slice header.
	replace(&s)
	fmt.Println("replaced header:", s)

	// Array pointer indexing is special shorthand for (*p)[i] in the Go spec.
	array := [3]int{1, 2, 3}
	arrayPointer := &array
	arrayPointer[0] = 9
	fmt.Println("array pointer:", array)
}

func setFirst(s []int) {
	s[0] = 11
}

func replace(s *[]int) {
	*s = []int{7, 8}
}
