package generics

//
// 3. Generic Queue
//

type DLLNode[T any] struct {
	Value T
	next  *DLLNode[T]
	prev  *DLLNode[T]
}

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	head, tail *DLLNode[T]
	size       int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	node := &DLLNode[T]{Value: value}
	node.prev = q.tail
	if q.tail == nil {
		q.head = node
	} else {
		q.tail.next = node
	}

	q.tail = node
	q.size++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	value, err := q.Front()

	if err == nil {
		node := q.head
		if node.next == nil {
			q.head = nil
			q.tail = nil
		} else {
			node.next.prev = nil
			q.head = node.next
			node.next = nil
		}
		q.size--
	}
	return value, err
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.head == nil {
		return zero, ErrEmptyCollection
	}
	return q.head.Value, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return q.size
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}
