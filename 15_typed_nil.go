//go:build ignore

package main

import "fmt"

type problem struct{}

func (*problem) Error() string { return "problem" }

// find succeeds, but returns a typed nil pointer as error.
func find() error {
	// p is nil; nothing went wrong.
	var p *problem
	// Returning it types the error, so callers see != nil.
	return p
}

func main() {
	err := find()

	// find "returned nil" yet err is non-nil: the trap.
	fmt.Println("err != nil:", err != nil)
	fmt.Printf("dynamic type: %T\n", err)
	fmt.Println("ptr is nil:", err.(*problem) == nil)

	// An interface is nil only with no type and no value.
	var ok error
	fmt.Println("bare error is nil:", ok == nil)
}
