//go:build ignore

// datastructures.go — using the ds package: Heap (priority queue), Stack, Queue, Set.
// Go's stdlib has no generic versions of these, so you ship your own (see ds/ds.go)
// or pull a library. Every op returns (value, ok) instead of panicking on empty.
package main

import (
	"fmt"

	"gofrs/ds"
)

func main() {
	// priority queue: min-heap of ints (smallest pops first)
	h := ds.NewOrderedHeap[int]()
	for _, x := range []int{5, 1, 3} {
		h.Push(x)
	}
	if top, ok := h.Peek(); ok {
		fmt.Println("heap peek:", top)
	}
	fmt.Print("heap drain:")
	for h.Len() > 0 {
		v, _ := h.Pop()
		fmt.Printf(" %d", v)
	}
	fmt.Println()

	// custom priority: max-heap via a reversed less func
	mh := ds.NewHeap[int](func(a, b int) bool { return a > b })
	mh.Push(2)
	mh.Push(9)
	mh.Push(4)
	top, _ := mh.Pop()
	fmt.Println("max-heap top:", top)

	// stack (LIFO)
	st := &ds.Stack[string]{}
	st.Push("a")
	st.Push("b")
	v, _ := st.Pop()
	fmt.Println("stack pop:", v)

	// queue (FIFO)
	q := &ds.Queue[int]{}
	q.Enqueue(1)
	q.Enqueue(2)
	d, _ := q.Dequeue()
	fmt.Println("queue dequeue:", d, "remaining:", q.Len())

	// set
	s := ds.NewSet(1, 2, 2, 3) // duplicate 2 collapses
	fmt.Println("set len:", s.Len(), "has 2:", s.Has(2))
}
