//go:build ignore

package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Sentinel errors created with errors.New give callers a
// stable identity.
var ErrNegative = errors.New("negative count")

func parseCount(text string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("parse count %q: %w", text, err)
	}
	if number < 0 {
		return 0, fmt.Errorf("count %d: %w", number, ErrNegative)
	}
	return number, nil
}

func main() {
	// Rust often uses enum variants. Go uses sentinel errors
	// when callers need identity, then adds context with %w
	// without losing that identity. Atoi's *NumError wraps
	// strconv.ErrSyntax, so errors.Is finds it through both
	// wraps.
	for _, text := range []string{"42", "nope", "-1"} {
		number, err := parseCount(text)
		switch {
		case err == nil:
			fmt.Println("count:", number)
		case errors.Is(err, ErrNegative):
			fmt.Println("domain error:", err)
		case errors.Is(err, strconv.ErrSyntax):
			fmt.Println("syntax error:", err)
		default:
			fmt.Println("other error:", err)
		}
	}
}
