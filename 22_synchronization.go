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

const readers = 8

func main() {
	// atomic.Pointer publishes a whole snapshot in one
	// store, so a reader never sees a half-written value.
	var published atomic.Pointer[config]
	published.Store(&config{Host: "localhost", Port: 8080})

	// RWMutex: many readers at once, or one writer alone.
	// One goroutine writes while eight read, so deleting
	// these locks makes `go run -race` report a race.
	var readWrite sync.RWMutex
	settings := config{Host: "old", Port: 1}

	// atomic.Int64 counts from every goroutine at once.
	// A plain int++ would lose updates here.
	var seen atomic.Int64

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		readWrite.Lock()
		settings = config{Host: "localhost", Port: 8080}
		readWrite.Unlock()
	}()

	for range readers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readWrite.RLock()
			_ = settings.Host
			readWrite.RUnlock()
			seen.Add(1)
		}()
	}
	wg.Wait()

	// Only one writer ran, so the final state is fixed
	// even though the interleaving was not.
	snapshot := published.Load()
	fmt.Println("readers counted:", seen.Load())
	fmt.Println("settings:", settings.Host, settings.Port)
	fmt.Println("published:", snapshot.Host, snapshot.Port)
}
