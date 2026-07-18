//go:build ignore

package main

import "fmt"

func main() {
	n := 7

	// 1. Expression switches accept comma-separated cases and do not fall through.
	switch n {
	case 1, 3, 5, 7, 9:
		fmt.Println("small odd")
	case 2, 4, 6, 8:
		fmt.Println("small even")
	default:
		fmt.Println("other integer")
	}

	// 2. A switch without an expression replaces an if/else-if chain.
	switch {
	case n < 0:
		fmt.Println("negative")
	case n%2 == 0:
		fmt.Println("even")
	default:
		fmt.Println("odd")
	}

	// 3. A type switch inspects an interface's dynamic type. It is open,
	// not exhaustive like matching a Rust enum.
	var value any = "go"
	switch value := value.(type) {
	case string:
		fmt.Printf("string %q\n", value)
	default:
		fmt.Printf("other %v\n", value)
	}
}
