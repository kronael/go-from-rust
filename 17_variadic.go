//go:build ignore

package main

import "fmt"

func sum(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func main() {
	// Inside sum, numbers is []int; args build the slice.
	fmt.Println("individual:", sum(1, 2, 3))

	// slice... passes an existing slice, no copy made.
	numbers := []int{4, 5, 6}
	fmt.Println("slice:", sum(numbers...))
}
