//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// A stored heterogeneous Rust tuple usually becomes a
	// named struct. This is a keyed literal: its field names
	// make {Ana 25} meaningful when printed.
	person := Person{Name: "Ana", Age: 25}
	fmt.Println("struct:", person)

	// [N]T stores a fixed number of values of one type.
	point := [2]int{3, 4}
	fmt.Println("array:", point)
}
