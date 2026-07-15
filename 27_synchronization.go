//go:build ignore

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// Mutex is the default for shared state and multi-value invariants.
	var mutex sync.Mutex
	value := 0
	mutex.Lock()
	value++
	mutex.Unlock()

	// RWMutex permits concurrent readers and one writer. It is not inherently
	// faster than Mutex; use it only after measuring read-heavy contention.
	var readWrite sync.RWMutex
	configuration := "old"
	readWrite.Lock()
	configuration = "new"
	readWrite.Unlock()
	readWrite.RLock()
	snapshot := configuration
	readWrite.RUnlock()

	// Typed atomic counters fit one independent word. atomic.Pointer and Value
	// can publish immutable snapshots, but not coordinate separate mutations.
	var requests atomic.Int64
	requests.Add(1)
	fmt.Println("mutex:", value, "rwmutex:", snapshot, "atomic:", requests.Load())

	// The standard library exposes no spinlock. A busy-waiting lock is not a
	// general shortcut around Mutex or the goroutine scheduler.
}
