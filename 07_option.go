//go:build ignore

package main

import "fmt"

func main() {
	values := map[string]int{"zero": 0}

	// Bare index yields the zero value, not a None.
	fmt.Println("bare missing:", values["missing"])

	// comma-ok tells a stored zero apart from absent.
	value, ok := values["zero"]
	fmt.Println("stored zero:", value, ok)
	if _, ok := values["missing"]; !ok {
		fmt.Println("missing: absent")
	}

	// *T: nullable pointer for optional/reference use.
	var absent *int
	zero := 0
	present := &zero
	fmt.Println("absent:", absent == nil)
	fmt.Println("present zero:", present != nil, *present)
}
