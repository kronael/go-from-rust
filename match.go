//go:build ignore

// match.go — Go's answer to Rust's match/destructuring.
// No pattern destructuring: slice heads/tails are extracted by index.
// Type switch is the closest to Rust's match on enums; conditionless switch
// replaces chains of if/else if. Cases do NOT fall through by default.
// Takeaway: switch on type for interface dispatch, switch {} for bool chains.

package main

import "fmt"

func main() {
	// No destructuring — extract by index manually.
	s := []int{1, 2, 3, 4}
	head, tail := s[0], s[1:]
	fmt.Println("head:", head, "tail:", tail)

	a, b := s[0], s[1]
	fmt.Println("a,b:", a, b)

	// Multi-assignment exists (functions, swaps) but it's not slice destructuring.
	x, y := pair()
	x, y = y, x
	fmt.Println("swapped:", x, y)

	// Type switch: closest to Rust's `match` on an enum / dyn Trait.
	for _, v := range []any{42, "hi", true, 3.14} {
		switch t := v.(type) {
		case int:
			fmt.Printf("int %d\n", t)
		case string:
			fmt.Printf("string %q\n", t)
		case bool:
			fmt.Printf("bool %v\n", t)
		default:
			fmt.Printf("other %v\n", t)
		}
	}

	// Conditionless switch: clean replacement for if/else chains.
	n := 7
	switch {
	case n < 0:
		fmt.Println("negative")
	case n == 0:
		fmt.Println("zero")
	case n%2 == 0:
		fmt.Println("even")
	default:
		fmt.Println("odd")
	}

	// Multi-value cases; no implicit fallthrough (unlike C).
	switch n {
	case 1, 3, 5, 7, 9:
		fmt.Println("small odd")
	case 2, 4, 6, 8:
		fmt.Println("small even")
	}
}

func pair() (int, int) { return 1, 2 }
