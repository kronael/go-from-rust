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

	// Comparable structs and arrays can be map keys.
	roles := map[Person]string{person: "admin"}
	visited := map[[2]int]bool{point: true}
	fmt.Println("struct key:", roles[person])
	fmt.Println("array key:", visited[point])
}
