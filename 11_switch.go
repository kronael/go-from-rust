//go:build ignore

package main

import "fmt"

func main() {
	n := 7

	// Expression switch: comma cases, no fallthrough.
	switch n {
	case 1, 3, 5, 7, 9:
		fmt.Println("small odd")
	case 2, 4, 6, 8:
		fmt.Println("small even")
	default:
		fmt.Println("other integer")
	}

	// Expressionless switch replaces an if/else-if chain.
	switch {
	case n < 0:
		fmt.Println("negative")
	case n%2 == 0:
		fmt.Println("even")
	default:
		fmt.Println("odd")
	}

	// Type switch: value.(type) matches the dynamic type.
	var value any = "go"
	switch value := value.(type) {
	case string:
		fmt.Printf("string %q\n", value)
	default:
		fmt.Printf("other %v\n", value)
	}
}
