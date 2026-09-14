//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func divmod(a, b int) (int, int) { return a / b, a % b }

func main() {
	// Keyed struct literal replaces a heterogeneous tuple.
	person := Person{Name: "Ana", Age: 25}
	fmt.Println("struct:", person)

	// [N]T holds a fixed count of one type.
	point := [2]int{3, 4}
	fmt.Println("array:", point)

	// Multiple results are a call feature, not a value.
	// Rust returns one tuple; Go hands back two results,
	// and there is no (int, int) value to store.
	quotient, remainder := divmod(17, 5)
	fmt.Println("assigned:", quotient, remainder)

	// A result list can fill another call's arguments.
	fmt.Println(divmod(20, 6))
}
