//go:build ignore

package main

import "fmt"

type addition struct {
	amount int
	result chan<- int
}

func main() {
	// A channel can route every update through one goroutine. That goroutine
	// owns count, so count needs no mutex.
	additions := make(chan addition)
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		count := 0
		for addition := range additions {
			count += addition.amount
			addition.result <- count
		}
	}()

	results := make(chan int)
	additions <- addition{amount: 2, result: results}
	fmt.Println("owned by one goroutine:", <-results)
	close(additions)
	<-stopped

	// Sending copies a value. A slice copy still points at the same backing
	// array, so ownership transfer is only a convention; Go does not enforce it.
	shared := []int{10}
	handoff := make(chan []int, 1)
	handoff <- shared
	received := <-handoff
	received[0] = 20
	fmt.Println("slice storage still shared:", shared[0])

	values := make(chan int, 2)
	values <- 10
	values <- 20
	close(values)

	// range drains buffered values and then stops when the channel is closed.
	for value := range values {
		fmt.Println("value:", value)
	}

	// Rust receivers report disconnection. Go returns the zero value and false.
	value, ok := <-values
	fmt.Println("closed receive:", value, ok)

	// Channels coordinate and may block or schedule goroutines. A Mutex is
	// usually cheaper for a short critical section. With one goroutine touching
	// state, use neither; choose channels when the ownership boundary helps.
	// One producer plus one consumer is still concurrent: a channel is the
	// simple queue. Consider a specialized SPSC ring only after profiling.
}
