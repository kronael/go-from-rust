//go:build ignore

// multireturn.go — multiple return values are NOT a storable tuple.
// Rust: fn divmod() -> (i32, i32) returns a real (i32, i32) you can bind and pass.
// Go: (int, int) is not a type — it's a syntactic position. Two valid consumers:
//   1. multi-variable assignment    q, r := divmod(...)
//   2. full arg-list forwarding     add(divmod(...))
// To store, destructure into locals or box into a struct (see tuples.go).

package main

import "fmt"

func divmod(a, b int) (int, int) { return a / b, a % b }
func add(x, y int) int           { return x + y }

func main() {
	q, r := divmod(17, 5) // consumer 1: multi-assign
	fmt.Println("destructured:", q, r)

	fmt.Println("forwarded:", add(divmod(17, 5))) // consumer 2: full arg-list forward

	// Both lines below are COMPILE ERRORS — uncomment to verify:
	//   var x (int, int) = divmod(1, 1)        // syntax error: no such type
	//   xs := append([]int{}, divmod(1, 1))    // multiple-value in single-value context
}
