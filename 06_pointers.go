//go:build ignore

package main

import "fmt"

func main() {
	// Passing s copies its slice descriptor, not its elements. Both slice values
	// still describe shared storage, so setFirst changes the first element.
	s := []int{10, 20, 30}
	setFirst(s)
	fmt.Println("slice argument:", s)

	// A *[]int lets replace assign a new slice value back to the caller, roughly
	// like &mut Vec<_>. An ordinary []int already permits element writes; return
	// a new slice instead when only the caller's length or storage must change.
	replace(&s)
	fmt.Println("replaced header:", s)

	// Go lets an array pointer use p[i] as shorthand for (*p)[i].
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
