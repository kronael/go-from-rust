//go:build ignore

package main

import "fmt"

func main() {
	// Passing a slice copies the header; elements are shared.
	s := []int{10, 20, 30}
	setFirst(s)
	fmt.Println("slice argument:", s)

	// *[]int lets a callee replace the caller's slice value.
	replace(&s)
	fmt.Println("replaced header:", s)

	// p[i] on an array pointer is shorthand for (*p)[i].
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
