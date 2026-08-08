//go:build ignore

package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// Rust: String and &str are validated UTF-8 by
	// construction. A Go string is an immutable byte
	// slice with no validity guarantee at all.
	s := "héllo"
	fmt.Println("bytes:", len(s))
	fmt.Println("runes:", utf8.RuneCountInString(s))

	// Indexing yields one byte, not a character. Rust
	// refuses to index a str for exactly this reason.
	fmt.Printf("s[1] = %d %q\n", s[1], s[1])

	// Invalid UTF-8 is still a string. Nothing rejects it,
	// so validation is yours to do at the boundary.
	bad := string([]byte{0x41, 0xff, 0x42})
	fmt.Println("valid:", utf8.ValidString(bad))

	// range decodes, substituting U+FFFD and advancing a
	// single byte over the bad one — it never fails.
	for i, r := range bad {
		fmt.Printf("  %d: %U\n", i, r)
	}

	// []rune substitutes the same way, so the round trip
	// is lossy rather than an error.
	fmt.Printf("round trip: %q\n", string([]rune(bad)))

	// rune is int32 and holds anything an int32 holds,
	// including values Rust's char cannot represent.
	half := rune(0xD800)
	fmt.Println("surrogate half:", half)
	fmt.Println("valid rune:", utf8.ValidRune(half))
}
