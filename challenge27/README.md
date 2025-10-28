# Challenge 27: Go Generics Data Structures

## Problem Statement

In this challenge, you will implement a set of generic data structures and algorithms in Go. This will allow you to practice using Go's generics feature (introduced in Go 1.18) to create reusable, type-safe code.

Your task is to implement several generic data structures that can work with any appropriate types:

1. A generic `Pair[T, U]` type that can hold two values of different types
2. A generic `Stack[T]` data structure with standard stack operations
3. A generic `Queue[T]` data structure with standard queue operations
4. A generic `Set[T]` data structure with basic set operations
5. A collection of generic utility functions for working with slices

These implementations will demonstrate how to use Go's generics to create flexible and type-safe code.

## Function Signatures

### 1. Generic Pair

```go
// Pair represents a generic pair of values of potentially different types
type Pair[T, U any] struct {
    First  T
    Second U
}

// NewPair creates a new pair with the given values
func NewPair[T, U any](first T, second U) Pair[T, U]

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T]
```

### 2. Generic Stack

```go
// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
    // Private implementation details
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T]

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T)

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error)

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error)

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool
```

### 3. Generic Queue

```go
// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
    // Private implementation details
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T]

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T)

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error)

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error)

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool
```

### 4. Generic Set

```go
// Set is a generic collection of unique elements
type Set[T comparable] struct {
    // Private implementation details
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T]

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T)

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T)

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool

// Size returns the number of elements in the set
func (s *Set[T]) Size() int

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T]

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T]

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T]
```

### 5. Generic Utility Functions

```go
// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T
```

## Input/Output Examples

### Pair Example
```go
// Creating a pair
p := NewPair("answer", 42)
fmt.Println(p.First)  // "answer"
fmt.Println(p.Second) // 42

// Swapping
swapped := p.Swap()
fmt.Println(swapped.First)  // 42
fmt.Println(swapped.Second) // "answer"
```

### Stack Example
```go
// Creating a stack of integers
stack := NewStack[int]()
stack.Push(1)
stack.Push(2)
stack.Push(3)

val, err := stack.Peek()
// val == 3, err == nil

val, err = stack.Pop()
// val == 3, err == nil

val, err = stack.Pop()
// val == 2, err == nil

size := stack.Size()
// size == 1

isEmpty := stack.IsEmpty()
// isEmpty == false
```

### Queue Example
```go
// Creating a queue of strings
queue := NewQueue[string]()
queue.Enqueue("first")
queue.Enqueue("second")
queue.Enqueue("third")

val, err := queue.Front()
// val == "first", err == nil

val, err = queue.Dequeue()
// val == "first", err == nil

val, err = queue.Dequeue()
// val == "second", err == nil

size := queue.Size()
// size == 1

isEmpty := queue.IsEmpty()
// isEmpty == false
```

### Set Example
```go
// Creating sets of integers
set1 := NewSet[int]()
set1.Add(1)
set1.Add(2)
set1.Add(3)
set1.Add(2) // Duplicate, won't be added

set2 := NewSet[int]()
set2.Add(2)
set2.Add(3)
set2.Add(4)

contains := set1.Contains(2)
// contains == true

union := Union(set1, set2)
// union contains 1, 2, 3, 4

intersection := Intersection(set1, set2)
// intersection contains 2, 3

difference := Difference(set1, set2)
// difference contains 1
```

### Utility Functions Example
```go
// Filter
numbers := []int{1, 2, 3, 4, 5, 6}
evens := Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// evens == [2, 4, 6]

// Map
squares := Map(numbers, func(n int) int {
    return n * n
})
// squares == [1, 4, 9, 16, 25, 36]

// Reduce
sum := Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
})
// sum == 21

// Contains
hasThree := Contains(numbers, 3)
// hasThree == true

// FindIndex
index := FindIndex(numbers, 4)
// index == 3

// RemoveDuplicates
withDuplicates := []int{1, 2, 2, 3, 1, 4, 5, 5}
uniques := RemoveDuplicates(withDuplicates)
// uniques == [1, 2, 3, 4, 5]
```

## Constraints

- Your implementation must make proper use of Go's generics features
- All functions must have appropriate type constraints
- The implementation should be efficient in terms of time and space complexity
- You must handle edge cases (empty collections, etc.) appropriately
- Your code should compile with Go 1.18 or later

## Evaluation Criteria

- Correctness: Does your solution correctly implement all the required data structures and functions?
- Proper use of generics: Are you using generics effectively with appropriate type constraints?
- Code quality: Is your code well-structured, readable, and maintainable?
- Performance: Are your implementations efficient in terms of time and space complexity?
- Error handling: Does your code handle error cases gracefully?

## Learning Resources

See the [learning.md](learning.md) document for a comprehensive guide on using generics in Go.

## Hints

1. Remember that the `comparable` constraint is required for types that will be used as map keys or compared with `==` and `!=`
2. For operations that may fail (like popping from an empty stack), return both the zero value and an error
3. The `any` constraint means any type is allowed, but more specific constraints can provide better type safety
4. Type parameters can be inferred in many cases, but explicit type arguments improve clarity in some situations
5. Go's generics implementation is designed to be straightforward and performant - focus on clean, idiomatic code rather than complex type manipulations 