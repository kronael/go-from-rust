//go:build ignore

// badidx.go — DELIBERATELY does not compile.
// Demonstrates: "invalid operation: cannot index p (variable of type *[]int)".
// In Rust, *p is auto-derefed through Deref<Target=[i32]>; in Go there is no
// auto-deref for indexing — you must write (*p)[0] explicitly.

package main

func main() {
	s := []int{1, 2, 3}
	p := &s
	_ = p[0] // compile error: cannot index *[]int — use (*p)[0]
}
