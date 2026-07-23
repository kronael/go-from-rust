//go:build ignore

package main

import "fmt"

func main() {
	values := map[string]int{"zero": 0}

	// comma-ok: (value, bool) distinguishes zero from absent.
	value, ok := values["zero"]
	fmt.Println("stored zero:", value, ok)
	value, ok = values["missing"]
	fmt.Println("missing:", value, ok)

	// *T is a nullable pointer, not a general Option<T>.
	var absent *int
	zero := 0
	present := &zero
	fmt.Println("absent:", absent == nil)
	fmt.Println("present zero:", present != nil, *present)
}
