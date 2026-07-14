//go:build ignore

package main

import "fmt"

func main() {
	n := 7

	// Rust's if is an expression; Go's if is a statement.
	var sign string
	if n >= 0 {
		sign = "pos"
	} else {
		sign = "neg"
	}
	fmt.Println("sign:", sign)
}
