//go:build ignore

package main

import "fmt"

// Rust monomorphizes: one specialized copy per type
// argument, with the type erased by the time it runs.
// Go stencils by GC shape — every pointer-shaped type
// argument shares one instantiation and reaches its type
// through a runtime dictionary.
func describe[T any](v T) string {
	// %T works inside a generic body precisely because
	// the instantiation still carries its type argument.
	return fmt.Sprintf("%T", v)
}

type Meters float64

// Specialization is therefore a runtime type switch, not
// a second impl block: one body that asks, rather than N
// bodies the compiler picked between.
func unit[T any](v T) string {
	switch any(v).(type) {
	case Meters:
		return "m"
	case int:
		return "count"
	default:
		return "?"
	}
}

func main() {
	fmt.Println(describe(1), describe("a"))
	fmt.Println(describe(Meters(2)))

	// *Meters and *int are the same GC shape, so these two
	// calls share one stenciled body — and still report
	// themselves apart, which is the dictionary talking.
	m := Meters(2)
	fmt.Println(describe(&m), describe(new(int)))

	fmt.Println(unit(Meters(5)), unit(3), unit("x"))
}
