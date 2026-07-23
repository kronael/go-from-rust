//go:build ignore

package main

import (
	"fmt"
	"os"
)

type point struct {
	X int
	Y int
}

func main() {
	// Print writes arguments as-is here. Println inserts spaces and a
	// final newline.
	fmt.Print("Print")
	fmt.Println(" + Println", 1, 2)

	// Printf uses verbs, and the newline is explicit.
	fmt.Printf("Printf: %s %d\n", "answer", 42)

	// %v prints values, %+v adds struct field names, and %#v prints Go syntax.
	// These do not split directly into Rust's Display and Debug.
	p := point{X: 3, Y: 4}
	fmt.Printf("%%v: %v\n", p)
	fmt.Printf("%%+v: %+v\n", p)
	fmt.Printf("%%#v: %#v\n", p)

	message := fmt.Sprintf("Sprintf: point %v", p)
	fmt.Println(message)

	// Fprintln targets an io.Writer. This line uses stderr, whose ordering against
	// the stdout lines is not guaranteed when a caller captures them separately.
	fmt.Fprintln(os.Stderr, "Fprintln: stderr")
}
