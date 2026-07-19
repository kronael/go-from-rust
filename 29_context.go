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
	// Go passes cancellation explicitly, conventionally as the first argument,
	// like passing a cancellation token; Rust has no direct std equivalent.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go work(ctx, done)

	// cancel closes ctx.Done. The worker wakes and ctx.Err reports
	// context.Canceled, then main waits for its done notification.
	cancel()
	<-done
	fmt.Println("main: stopped worker")
}
