//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	// Barrier: WaitGroup counts, channel close broadcasts.
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

	// Wait for all at the gate; close releases them all.
	ready.Wait()
	close(start)
	done.Wait()
	fmt.Println("after barrier:", results)
}
