//go:build ignore

package main

import "fmt"

func main() {
	label := "before"

	// Defer args evaluate now; closures read vars later.
	defer fmt.Println("argument:", label)
	defer func() { fmt.Println("closure:", label) }()
	defer fmt.Println("last defer runs first")

	label = "after"
	fmt.Println("body:", label)
}
