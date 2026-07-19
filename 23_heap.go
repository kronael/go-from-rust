//go:build ignore

package main

import (
	"container/heap"
	"fmt"
)

type intHeap []int

func (numbers intHeap) Len() int           { return len(numbers) }
func (numbers intHeap) Less(i, j int) bool { return numbers[i] < numbers[j] }
func (numbers intHeap) Swap(i, j int)      { numbers[i], numbers[j] = numbers[j], numbers[i] }

func (numbers *intHeap) Push(value any) {
	*numbers = append(*numbers, value.(int))
}

func (numbers *intHeap) Pop() any {
	// heap.Pop first moves the root to the end. This adapter removes that last
	// slot; it does not search for the minimum itself.
	old := *numbers
	last := old[len(old)-1]
	*numbers = old[:len(old)-1]
	return last
}

func main() {
	// container/heap supplies the algorithm. The adapter supplies storage and
	// ordering; Less makes this a min-heap rather than Rust's max BinaryHeap.
	numbers := &intHeap{5, 1, 3}
	heap.Init(numbers)
	heap.Push(numbers, 2)

	fmt.Print("ascending:")
	for numbers.Len() > 0 {
		fmt.Printf(" %d", heap.Pop(numbers).(int))
	}
	fmt.Println()
}
