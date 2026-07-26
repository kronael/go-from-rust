//go:build ignore

package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	// map[K]V is Go's built-in HashMap.
	scores := map[string]int{"go": 10, "rust": 9}
	score, found := scores["go"]
	fmt.Println("map lookup:", score, found)
	scores["zig"] = 8
	delete(scores, "zig")

	// map[T]struct{} is the zero-payload set.
	seen := map[string]struct{}{"go": {}, "rust": {}}
	_, hasGo := seen["go"]
	fmt.Println("set contains go:", hasGo)
	seen["zig"] = struct{}{}
	delete(seen, "zig")

	// Iteration order is unspecified; sort keys for output.
	fmt.Print("sorted order:")
	for _, language := range slices.Sorted(maps.Keys(scores)) {
		fmt.Printf(" %s=%d", language, scores[language])
	}
	fmt.Println()
}
