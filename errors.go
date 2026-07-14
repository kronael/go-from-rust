//go:build ignore

package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parseNumber(text string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		// %w preserves the cause for errors.Is instead of flattening it to text.
		return 0, fmt.Errorf("parse %q: %w", text, err)
	}
	return number, nil
}

func main() {
	// Rust's Result maps to a value and error returned as ordinary values.
	for _, text := range []string{"42", "nope"} {
		number, err := parseNumber(text)
		if err != nil {
			fmt.Println("error:", err)
			fmt.Println("is syntax error:", errors.Is(err, strconv.ErrSyntax))
			continue
		}
		fmt.Println("value:", number)
	}
}
