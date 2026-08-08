//go:build ignore

package main

import (
	"fmt"
	"strings"
)

// Rust splits Copy (implicit, bitwise) from Clone
// (explicit, and often deep). Go has one behavior:
// assignment, arguments, channel sends, and range values
// all copy the struct bitwise, always shallowly.
type Doc struct {
	Title string
	Tags  []string
}

// A type whose value must not be copied says so at run
// time or not at all. Builder records its own address.
type Report struct {
	body strings.Builder
}

func main() {
	original := Doc{Title: "spec", Tags: []string{"go"}}
	copied := original

	// The string field is independent; the slice header
	// was copied but still points at the same array.
	copied.Title = "draft"
	copied.Tags[0] = "rust"
	fmt.Println("original:", original.Title, original.Tags)
	fmt.Println("copied:  ", copied.Title, copied.Tags)

	var r Report
	r.body.WriteString("first")

	// Rust's move semantics make this a compile error.
	// Go compiles it and fails on the next write.
	broken := r
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		broken.body.WriteString("second")
	}()

	// sync.Mutex and sync.WaitGroup have no such check —
	// a copy silently guards nothing. `go vet` catches the
	// common cases (copylocks) but it is a heuristic, so
	// pointer receivers are the rule for those types.
	fmt.Println("original still works:", r.body.String())
}
