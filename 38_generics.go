//go:build ignore

package main

import "fmt"

// Constraint interface: the type set ~int | ~float64.
type Number interface {
	~int | ~float64
}

// sum works for any Number; += compiles for the whole set.
func sum[T Number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

// comparable: the built-in constraint for == and !=.
func index[T comparable](values []T, target T) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}

type Celsius float64

// Box holds one value of any type.
type Box[T any] struct {
	value T
}

// Go 1.27: MapTo declares its own T2, beyond Box's T.
// An interface method still cannot add a type param.
func (b Box[T]) MapTo[T2 any](f func(T) T2) Box[T2] {
	return Box[T2]{value: f(b.value)}
}

func main() {
	// Type arguments inferred; no turbofish.
	fmt.Println("ints:", sum([]int{1, 2, 3}))
	fmt.Println("floats:", sum([]float64{1.5, 2.5}))

	// A named type satisfies Number via ~float64.
	fmt.Println("named:", sum([]Celsius{20, 1.5}))

	// comparable admits strings; Number excludes them.
	fmt.Println("index:", index([]string{"a", "b", "c"}, "b"))

	// The method's T2 differs per call: int, then string.
	box := Box[int]{value: 21}
	doubled := box.MapTo(func(v int) int { return v * 2 })
	fmt.Println("mapped int:", doubled.value)
	label := box.MapTo(func(v int) string {
		return fmt.Sprintf("n=%d", v)
	})
	fmt.Println("mapped string:", label.value)
}
