package generics

import "slices"

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	values []T
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	if !s.Contains(value) {
		s.values = append(s.values, value)
	}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	s.values = slices.DeleteFunc(s.values, func(v T) bool {
		return value == v
	})
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	return slices.Contains(s.values, value)
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.values)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	out := make([]T, len(s.values))
	copy(out, s.values)
	return out
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	out := make([]T, len(s1.values))
	copy(out, s1.values)
	res := Set[T]{out}
	for _, v := range s2.values {
		res.Add(v)
	}
	return &res
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for i := range s1.values {
		if s2.Contains(s1.values[i]) {
			res.Add(s1.values[i])
		}
	}
	return res
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for i := range s1.values {
		if !s2.Contains(s1.values[i]) {
			res.Add(s1.values[i])
		}
	}
	return res
}
