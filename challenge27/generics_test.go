package generics

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

// TestPair tests the Pair implementation
func TestPair(t *testing.T) {
	t.Run("NewPair", func(t *testing.T) {
		p := NewPair("hello", 42)
		if p.First != "hello" {
			t.Errorf("Expected First to be 'hello', got %v", p.First)
		}
		if p.Second != 42 {
			t.Errorf("Expected Second to be 42, got %v", p.Second)
		}
	})

	t.Run("Swap", func(t *testing.T) {
		p := NewPair("hello", 42)
		swapped := p.Swap()
		if swapped.First != 42 {
			t.Errorf("Expected swapped First to be 42, got %v", swapped.First)
		}
		if swapped.Second != "hello" {
			t.Errorf("Expected swapped Second to be 'hello', got %v", swapped.Second)
		}
	})
}

// TestStack tests the Stack implementation
func TestStack(t *testing.T) {
	t.Run("NewStack", func(t *testing.T) {
		stack := NewStack[int]()
		if stack == nil {
			t.Error("Expected NewStack to return a non-nil stack")
		}
		if !stack.IsEmpty() {
			t.Error("Expected new stack to be empty")
		}
		if stack.Size() != 0 {
			t.Errorf("Expected size of new stack to be 0, got %d", stack.Size())
		}
	})

	t.Run("Push", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(1)
		if stack.IsEmpty() {
			t.Error("Expected stack to not be empty after Push")
		}
		if stack.Size() != 1 {
			t.Errorf("Expected size to be 1 after Push, got %d", stack.Size())
		}
	})

	t.Run("Peek", func(t *testing.T) {
		stack := NewStack[int]()
		_, err := stack.Peek()
		if err == nil {
			t.Error("Expected Peek on empty stack to return error")
		}

		stack.Push(1)
		stack.Push(2)
		val, err := stack.Peek()
		if err != nil {
			t.Errorf("Expected Peek on non-empty stack to not return error, got %v", err)
		}
		if val != 2 {
			t.Errorf("Expected Peek to return 2, got %v", val)
		}
		if stack.Size() != 2 {
			t.Errorf("Expected size to still be 2 after Peek, got %d", stack.Size())
		}
	})

	t.Run("Pop", func(t *testing.T) {
		stack := NewStack[int]()
		_, err := stack.Pop()
		if err == nil {
			t.Error("Expected Pop on empty stack to return error")
		}

		stack.Push(1)
		stack.Push(2)
		val, err := stack.Pop()
		if err != nil {
			t.Errorf("Expected Pop on non-empty stack to not return error, got %v", err)
		}
		if val != 2 {
			t.Errorf("Expected Pop to return 2, got %v", val)
		}
		if stack.Size() != 1 {
			t.Errorf("Expected size to be 1 after Pop, got %d", stack.Size())
		}

		val, err = stack.Pop()
		if err != nil {
			t.Errorf("Expected Pop on non-empty stack to not return error, got %v", err)
		}
		if val != 1 {
			t.Errorf("Expected Pop to return 1, got %v", val)
		}
		if !stack.IsEmpty() {
			t.Error("Expected stack to be empty after popping all elements")
		}
	})
}

// TestQueue tests the Queue implementation
func TestQueue(t *testing.T) {
	t.Run("NewQueue", func(t *testing.T) {
		queue := NewQueue[string]()
		if queue == nil {
			t.Error("Expected NewQueue to return a non-nil queue")
		}
		if !queue.IsEmpty() {
			t.Error("Expected new queue to be empty")
		}
		if queue.Size() != 0 {
			t.Errorf("Expected size of new queue to be 0, got %d", queue.Size())
		}
	})

	t.Run("Enqueue", func(t *testing.T) {
		queue := NewQueue[string]()
		queue.Enqueue("first")
		if queue.IsEmpty() {
			t.Error("Expected queue to not be empty after Enqueue")
		}
		if queue.Size() != 1 {
			t.Errorf("Expected size to be 1, got %d", queue.Size())
		}
	})

	t.Run("Front", func(t *testing.T) {
		queue := NewQueue[string]()
		_, err := queue.Front()
		if err == nil {
			t.Error("Expected Front on empty queue to return error")
		}

		queue.Enqueue("first")
		queue.Enqueue("second")
		val, err := queue.Front()
		if err != nil {
			t.Errorf("Expected Front on non-empty queue to not return error, got %v", err)
		}
		if val != "first" {
			t.Errorf("Expected Front to return 'first', got %v", val)
		}
		if queue.Size() != 2 {
			t.Errorf("Expected size to still be 2 after Front, got %d", queue.Size())
		}
	})

	t.Run("Dequeue", func(t *testing.T) {
		queue := NewQueue[string]()
		_, err := queue.Dequeue()
		if err == nil {
			t.Error("Expected Dequeue on empty queue to return error")
		}

		queue.Enqueue("first")
		queue.Enqueue("second")
		val, err := queue.Dequeue()
		if err != nil {
			t.Errorf("Expected Dequeue on non-empty queue to not return error, got %v", err)
		}
		if val != "first" {
			t.Errorf("Expected Dequeue to return 'first', got %v", val)
		}
		if queue.Size() != 1 {
			t.Errorf("Expected size to be 1 after Dequeue, got %d", queue.Size())
		}

		val, err = queue.Dequeue()
		if err != nil {
			t.Errorf("Expected Dequeue on non-empty queue to not return error, got %v", err)
		}
		if val != "second" {
			t.Errorf("Expected Dequeue to return 'second', got %v", val)
		}
		if !queue.IsEmpty() {
			t.Error("Expected queue to be empty after dequeuing all elements")
		}
	})
}

// TestSet tests the Set implementation
func TestSet(t *testing.T) {
	t.Run("NewSet", func(t *testing.T) {
		set := NewSet[int]()
		if set == nil {
			t.Error("Expected NewSet to return a non-nil set")
		}
		if set.Size() != 0 {
			t.Errorf("Expected size of new set to be 0, got %d", set.Size())
		}
	})

	t.Run("Add", func(t *testing.T) {
		set := NewSet[int]()
		set.Add(1)
		if !set.Contains(1) {
			t.Error("Expected set to contain 1 after Add")
		}
		if set.Size() != 1 {
			t.Errorf("Expected size to be 1 after Add, got %d", set.Size())
		}

		// Adding the same element again shouldn't change the set
		set.Add(1)
		if set.Size() != 1 {
			t.Errorf("Expected size to still be 1 after adding duplicate, got %d", set.Size())
		}

		set.Add(2)
		if !set.Contains(2) {
			t.Error("Expected set to contain 2 after Add")
		}
		if set.Size() != 2 {
			t.Errorf("Expected size to be 2 after adding second element, got %d", set.Size())
		}
	})

	t.Run("Remove", func(t *testing.T) {
		set := NewSet[int]()
		set.Add(1)
		set.Add(2)
		set.Add(3)

		set.Remove(2)
		if set.Contains(2) {
			t.Error("Expected set to not contain 2 after Remove")
		}
		if set.Size() != 2 {
			t.Errorf("Expected size to be 2 after Remove, got %d", set.Size())
		}

		// Removing a non-existent element shouldn't change the set
		set.Remove(4)
		if set.Size() != 2 {
			t.Errorf(
				"Expected size to still be 2 after removing non-existent element, got %d",
				set.Size(),
			)
		}
	})

	t.Run("Elements", func(t *testing.T) {
		set := NewSet[int]()
		elements := set.Elements()
		if len(elements) != 0 {
			t.Errorf("Expected empty set to have 0 elements, got %d", len(elements))
		}

		set.Add(1)
		set.Add(2)
		set.Add(3)
		elements = set.Elements()
		sort.Ints(elements) // Sort to make the test deterministic
		if len(elements) != 3 {
			t.Errorf("Expected set to have 3 elements, got %d", len(elements))
		}
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(elements, expected) {
			t.Errorf("Expected elements to be %v, got %v", expected, elements)
		}
	})

	t.Run("Union", func(t *testing.T) {
		set1 := NewSet[int]()
		set1.Add(1)
		set1.Add(2)
		set1.Add(3)

		set2 := NewSet[int]()
		set2.Add(3)
		set2.Add(4)
		set2.Add(5)

		union := Union(set1, set2)
		if union.Size() != 5 {
			t.Errorf("Expected union to have 5 elements, got %d", union.Size())
		}
		for i := 1; i <= 5; i++ {
			if !union.Contains(i) {
				t.Errorf("Expected union to contain %d", i)
			}
		}
	})

	t.Run("Intersection", func(t *testing.T) {
		set1 := NewSet[int]()
		set1.Add(1)
		set1.Add(2)
		set1.Add(3)

		set2 := NewSet[int]()
		set2.Add(3)
		set2.Add(4)
		set2.Add(5)

		intersection := Intersection(set1, set2)
		if intersection.Size() != 1 {
			t.Errorf("Expected intersection to have 1 element, got %d", intersection.Size())
		}
		if !intersection.Contains(3) {
			t.Error("Expected intersection to contain 3")
		}
	})

	t.Run("Difference", func(t *testing.T) {
		set1 := NewSet[int]()
		set1.Add(1)
		set1.Add(2)
		set1.Add(3)

		set2 := NewSet[int]()
		set2.Add(3)
		set2.Add(4)
		set2.Add(5)

		difference := Difference(set1, set2)
		if difference.Size() != 2 {
			t.Errorf("Expected difference to have 2 elements, got %d", difference.Size())
		}
		if !difference.Contains(1) || !difference.Contains(2) {
			t.Error("Expected difference to contain 1 and 2")
		}
		if difference.Contains(3) {
			t.Error("Expected difference to not contain 3")
		}
	})
}

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
