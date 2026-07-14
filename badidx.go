//go:build ignore

package main

func main() {
	s := []int{1, 2, 3}
	p := &s

	// Go does not index through *[]T; normally pass []T.
	_ = p[0]
}
