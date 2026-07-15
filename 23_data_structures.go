//go:build ignore

package main

import (
	"container/heap"
	"fmt"
)

type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *intHeap) Push(value any) {
	*h = append(*h, value.(int))
}

func (h *intHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]
	return last
}

func main() {
	// A slice is a Vec-like stack: append, then remove from the end.
	stack := []string{"a", "b"}
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	fmt.Println("stack:", top, stack)

	// Reslicing from the front is enough for a small, short-lived queue.
	queue := []int{1, 2, 3}
	front := queue[0]
	queue = queue[1:]
	fmt.Println("queue:", front, queue)

	// map[T]struct{} is the usual HashSet equivalent.
	set := map[int]struct{}{1: {}, 2: {}, 3: {}}
	_, hasTwo := set[2]
	fmt.Println("set:", len(set), hasTwo)

	// container/heap supplies the algorithm; intHeap supplies its ordering.
	numbers := &intHeap{5, 1, 3}
	heap.Init(numbers)
	heap.Push(numbers, 2)
	fmt.Print("min-heap:")
	for numbers.Len() > 0 {
		fmt.Printf(" %d", heap.Pop(numbers).(int))
	}
	fmt.Println()
}
