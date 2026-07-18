//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// A stored heterogeneous Rust tuple usually becomes a named struct.
	person := Person{Name: "Ana", Age: 25}
	fmt.Println("struct:", person)

	// [N]T stores a fixed number of values of one type.
	point := [2]int{3, 4}
	fmt.Println("array:", point)

	// Go has no general tuple value. The next lesson covers multiple results,
	// which are function-call syntax rather than stored tuples.
}
