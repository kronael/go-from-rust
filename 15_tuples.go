//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

// Employee embeds Person; its fields promote to Employee.
type Employee struct {
	Person
	Role string
}

func main() {
	// Keyed struct literal replaces a heterogeneous tuple.
	person := Person{Name: "Ana", Age: 25}
	fmt.Println("struct:", person)

	// [N]T holds a fixed count of one type.
	point := [2]int{3, 4}
	fmt.Println("array:", point)

	// Go 1.27: a promoted field is a valid literal key
	// directly; no need to nest Person{...} inside.
	worker := Employee{Name: "Bo", Age: 30, Role: "eng"}
	fmt.Println("promoted key:", worker)
}
