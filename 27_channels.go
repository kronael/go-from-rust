//go:build ignore

package main

import "fmt"

func main() {
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
}
