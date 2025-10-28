package generics

import "errors"

// ErrEmptyCollection is returned when an operation cannot be performed on an empty collection
var ErrEmptyCollection = errors.New("collection is empty")

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	res := make([]T, 0)
	return Reduce(slice, res, func(acc []T, element T) []T {
		if predicate(element) {
			return append(acc, element)
		}
		return acc
	})
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	res := make([]U, 0, len(slice))
	return Reduce(slice, res, func(acc []U, element T) []U {
		return append(acc, mapper(element))
	})
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	res := initial
	for i := range slice {
		res = reducer(res, slice[i])
	}
	return res
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	return FindIndex(slice, element) >= 0
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for i := range slice {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	elements := make(map[T]struct{})
	res := make([]T, 0)
	for i := range slice {
		if _, exists := elements[slice[i]]; !exists {
			elements[slice[i]] = struct{}{}
			res = append(res, slice[i])
		}
	}
	return res
}
