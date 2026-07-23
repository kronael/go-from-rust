//go:build ignore

package main

import "fmt"

func main() {
	// range over a slice yields index and copied value.
	letters := []string{"a", "b", "c"}
	for i, letter := range letters {
		fmt.Printf("slice %d=%s\n", i, letter)
	}

	// range over a string yields byte index and rune.
	for i, r := range "hé" {
		fmt.Printf("string byte %d: %U\n", i, r)
	}

	// range over an int n yields 0..n-1.
	for i := range 3 {
		fmt.Println("integer:", i)
	}

	// The range value is a copy; write via the index.
	nums := []int{1, 2, 3}
	for _, n := range nums {
		n *= 10
	}
	fmt.Println("value copies:", nums)
}
