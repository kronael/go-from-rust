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
	// Seq[int]: pushes ints into yield; false means stop.
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

	// range over All()'s Seq; break makes yield false.
	for value := range countdown.All() {
		fmt.Println("countdown:", value)
		if value == 3 {
			break
		}
	}

	// Standard helpers can consume the same Seq.
	fmt.Println("collected:", slices.Collect(countdown.All()))
}
