//go:build ignore

package main

import "fmt"

func divmod(a, b int) (int, int) { return a / b, a % b }

func main() {
	// Rust returns a tuple value; Go returns a list of results.
	quotient, remainder := divmod(17, 5)
	fmt.Println("assigned:", quotient, remainder)

	// A result list can directly fill another call's argument list.
	fmt.Println(divmod(20, 6))

	// Multiple results are not one storable tuple value; use a struct to store one.
}
