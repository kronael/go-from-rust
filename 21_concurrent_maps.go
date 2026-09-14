//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	// Typed map plus Mutex: the default, keeps static types.
	counts := map[string]int{}
	var mutex sync.Mutex
	var workers sync.WaitGroup
	for range 2 {
		// Add before go; defer Done/Unlock; Wait joins.
		workers.Add(1)
		go func() {
			defer workers.Done()
			mutex.Lock()
			defer mutex.Unlock()
			counts["go"]++
		}()
	}
	workers.Wait()
	fmt.Println("map plus mutex:", counts["go"])

	// sync.Map: for write-once/read-many; stores any.
	var cache sync.Map
	cache.Store("go", 10)
	value, found := cache.Load("go")
	fmt.Println("sync.Map:", value.(int), found)
}
