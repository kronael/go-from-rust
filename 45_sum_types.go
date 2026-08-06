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

func main() {
	circle := Shape{Circle: &Circle{R: 2}}
	rect := Shape{Rect: &Rect{W: 3, H: 4}}
	empty := Shape{}
	both := Shape{
		Circle: &Circle{R: 1},
		Rect:   &Rect{W: 2, H: 2},
	}

	fmt.Printf("circle: %.2f\n", area(circle))
	fmt.Printf("rect:   %.2f\n", area(rect))
	fmt.Printf("empty:  %.2f\n", area(empty))
	fmt.Printf("both:   %.2f\n", area(both))
}
