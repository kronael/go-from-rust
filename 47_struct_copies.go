//go:build ignore

package main

import (
	"fmt"
	"strings"
)

// Rust distinguishes an implicit Copy from an explicit
// Clone. Go assignment copies the struct value; fields that
// refer to other storage still refer to that same storage.
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
	fmt.Println("titles differ:", original.Title, copied.Title)
	fmt.Println("tags share storage:", original.Tags, copied.Tags)

	var r Report
	r.body.WriteString("first")

	// This copy is legal Go too. strings.Builder adds its own
	// runtime check, so the later write exposes the bad copy.
	broken := r
	func() {
		defer func() {
			fmt.Println("copy compiled; write panicked:", recover())
		}()
		broken.body.WriteString("second")
	}()

	fmt.Println("original still works:", r.body.String())
}
