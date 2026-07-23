//go:build ignore

package main

import (
	"fmt"
	"os"
)

// The function under test.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// A test case is a named row; the slice drives one loop.
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

	// Real tests use t.Run per row; t.Errorf on mismatch.
	allPassed := true
	for _, c := range cases {
		got := abs(c.input)
		status := "PASS"
		if got != c.want {
			status = "FAIL"
			allPassed = false
		}
		fmt.Printf("%s %-8s abs(%d) = %d\n",
			status, c.name, c.input, got)
	}
	fmt.Println("all passed:", allPassed)

	// A real harness exits nonzero on any mismatch.
	if !allPassed {
		os.Exit(1)
	}
}
