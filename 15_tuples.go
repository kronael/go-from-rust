//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// Keyed struct literal replaces a heterogeneous tuple.
	person := Person{Name: "Ana", Age: 25}
	fmt.Println("struct:", person)

	// [N]T holds a fixed count of one type.
	point := [2]int{3, 4}
	fmt.Println("array:", point)
}
