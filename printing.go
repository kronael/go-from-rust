//go:build ignore

// printing.go — fmt package: Print/Println/Printf/Sprintf/Fprintln and common verbs.
// Rust: println!("{}", x) / format!("{:?}", x). Go: fmt.Printf("%v", x) / %+v / %#v.
// Key verbs: %v default, %+v with field names, %#v Go-syntax repr, %T type name,
//            %d int, %s string, %q quoted string, %x hex, %t bool, %f float.
// Takeaway: Fprintln(w, ...) writes to any io.Writer — use os.Stderr for errors.

package main

import (
	"fmt"
	"os"
)

type Point struct{ X, Y int }

func main() {
	fmt.Println("Println: args spaced, newline added", 1, 2, 3)
	fmt.Print("Print: no newline, ", "no auto-spaces between strings\n")
	fmt.Printf("Printf: formatted with verbs: %d and %s\n", 42, "hi")

	s := fmt.Sprintf("user #%d", 7) // format into a string, don't print
	fmt.Println("Sprintf ->", s)

	p := Point{1, 2}
	fmt.Printf("%%v   default:   %v\n", p)
	fmt.Printf("%%+v  w/ fields:  %+v\n", p)
	fmt.Printf("%%#v  Go syntax:  %#v\n", p)
	fmt.Printf("%%T   type:       %T\n", p)
	fmt.Printf("%%d %%s %%q:      %d %s %q\n", 10, "go", "go")
	fmt.Printf("%%.2f %%t:        %.2f %t\n", 3.14159, true)
	fmt.Printf("%%x hex:         %x\n", 255)

	fmt.Fprintln(os.Stderr, "this goes to stderr") // any io.Writer works
}
