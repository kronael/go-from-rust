//go:build ignore

package main

import "fmt"

func divmod(a, b int) (int, int) { return a / b, a % b }

func main() {
	// Multiple results are a call feature, not a value.
	quotient, remainder := divmod(17, 5)
	fmt.Println("assigned:", quotient, remainder)

	// A result list can fill another call's arguments.
	fmt.Println(divmod(20, 6))
}
