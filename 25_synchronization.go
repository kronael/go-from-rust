//go:build ignore

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// The prior lesson uses Mutex for shared invariants. RWMutex permits
	// concurrent readers, but use it only after measuring read-heavy contention.
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
	fmt.Println("rwmutex:", snapshot, "atomic:", requests.Load())

	// The standard library exposes no spinlock. A busy-waiting lock is not a
	// general shortcut around Mutex or the goroutine scheduler.
}
