//go:build ignore

package main

import "fmt"

// A Go type parameter can be converted to any without
// losing its concrete type. %T observes that type at run time.
func describe[T any](v T) string {
	return fmt.Sprintf("%T", v)
}

type Meters float64

// Go has no generic specialization syntax. When behavior
// really depends on T, a type switch asks at run time.
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
	fmt.Println("concrete types remain:",
		describe(1), describe("a"), describe(Meters(2)))
	fmt.Println("runtime type switch:",
		unit(Meters(5)), unit(3), unit("x"))
}
