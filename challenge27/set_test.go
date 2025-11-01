package generics

import (
	"reflect"
	"sort"
	"testing"
)

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
