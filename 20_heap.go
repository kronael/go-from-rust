//go:build ignore

package main

import (
	"container/heap"
	"fmt"
)

type intHeap []int

func (numbers intHeap) Len() int { return len(numbers) }
func (numbers intHeap) Less(i, j int) bool {
	return numbers[i] < numbers[j]
}
func (numbers intHeap) Swap(i, j int) {
	numbers[i], numbers[j] = numbers[j], numbers[i]
}

func (numbers *intHeap) Push(value any) {
	*numbers = append(*numbers, value.(int))
}

func (numbers *intHeap) Pop() any {
	// heap.Pop moves root to the end; Pop drops that slot.
	old := *numbers
	last := old[len(old)-1]
	*numbers = old[:len(old)-1]
	return last
}

func main() {
	// heap supplies the algorithm; the type supplies storage.
	numbers := &intHeap{5, 1, 3}
	heap.Init(numbers)
	heap.Push(numbers, 2)

	fmt.Print("ascending:")
	for numbers.Len() > 0 {
		fmt.Printf(" %d", heap.Pop(numbers).(int))
	}
	fmt.Println()
}
