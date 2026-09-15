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
	// Print spaces two operands only when neither is a
	// string; Println always spaces them, and adds \n.
	fmt.Print("Print", 1, 2, "\n")
	fmt.Println("Println", 1, 2)

	// Printf uses verbs; the newline is explicit.
	fmt.Printf("Printf: %s %d\n", "answer", 42)

	// %v value, %+v with field names, %#v Go syntax.
	p := point{X: 3, Y: 4}
	fmt.Printf("%%v: %v\n", p)
	fmt.Printf("%%+v: %+v\n", p)
	fmt.Printf("%%#v: %#v\n", p)

	message := fmt.Sprintf("Sprintf: point %v", p)
	fmt.Println(message)

	// Fprintln targets an io.Writer; here stderr.
	fmt.Fprintln(os.Stderr, "Fprintln: stderr")
}
