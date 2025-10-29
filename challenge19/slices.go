package slices

import "math"

// FindMax returns the maximum value in a slice of integers.
// If the slice is empty, it returns 0.
func FindMax(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	maxVal := math.MinInt32
	for i := range numbers {
		maxVal = max(maxVal, numbers[i])
	}
	return maxVal
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	seen := make(map[int]struct{})
	res := make([]int, 0)
	for i := range numbers {
		if _, exists := seen[numbers[i]]; !exists {
			res = append(res, numbers[i])
			seen[numbers[i]] = struct{}{}
		}
	}
	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	n := len(slice)
	res := make([]int, n)
	for i := range slice {
		res[n-i-1] = slice[i]
	}
	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	res := make([]int, 0)
	for i := range numbers {
		if numbers[i]%2 == 0 {
			res = append(res, numbers[i])
		}
	}
	return res
}
