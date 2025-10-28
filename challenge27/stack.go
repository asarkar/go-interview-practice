package generics

//
// 2. Generic Stack
//

type ListNode[T any] struct {
	Value T
	next  *ListNode[T]
}

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	head *ListNode[T]
	size int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	node := &ListNode[T]{Value: value}
	node.next = s.head
	s.head = node
	s.size++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	value, err := s.Peek()
	if err == nil {
		s.head = s.head.next
		s.size--
	}

	return value, err
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if s.head == nil {
		return zero, ErrEmptyCollection
	}
	return s.head.Value, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return s.size
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return s.size == 0
}
