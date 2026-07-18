//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
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

	// Wait until every worker reaches the gate, then close broadcasts release.
	ready.Wait()
	close(start)
	done.Wait()
	fmt.Println("after barrier:", results)

	// Go has no reusable Barrier type. WaitGroup plus a closed channel is a
	// simple one-shot barrier; cyclic barriers need careful sync.Cond state.
}
