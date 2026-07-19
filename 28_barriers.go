//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	// A barrier has two roles: count arrivals, then release every waiter.
	// WaitGroup is a counter; closing a channel is a one-to-many broadcast.
	const workerCount = 3
	ready := sync.WaitGroup{}
	ready.Add(workerCount)
	start := make(chan struct{})
	results := make([]int, workerCount)

	var done sync.WaitGroup
	for worker := range workerCount {
		done.Go(func() {
			ready.Done()
			<-start
			results[worker] = worker * worker
		})
	}

	// Wait until every worker reaches the gate. A coordinator could update shared
	// state here; close then releases all blocked receives at once.
	ready.Wait()
	close(start)
	done.Wait()
	fmt.Println("after barrier:", results)

	// A WaitGroup alone can release immediately at the last arrival. Channels
	// alone can count N arrival messages and close a separate release channel.
	// This combination makes both roles explicit. Unlike Rust's reusable Barrier,
	// it is one-shot; cyclic barriers need careful sync.Cond state.

	// Primitive map: Mutex protects mutable invariants; atomics publish one value
	// or immutable snapshot; channels transfer work, apply backpressure, or signal;
	// WaitGroup joins finite work; Once initializes once; Cond waits for a repeated
	// state condition. Prefer the primitive that states the coordination rule.
}
