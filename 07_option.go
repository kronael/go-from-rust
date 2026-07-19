//go:build ignore

package main

import "fmt"

func main() {
	values := map[string]int{"zero": 0}

	// ok distinguishes a missing key from a stored zero. This (T, bool) pattern
	// is Go's common Result/Option-like lookup form and needs no pointer itself.
	value, ok := values["zero"]
	fmt.Println("stored zero:", value, ok)
	value, ok = values["missing"]
	fmt.Println("missing:", value, ok)

	// *T is a nullable pointer, not a general Option<T>. Reading it adds an
	// indirection. It does not inherently allocate; the compiler may move an
	// escaping pointee to the heap, so measure if this representation is hot.
	var absent *int
	zero := 0
	present := &zero
	fmt.Println("absent:", absent == nil)
	fmt.Println("present zero:", present != nil, *present)
}
