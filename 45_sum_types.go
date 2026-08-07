//go:build ignore

package main

import (
	"fmt"
	"math"
)

// Rust: enum Shape { Circle(f64), Rect(f64, f64) }
// Variant payloads: plain structs, no shared base type.
type Circle struct{ R float64 }
type Rect struct{ W, H float64 }

// Shape holds each variant behind its own optional
// pointer. No Kind field, so presence is the tag: it
// can never disagree with which payload is set.
type Shape struct {
	Circle *Circle
	Rect   *Rect
}

// area checks each variant independently, not with a
// switch: none-set and both-set are ordinary states
// here, not bugs to guard against.
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

// Match is the one encoding Go checks exhaustively: a
// new variant adds a parameter, so every call site
// stops compiling until it handles the new case. The
// price is that neither-set and both-set now need an
// answer, and no return here leaves the caller.
func Match[R any](
	shape Shape,
	onCircle func(Circle) R,
	onRect func(Rect) R,
	otherwise func() R,
) R {
	switch {
	case shape.Circle != nil && shape.Rect == nil:
		return onCircle(*shape.Circle)
	case shape.Rect != nil && shape.Circle == nil:
		return onRect(*shape.Rect)
	default:
		return otherwise()
	}
}

func main() {
	both := Shape{
		Circle: &Circle{R: 1},
		Rect:   &Rect{W: 2, H: 2},
	}
	shapes := []Shape{
		{Circle: &Circle{R: 2}},
		{Rect: &Rect{W: 3, H: 4}},
		{},
		both,
	}

	for _, shape := range shapes {
		name := Match(shape,
			func(Circle) string { return "circle" },
			func(Rect) string { return "rect" },
			func() string { return "not exactly one" },
		)
		fmt.Printf("%-15s %.2f\n", name, area(shape))
	}
}
