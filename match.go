//go:build ignore

package main

import "fmt"

func main() {
	// Rust slice patterns become explicit indexing and slicing.
	values := []int{1, 2, 3, 4}
	head, tail := values[0], values[1:]
	fmt.Println("head:", head, "tail:", tail)

	// Interface type switches are open, not exhaustive like enum matches.
	for _, value := range []any{42, "hi", true, 3.14} {
		switch value := value.(type) {
		case int:
			fmt.Printf("int %d\n", value)
		case string:
			fmt.Printf("string %q\n", value)
		default:
			fmt.Printf("other %v\n", value)
		}
	}

	// A conditionless switch replaces an if/else-if chain.
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

	// Value cases accept comma-separated values and do not fall through.
	switch n {
	case 1, 3, 5, 7, 9:
		fmt.Println("small odd")
	case 2, 4, 6, 8:
		fmt.Println("small even")
	}
}
