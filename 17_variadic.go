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
	// Inside sum, numbers is []int. Separate arguments create
	// its slice.
	fmt.Println("individual:", sum(1, 2, 3))

	// numbers... passes this existing slice as the variadic
	// arguments without a new slice. A []T parameter is
	// usually closest to Rust's &[T].
	numbers := []int{4, 5, 6}
	fmt.Println("slice:", sum(numbers...))
}
