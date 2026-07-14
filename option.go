//go:build ignore

package main

import "fmt"

func main() {
	values := map[string]int{"zero": 0}

	// A Rust map lookup returning Option maps to Go's comma-ok form.
	value, ok := values["zero"]
	fmt.Println("stored zero:", value, ok)
	value, ok = values["missing"]
	fmt.Println("missing:", value, ok)

	// Some APIs use *T when nil must differ from T's zero value.
	var absent *int
	zero := 0
	present := &zero
	fmt.Println("absent:", absent == nil)
	fmt.Println("present zero:", present != nil, *present)
}
