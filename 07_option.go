//go:build ignore

package main

import "fmt"

// maybeDouble reports a value and whether it exists. The
// (T, bool) pair is Go's inline optional; no pointer.
func maybeDouble(n int, ok bool) (int, bool) {
	if !ok {
		return 0, false
	}
	return n * 2, true
}

func main() {
	// The bool tells a real 0 apart from "not present".
	value, ok := maybeDouble(0, true)
	fmt.Println("present zero:", value, ok)
	value, ok = maybeDouble(5, false)
	fmt.Println("absent:", value, ok)

	// *T represents optional by pointer: nil is absent.
	// Use it for optional fields or reference semantics.
	var missing *int
	n := 0
	set := &n
	fmt.Println("nil is absent:", missing == nil)
	fmt.Println("ptr present:", set != nil, *set)
}
