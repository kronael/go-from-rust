//go:build ignore

package main

import "fmt"

func main() {
	done := make(chan struct{})

	go func() {
		// Goroutines do not swallow panics. Without this
		// same-goroutine recover, the unrecovered panic
		// terminates the whole process. recover works only when
		// called directly by a deferred function in this
		// panicking goroutine.
		defer func() {
			recovered := recover()
			fmt.Println("worker recovered:", recovered)
			close(done)
		}()

		fmt.Println("worker panicking")
		// Expected failures use error; panic is for violated
		// assumptions.
		panic("broken invariant")
	}()

	<-done
	fmt.Println("main still running")

	// Like Rust catch_unwind, recovery is an exceptional
	// boundary, not normal error handling. Go has no global
	// panic handler. Wrap each goroutine if recovery is
	// required; re-panic after logging when the process should
	// still crash. Libraries may recover internally: net/http
	// does this around handlers.
}
