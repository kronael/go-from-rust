//go:build ignore

package main

import "fmt"

func main() {
	// Rust for consumes IntoIterator; Go range accepts specific operand types.
	letters := []string{"a", "b", "c"}
	for i, letter := range letters {
		fmt.Printf("slice %d=%s\n", i, letter)
	}

	// Map iteration order is unspecified.
	counts := map[string]int{"x": 1, "y": 2}
	for key, count := range counts {
		fmt.Printf("map %s=%d\n", key, count)
	}

	// A string yields byte indexes and runes.
	for i, r := range "hé" {
		fmt.Printf("string byte %d: %U\n", i, r)
	}

	for i := range 3 {
		fmt.Println("integer:", i)
	}

	// The range value is a copy; assign through the index to change the slice.
	nums := []int{1, 2, 3}
	for _, n := range nums {
		n *= 10
	}
	fmt.Println("value copies:", nums)
}
