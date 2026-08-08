//go:build ignore

package main

import (
	"errors"
	"fmt"
	"math"
)

// Rust: enum Shape { Circle(f64), Rect(f64, f64) }
// Variants are plain structs, no shared base type.
type Circle struct{ R float64 }
type Rect struct{ W, H float64 }

// No Kind field: pointer presence is the tag, so it
// cannot disagree with the payload.
type Shape struct {
	Circle *Circle
	Rect   *Rect
}

// Independent checks, not a switch: none-set and
// both-set are ordinary states, not bugs.
func area(shape Shape) float64 {
	total := 0.0
	if shape.Circle != nil {
		r := shape.Circle.R
		total += math.Pi * r * r
	}
	if shape.Rect != nil {
		total += shape.Rect.W * shape.Rect.H
	}
	return total
}

// One handler per variant is what Go checks
// exhaustively: a new variant adds a parameter, so
// every call site stops compiling. Exactly-one is not
// checkable, so Match reports it as an error.
func Match[R any](
	shape Shape,
	onCircle func(Circle) R,
	onRect func(Rect) R,
) (R, error) {
	var zero R
	switch {
	case shape.Circle != nil && shape.Rect != nil:
		return zero, errors.New("both variants set")
	case shape.Circle != nil:
		return onCircle(*shape.Circle), nil
	case shape.Rect != nil:
		return onRect(*shape.Rect), nil
	default:
		return zero, errors.New("no variant set")
	}
}

func main() {
	shapes := []Shape{
		{Circle: &Circle{R: 2}},
		{Rect: &Rect{W: 3, H: 4}},
		{},
		{Circle: &Circle{R: 1}, Rect: &Rect{W: 2, H: 2}},
	}

	for _, shape := range shapes {
		label, err := Match(shape,
			func(Circle) string { return "circle" },
			func(Rect) string { return "rect" },
		)
		if err != nil {
			label = err.Error()
		}
		fmt.Printf("%-17s %.2f\n", label, area(shape))
	}
}
