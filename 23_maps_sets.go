//go:build ignore

package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	// map[K]V is Go's built-in HashMap-like dictionary.
	scores := map[string]int{"go": 10, "rust": 9}
	score, found := scores["go"]
	fmt.Println("map lookup:", score, found)
	scores["zig"] = 8
	delete(scores, "zig")

	// map[T]struct{} is the usual zero-payload HashSet equivalent.
	seen := map[string]struct{}{"go": {}, "rust": {}}
	_, hasGo := seen["go"]
	fmt.Println("set contains go:", hasGo)
	seen["zig"] = struct{}{}
	delete(seen, "zig")

	// Map iteration order is unspecified. Sort maps.Keys when output merely needs
	// to be deterministic; this allocates a key slice and sorts it.
	fmt.Print("sorted order:")
	for _, language := range slices.Sorted(maps.Keys(scores)) {
		fmt.Printf(" %s=%d", language, scores[language])
	}
	fmt.Println()

	// Go has no ordered-map type. If insertion order is part of the data model,
	// keep a []K beside the map or choose a specialized third-party collection.
}
