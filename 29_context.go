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
	// context carries cancellation, passed as the first arg.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go work(ctx, done)

	// cancel closes ctx.Done; ctx.Err reports why.
	cancel()
	<-done
	fmt.Println("main: stopped worker")
}
