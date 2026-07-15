//go:build ignore

package main

import (
	"context"
	"fmt"
)

func work(ctx context.Context, done chan<- struct{}) {
	<-ctx.Done()
	fmt.Println("worker:", ctx.Err())
	close(done)
}

func main() {
	// Go passes cancellation explicitly, conventionally as the first argument.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go work(ctx, done)

	cancel()
	<-done
	fmt.Println("main: stopped worker")
}
