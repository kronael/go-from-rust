//go:build ignore

package main

import "fmt"

// List is generic over element type E.
type List[E any] []E

// Apply: a generic method — its own type
// parameter F, beyond the receiver's E.
// Go 1.27 lets a method add type params;
// earlier Go allowed this only on funcs.
func (l List[E]) Apply[F any](
	f func(E) F,
) List[F] {
	out := make(List[F], len(l))
	for i, v := range l {
		out[i] = f(v)
	}
	return out
}

func main() {
	nums := List[int]{1, 2, 3}

	// Instantiation 1: F inferred as int.
	doubled := nums.Apply(func(n int) int {
		return n * 2
	})
	fmt.Println("doubled:", doubled)

	// Instantiation 2: F inferred as
	// string, same method, same receiver.
	labels := nums.Apply(func(n int) string {
		return fmt.Sprintf("n%d", n)
	})
	fmt.Println("labels:", labels)
}
