//go:build ignore

package main

import "fmt"

func main() {
	label := "before"

	// Defer arguments are evaluated now; closures read variables when run.
	defer fmt.Println("argument:", label)
	defer func() { fmt.Println("closure:", label) }()
	defer fmt.Println("last defer runs first")

	label = "after"
	fmt.Println("body:", label)

	// Rust Drop follows lexical scope. Go defers run at function return, LIFO.
}
