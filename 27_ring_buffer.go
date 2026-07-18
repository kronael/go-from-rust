//go:build ignore

package main

import "fmt"

type ringDeque[T any] struct {
	values []T
	head   int
	size   int
}

func (deque *ringDeque[T]) pushBack(value T) bool {
	if deque.size == len(deque.values) {
		return false
	}
	index := (deque.head + deque.size) % len(deque.values)
	deque.values[index] = value
	deque.size++
	return true
}

func (deque *ringDeque[T]) pushFront(value T) bool {
	if deque.size == len(deque.values) {
		return false
	}
	deque.head = (deque.head - 1 + len(deque.values)) % len(deque.values)
	deque.values[deque.head] = value
	deque.size++
	return true
}

func (deque *ringDeque[T]) popFront() (T, bool) {
	var zero T
	if deque.size == 0 {
		return zero, false
	}
	value := deque.values[deque.head]
	deque.values[deque.head] = zero
	deque.head = (deque.head + 1) % len(deque.values)
	deque.size--
	return value, true
}

func (deque *ringDeque[T]) popBack() (T, bool) {
	var zero T
	if deque.size == 0 {
		return zero, false
	}
	index := (deque.head + deque.size - 1) % len(deque.values)
	value := deque.values[index]
	deque.values[index] = zero
	deque.size--
	return value, true
}

func main() {
	// Unlike a channel, this bounded ring is storage only: it never blocks and
	// is not goroutine-safe. It stays contiguous and reuses one allocation.
	deque := ringDeque[string]{values: make([]string, 3)}
	fmt.Println("push:", deque.pushBack("middle"), deque.pushFront("front"), deque.pushBack("back"))
	fmt.Println("full:", deque.pushBack("overflow"))

	front, frontOK := deque.popFront()
	back, backOK := deque.popBack()
	fmt.Println("ends:", front, frontOK, back, backOK)

	middle, middleOK := deque.popFront()
	_, emptyOK := deque.popFront()
	fmt.Println("drained:", middle, middleOK, "empty:", emptyOK)

	// Go has no generic VecDeque. container/list allocates nodes; an unbounded
	// ring needs growth and relinearization. Prefer a maintained implementation.
}
