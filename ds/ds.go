// Package ds provides the small generic data structures Go's stdlib omits:
// a priority-queue Heap, a Stack, a Queue, and a Set. Dependency-free — copy
// what you need, or reach for a library (see the README's data-structures note).
//
// All four use value-type storage and zero freed slots so pointer elements can
// be garbage-collected. The empty-result methods return (T, ok) rather than
// panicking — Go's comma-ok idiom standing in for Rust's Option<T>.
package ds

import "cmp"

// ---- Heap: binary min-heap / priority queue ----

// Heap is a binary heap ordered by less. less(a, b)==true means a has higher
// priority (pops first). Rust analogue: BinaryHeap<T> (which is a max-heap).
type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// NewHeap builds a heap with a custom ordering (e.g. a > b for a max-heap).
func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

// NewOrderedHeap builds a min-heap for any ordered type (smallest pops first).
func NewOrderedHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{less: func(a, b T) bool { return a < b }}
}

func (h *Heap[T]) Len() int { return len(h.data) }

// Push adds x and sifts it up — O(log n).
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	i := len(h.data) - 1
	for i > 0 {
		p := (i - 1) / 2
		if !h.less(h.data[i], h.data[p]) {
			break
		}
		h.data[i], h.data[p] = h.data[p], h.data[i]
		i = p
	}
}

// Pop removes and returns the highest-priority element — O(log n).
// ok is false if the heap is empty.
func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	if len(h.data) == 0 {
		return zero, false
	}
	top := h.data[0]
	n := len(h.data) - 1
	h.data[0] = h.data[n]
	h.data[n] = zero // drop the moved slot's reference for GC
	h.data = h.data[:n]
	i := 0
	for {
		l, r, best := 2*i+1, 2*i+2, i
		if l < n && h.less(h.data[l], h.data[best]) {
			best = l
		}
		if r < n && h.less(h.data[r], h.data[best]) {
			best = r
		}
		if best == i {
			break
		}
		h.data[i], h.data[best] = h.data[best], h.data[i]
		i = best
	}
	return top, true
}

// Peek returns the highest-priority element without removing it — O(1).
func (h *Heap[T]) Peek() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}
	return h.data[0], true
}

// ---- Stack: LIFO ----

type Stack[T any] struct{ data []T }

func (s *Stack[T]) Push(x T) { s.data = append(s.data, x) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	n := len(s.data)
	if n == 0 {
		return zero, false
	}
	x := s.data[n-1]
	s.data[n-1] = zero
	s.data = s.data[:n-1]
	return x, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Len() int { return len(s.data) }

// ---- Queue: FIFO ----

// Queue is a slice-backed FIFO with amortized O(1) ops. It tracks a head index
// and compacts when the dead prefix grows — avoiding the "s = s[1:] pins the
// whole backing array forever" leak that a naive front-pop slice has.
type Queue[T any] struct {
	data []T
	head int
}

func (q *Queue[T]) Enqueue(x T) { q.data = append(q.data, x) }

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.head >= len(q.data) {
		return zero, false
	}
	x := q.data[q.head]
	q.data[q.head] = zero
	q.head++
	if q.head > len(q.data)/2 { // compact once the dead prefix dominates
		n := copy(q.data, q.data[q.head:])
		q.data = q.data[:n]
		q.head = 0
	}
	return x, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.head >= len(q.data) {
		var zero T
		return zero, false
	}
	return q.data[q.head], true
}

func (q *Queue[T]) Len() int { return len(q.data) - q.head }

// ---- Set: unordered set of comparable values ----

// Set is a map-backed set. The struct{} value costs zero bytes — the idiomatic
// Go set. Rust analogue: HashSet<T>.
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, x := range items {
		s[x] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(x T)      { s[x] = struct{}{} }
func (s Set[T]) Has(x T) bool { _, ok := s[x]; return ok }
func (s Set[T]) Delete(x T)   { delete(s, x) }
func (s Set[T]) Len() int     { return len(s) }
