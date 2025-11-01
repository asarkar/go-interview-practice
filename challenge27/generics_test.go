package generics

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

// TestUtilityFunctions tests the generic utility functions
func TestUtilityFunctions(t *testing.T) {
	t.Run("Filter", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		evens := Filter(numbers, func(n int) bool {
			return n%2 == 0
		})
		if len(evens) != 4 {
			t.Errorf("Expected 4 even numbers, got %d", len(evens))
		}
		expected := []int{2, 4, 6, 8}
		if !reflect.DeepEqual(evens, expected) {
			t.Errorf("Expected %v, got %v", expected, evens)
		}

		// Test with empty slice
		empty := []int{}
		filtered := Filter(empty, func(n int) bool { return true })
		if len(filtered) != 0 {
			t.Errorf(
				"Expected filtering empty slice to return empty slice, got length %d",
				len(filtered),
			)
		}
	})

	t.Run("Map", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4}
		squares := Map(numbers, func(n int) int {
			return n * n
		})
		expected := []int{1, 4, 9, 16}
		if !reflect.DeepEqual(squares, expected) {
			t.Errorf("Expected %v, got %v", expected, squares)
		}

		// Test mapping to different type
		strings := Map(numbers, func(n int) string {
			return strconv.Itoa(n)
		})
		expectedStrings := []string{"1", "2", "3", "4"}
		if !reflect.DeepEqual(strings, expectedStrings) {
			t.Errorf("Expected %v, got %v", expectedStrings, strings)
		}

		// Test with empty slice
		empty := []int{}
		mapped := Map(empty, func(n int) int { return n * 2 })
		if len(mapped) != 0 {
			t.Errorf(
				"Expected mapping empty slice to return empty slice, got length %d",
				len(mapped),
			)
		}
	})

	t.Run("Reduce", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		sum := Reduce(numbers, 0, func(acc, n int) int {
			return acc + n
		})
		if sum != 15 {
			t.Errorf("Expected sum to be 15, got %d", sum)
		}

		// Test with different types
		product := Reduce(numbers, 1, func(acc, n int) int {
			return acc * n
		})
		if product != 120 {
			t.Errorf("Expected product to be 120, got %d", product)
		}

		// Test with empty slice
		empty := []int{}
		result := Reduce(empty, 42, func(acc, n int) int { return acc + n })
		if result != 42 {
			t.Errorf("Expected reducing empty slice to return initial value 42, got %d", result)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		if !Contains(numbers, 3) {
			t.Error("Expected numbers to contain 3")
		}
		if Contains(numbers, 6) {
			t.Error("Expected numbers to not contain 6")
		}

		// Test with empty slice
		empty := []int{}
		if Contains(empty, 1) {
			t.Error("Expected empty slice to not contain 1")
		}
	})

	t.Run("FindIndex", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		if FindIndex(numbers, 3) != 2 {
			t.Errorf("Expected index of 3 to be 2, got %d", FindIndex(numbers, 3))
		}
		if FindIndex(numbers, 6) != -1 {
			t.Errorf("Expected index of 6 to be -1, got %d", FindIndex(numbers, 6))
		}

		// Test with empty slice
		empty := []int{}
		if FindIndex(empty, 1) != -1 {
			t.Errorf("Expected index in empty slice to be -1, got %d", FindIndex(empty, 1))
		}
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		withDuplicates := []int{1, 2, 2, 3, 1, 4, 5, 5}
		unique := RemoveDuplicates(withDuplicates)
		// Since order can vary with map iteration, we'll sort the results
		sort.Ints(unique)
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(unique, expected) {
			t.Errorf("Expected %v, got %v", expected, unique)
		}

		// Test with no duplicates
		noDuplicates := []int{1, 2, 3, 4, 5}
		result := RemoveDuplicates(noDuplicates)
		sort.Ints(result)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with empty slice
		empty := []int{}
		emptyResult := RemoveDuplicates(empty)
		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}
	})
}
