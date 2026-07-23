//go:build ignore

package main

import (
	"fmt"
	"os"
)

// The function under test. A real test file would sit beside it as
// abs_test.go and run under `go test`, not `go run`.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// A test case is a named row. Go's idiom is a slice of these structs, one
// row per scenario. Vanilla Rust writes a separate #[test] per case (or
// reaches for a macro like rstest).
type testCase struct {
	name  string
	input int
	want  int
}

func main() {
	cases := []testCase{
		{name: "positive", input: 3, want: 3},
		{name: "negative", input: -4, want: 4},
		{name: "zero", input: 0, want: 0},
	}

	// In a real test this loop body is t.Run(c.name, func(t *testing.T){...}),
	// and a mismatch calls t.Errorf instead of printing. `go test` reports the
	// failing subtest name. Rust asserts with assert_eq!, which panics and
	// stops that test rather than tallying rows like this.
	allPassed := true
	for _, c := range cases {
		got := abs(c.input)
		status := "PASS"
		if got != c.want {
			status = "FAIL"
			allPassed = false
		}
		fmt.Printf("%s %-8s abs(%d) = %d\n", status, c.name, c.input, got)
	}
	fmt.Println("all passed:", allPassed)

	// A real harness fails the process on any mismatch, like `go test`.
	if !allPassed {
		os.Exit(1)
	}
}
