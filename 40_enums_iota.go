//go:build ignore

package main

import "fmt"

// A typed int plus a const block is Go's enum. iota starts
// at 0 and adds 1 per constant, so Red=0, Green=1, Blue=2.
// Rust's `enum Color` is a real sum type; a Go enum is just
// an int, with no exhaustiveness checking.
type Color int

const (
	Red Color = iota
	Green
	Blue
)

// String makes Color satisfy fmt.Stringer, Go's analogue of
// Rust's Display impl (not Debug, which would print the
// variant name). fmt calls it for %v and %s automatically.
func (color Color) String() string {
	switch color {
	case Red:
		return "red"
	case Green:
		return "green"
	case Blue:
		return "blue"
	default:
		return fmt.Sprintf("Color(%d)", int(color))
	}
}

// Bit flags shift 1 left by iota, so each constant is a
// distinct bit: the Go analogue of the bitflags crate.
// Combine with |, test with &.
type Perm int

const (
	Read Perm = 1 << iota
	Write
	Exec
)

func main() {
	// Stringer is invoked for printing; no explicit call
	// needed.
	fmt.Println("colors:", Red, Green, Blue)

	// An out-of-range value is still valid; the default arm
	// handles it. Go won't stop you making Color(9), unlike a
	// Rust match on a real enum.
	fmt.Println("unknown:", Color(9))

	perms := Read | Exec
	fmt.Println("can write:", perms&Write != 0)
	fmt.Println("can exec:", perms&Exec != 0)
}
