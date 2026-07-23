//go:build ignore

package main

import "fmt"

type problem struct{}

func (*problem) Error() string { return "problem" }

func main() {
	var pointer *problem
	var err error = pointer

	// Storing a nil *problem sets the type, so err != nil.
	fmt.Println("plain pointer is nil:", pointer == nil)
	fmt.Printf("interface dynamic type: %T\n", err)
	fmt.Println("interface value is nil pointer:",
		err.(*problem) == nil)
	fmt.Println("interface itself is nil:", err == nil)

	// An unassigned interface has neither part; it is nil.
	var empty error
	fmt.Printf("empty dynamic type: %T\n", empty)
	fmt.Println("empty interface is nil:", empty == nil)
}
