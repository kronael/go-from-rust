//go:build ignore

package main

import "fmt"

// Typed int + const block + iota (0,1,2...) is Go's enum.
type Color int

const (
	Red Color = iota
	Green
	Blue
)

// String makes Color a fmt.Stringer; fmt calls it for %v.
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

// Bit flags: 1 << iota gives each constant a distinct bit.
type Perm int

const (
	Read Perm = 1 << iota
	Write
	Exec
)

func main() {
	// Stringer is called automatically when printing.
	fmt.Println("colors:", Red, Green, Blue)

	// Color(9) is still valid; enums are not exhaustive.
	fmt.Println("unknown:", Color(9))

	perms := Read | Exec
	fmt.Println("can write:", perms&Write != 0)
	fmt.Println("can exec:", perms&Exec != 0)
}
