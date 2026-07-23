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
	// RWMutex: many readers or one writer.
	var readWrite sync.RWMutex
	configuration := "old"
	readWrite.Lock()
	configuration = "new"
	readWrite.Unlock()
	readWrite.RLock()
	snapshot := configuration
	readWrite.RUnlock()

	// atomic.Int64: one indivisible update, no ordering arg.
	var requests atomic.Int64
	requests.Add(1)

	// atomic.Pointer: publish an immutable snapshot.
	var current atomic.Pointer[config]
	current.Store(&config{Host: "localhost", Port: 8080})
	published := current.Load()
	fmt.Println(
		"rwmutex:", snapshot,
		"atomic:", requests.Load(),
		"config:", published.Host, published.Port)
}
