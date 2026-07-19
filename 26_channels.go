//go:build ignore

package main

import "fmt"

type addition struct {
	amount int
	result chan<- int
}

func main() {
	// A channel resembles Rust mpsc: send and receive copy a value. chan<- int is
	// send-only; <-ch receives. Unbuffered operations block until a peer is ready.
	// Routing every update through one goroutine gives it sole ownership of count.
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

	// Capacity 2 stores these sends in FIFO order, so neither needs a waiting
	// receiver. The sender closes; receivers drain values and observe completion.
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

	// Sending copies a value; sending a slice still copies only its header.
	// Channels add synchronization and blocking. The next lesson's ring buffer
	// is storage only and needs external synchronization when shared.
}
