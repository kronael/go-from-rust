//go:build ignore

package main

import "fmt"

type addition struct {
	amount int
	result chan<- int
}

func main() {
	// Unbuffered channel: send/receive block until paired.
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

	// Buffered (cap 2): FIFO, sends don't block until full.
	values := make(chan int, 2)
	values <- 10
	values <- 20
	close(values)

	// range drains buffered values, stops on close.
	for value := range values {
		fmt.Println("value:", value)
	}

	// Closed receive returns the zero value and false.
	value, ok := <-values
	fmt.Println("closed receive:", value, ok)
}
