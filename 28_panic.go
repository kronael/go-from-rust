//go:build ignore

package main

import "fmt"

func main() {
	done := make(chan struct{})

	go func() {
		// recover works only in a deferred func, same goroutine.
		defer func() {
			recovered := recover()
			fmt.Println("worker recovered:", recovered)
			close(done)
		}()

		fmt.Println("worker panicking")
		// panic is for broken invariants, not expected failures.
		panic("broken invariant")
	}()

	<-done
	fmt.Println("main still running")
}
