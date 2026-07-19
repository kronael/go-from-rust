//go:build ignore

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type config struct {
	Host string
	Port int
}

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

	// Typed atomics make one operation on one value indivisible. Go's standard
	// atomic operations are sequentially consistent; unlike Rust, no ordering
	// argument appears at each call.
	var requests atomic.Int64
	requests.Add(1)

	// With one writer, atomic.Pointer can cheaply publish an immutable snapshot
	// to many readers. Do not mutate a config after Store: the pointer operation
	// does not protect its fields.
	var current atomic.Pointer[config]
	current.Store(&config{Host: "localhost", Port: 8080})
	published := current.Load()
	fmt.Println("rwmutex:", snapshot, "atomic:", requests.Load(), "config:", published.Host, published.Port)

	// GC reclaims an old pure-data snapshot after readers release it. If snapshots
	// own files or sockets, Store cannot say when every old reader is finished;
	// use an explicit lifetime protocol instead of relying on this pattern alone.

	// Use Mutex when several mutable fields form one invariant. Use RWMutex only
	// after measured reader contention; use atomics for one value or snapshot.
}
