//go:build ignore

// option.go — Go has no Option<T>; absence is expressed two ways.
// Rust: Option<i32> with Some/None. Go: (T, bool) for lookup/conversion,
// *T for struct fields where nil = absent and zero-value must stay reachable.
// A generic Option type is possible with generics but not idiomatic.
// Takeaway: prefer comma-ok; reach for *T only when zero ≠ absent matters.

package main

import "fmt"

// 1. comma-ok: the idiomatic "optional T" for lookup / conversion.
func find(m map[string]int, k string) (int, bool) {
	v, ok := m[k]
	return v, ok
}

// generic Option — valid, but you'll rarely see it in real Go code.
type Option[T any] struct {
	val T
	ok  bool
}

func Some[T any](v T) Option[T] { return Option[T]{v, true} }
func None[T any]() Option[T]    { return Option[T]{} }

func main() {
	m := map[string]int{"a": 1}

	if v, ok := find(m, "a"); ok {
		fmt.Println("found:", v)
	}
	if _, ok := find(m, "z"); !ok {
		fmt.Println("absent: z")
	}

	// 2. *T as nullable: nil = absent, non-nil = present (even if the value is 0).
	var none *int
	zero := 0
	some := &zero // &0 is present; its value happens to be 0
	fmt.Println("none is nil?", none == nil)
	fmt.Println("some is nil?", some == nil, "value:", *some)

	report := func(p *int) {
		if p == nil {
			fmt.Println("  field not set")
		} else {
			fmt.Println("  field set to", *p)
		}
	}
	report(none)
	report(some)

	// 3. hand-rolled Option (works, rarely seen in practice)
	o := Some(42)
	fmt.Println("Option ok?", o.ok, "val:", o.val)
}
