//go:build ignore

// copy built-in and variadic / spread (...).
// Rust analogue: copy ~ slice::copy_from_slice; variadic ~ fn f(args: &[T]) called with a slice.
// Gotcha: copy is bounded by min(len(dst), len(src)) — cap is irrelevant; len-0 dst copies nothing.
// Takeaway: ... both declares variadic params and spreads a slice at the call site.

package main

import "fmt"

func main() {
	src := []int{1, 2, 3}

	dst := make([]int, 3)
	n := copy(dst, src)
	fmt.Printf("copy full: copied %d -> %v\n", n, dst)

	small := make([]int, 2)
	n = copy(small, src)
	fmt.Printf("copy into smaller: copied %d -> %v\n", n, small)

	// classic bug: cap=10 but len=0, so copy copies nothing
	bug := make([]int, 0, 10)
	n = copy(bug, src)
	fmt.Printf("copy into len-0: copied %d -> %v (classic bug)\n", n, bug)

	b := make([]byte, 5)
	copy(b, "hello") // string -> []byte is a special-cased overload
	fmt.Printf("copy from string: %s\n", b)

	fmt.Println("sum(1,2,3) =", sum(1, 2, 3))
	xs := []int{4, 5, 6}
	fmt.Println("sum(xs...) =", sum(xs...)) // spread slice into variadic

	a := []int{1, 2}
	a = append(a, xs...) // spread also works in append
	fmt.Println("append spread:", a)
}

// nums is []int inside the body — Rust: fn sum(nums: &[i32]) -> i32
func sum(nums ...int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}
