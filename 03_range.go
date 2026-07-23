//go:build ignore

package main

import "fmt"

func main() {
	// Rust for consumes IntoIterator; Go range accepts
	// specific operand types.
	letters := []string{"a", "b", "c"}
	for i, letter := range letters {
		fmt.Printf("slice %d=%s\n", i, letter)
	}

	// A string yields byte indexes and runes: Unicode code
	// points. %U prints the code point as U+XXXX, so the
	// second index is 1 even though é uses two bytes.
	for i, r := range "hé" {
		fmt.Printf("string byte %d: %U\n", i, r)
	}

	// Ranging over 3 yields 0, 1, 2, like Rust's 0..3.
	for i := range 3 {
		fmt.Println("integer:", i)
	}

	// The range value is a copy; assign through the index to
	// change the slice.
	nums := []int{1, 2, 3}
	for _, n := range nums {
		n *= 10
	}
	fmt.Println("value copies:", nums)
}
