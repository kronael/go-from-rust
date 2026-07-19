//go:build ignore

package main

import (
	"fmt"
	"iter"
	"slices"
)

type Countdown struct {
	Start int
}

func (countdown Countdown) All() iter.Seq[int] {
	// Seq[int] is a function that pushes ints into yield. Rust Iterator::next
	// is pull-based instead. false tells the producer that the caller stopped.
	return func(yield func(int) bool) {
		for value := countdown.Start; value > 0; value-- {
			if !yield(value) {
				return
			}
		}
	}
}

func main() {
	countdown := Countdown{Start: 5}

	// Go cannot range over an arbitrary type directly. Expose an iterator method,
	// conventionally named All, and range over its Seq; break makes yield false.
	for value := range countdown.All() {
		fmt.Println("countdown:", value)
		if value == 3 {
			break
		}
	}

	// Standard helpers can consume the same Seq.
	fmt.Println("collected:", slices.Collect(countdown.All()))
}
