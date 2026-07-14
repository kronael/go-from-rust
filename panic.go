//go:build ignore

package main

import "fmt"

func main() {
	done := make(chan struct{})

	go func() {
		// recover only catches a panic in the same goroutine.
		defer func() {
			recovered := recover()
			fmt.Println("worker recovered:", recovered)
			close(done)
		}()

		fmt.Println("worker panicking")
		// Expected failures use error; panic is for violated assumptions.
		panic("broken invariant")
	}()

	<-done
	fmt.Println("main still running")
}
