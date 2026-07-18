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
	// Rust usually takes &[T]. Go variadics are call-site sugar for []T.
	fmt.Println("individual:", sum(1, 2, 3))

	numbers := []int{4, 5, 6}
	fmt.Println("slice:", sum(numbers...))
}
