//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	// A typed map plus Mutex is the default: it preserves
	// static types and can protect invariants spanning more
	// than one operation or value. Ordinary maps must not be
	// written without synchronization while another goroutine
	// reads or writes them; this resembles
	// Arc<Mutex<HashMap<...>>> in Rust.
	counts := map[string]int{}
	var mutex sync.Mutex
	var workers sync.WaitGroup
	for range 2 {
		// Add registers work before go starts a goroutine. Each
		// goroutine defers Done and Unlock; Wait blocks until
		// both have called Done.
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

	// sync.Map is specialized for write-once/read-many caches
	// or concurrent work on disjoint keys. It stores any, so
	// Load loses the map's static type.
	var cache sync.Map
	cache.Store("go", 10)
	value, found := cache.Load("go")
	fmt.Println("sync.Map:", value.(int), found)
}
