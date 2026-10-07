// Package gopq provides a generic binary heap ordered by a caller's
// comparison function. It has no external dependencies and requires Go 1.18+.
package gopq

// Queue holds values of type T in priority order.
//
// Create a Queue with New. Do not copy a Queue after use. Queues are not safe for
// concurrent use without external synchronization. Values with equal priority
// may be popped in any order.
type Queue[T any] struct {
	items []T
	less  func(a, b T) bool
}

// New creates a queue where less(a, b) reports whether a should be popped before
// b. Use a < b for a min-heap or a > b for a max-heap of ordered values.
//
// less must define a strict weak ordering: it must be consistent, transitive,
// and false for equal values. Do not change priorities or comparison behavior
// while values are in the queue.
//
// Optional initial values are copied into a new slice, then heapified in O(n)
// time. The copy is shallow: pointed-to data remains shared with the caller.
// New panics if less is nil.
func New[T any](less func(a, b T) bool, values ...T) *Queue[T] {
	if less == nil {
		panic("gopq: nil comparison function")
	}
	items := append([]T(nil), values...)
	q := &Queue[T]{items: items, less: less}
	for i := len(items)/2 - 1; i >= 0; i-- {
		q.down(i)
	}
	return q
}

// Len returns the number of values in the queue in O(1) time.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Peek returns the highest-priority value without removing it in O(1) time.
// If the queue is empty, it returns the zero value of T and false.
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Push adds a value in O(log n) heap operations, with amortized slice growth.
// Push panics if the queue was not created with New.
func (q *Queue[T]) Push(value T) {
	if q.less == nil {
		panic("gopq: create the queue with New before pushing")
	}
	q.items = append(q.items, value)
	for i := len(q.items) - 1; i > 0; {
		parent := (i - 1) / 2
		if !q.less(q.items[i], q.items[parent]) {
			break
		}
		q.items[i], q.items[parent] = q.items[parent], q.items[i]
		i = parent
	}
}

// Pop removes and returns the highest-priority value in O(log n) time.
// If the queue is empty, it returns the zero value of T and false.
func (q *Queue[T]) Pop() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	value := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	// Release references held by the unused portion of the backing slice.
	var zero T
	q.items[last] = zero
	q.items = q.items[:last]
	if last > 0 {
		q.down(0)
	}
	return value, true
}

// down restores the heap invariant below i: no child precedes its parent.
func (q *Queue[T]) down(i int) {
	// Only internal nodes have children. This also avoids index overflow.
	for i < len(q.items)/2 {
		child := 2*i + 1
		right := child + 1
		if right < len(q.items) && q.less(q.items[right], q.items[child]) {
			child = right
		}
		if !q.less(q.items[child], q.items[i]) {
			return
		}
		q.items[i], q.items[child] = q.items[child], q.items[i]
		i = child
	}
}
