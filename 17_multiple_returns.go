//go:build ignore

package main

import "fmt"

func divmod(a, b int) (int, int) { return a / b, a % b }

func main() {
	// Multiple results are a function-call language feature, not a data
	// structure. Rust returns one tuple value here; Go returns two results.
	quotient, remainder := divmod(17, 5)
	fmt.Println("assigned:", quotient, remainder)

	// A result list can directly fill another call's argument list.
	fmt.Println(divmod(20, 6))

	// There is no (int, int) value to store; use a struct for stored data.
}
