//go:build ignore

// ternary.go — Go deliberately omits the ?: operator.
// Rust: let s = if n >= 0 { "pos" } else { "neg" }; (if is an expression).
// Go: if/else is a statement; assign inside each branch.
// For numbers, min/max builtins (Go 1.21+) cover the common case.
// A generic cond() helper works but evaluates BOTH args before the call —
// no short-circuit, so a panicking or side-effecting branch is always executed.

package main

import "fmt"

func cond[T any](c bool, a, b T) T {
	if c {
		return a
	}
	return b
}

func main() {
	n := 7

	var sign string
	if n >= 0 {
		sign = "pos"
	} else {
		sign = "neg"
	}
	fmt.Println("if/else:", sign)

	fmt.Println("max/min builtins:", max(3, 9), min(3, 9)) // Go 1.21+

	// Both branch args are evaluated before cond runs — the print proves it.
	pick := func(label string) string {
		fmt.Println("  evaluated branch:", label)
		return label
	}
	fmt.Println("generic helper:", cond(n > 0, pick("pos"), pick("neg")))
}
