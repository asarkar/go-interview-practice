package generics

import "testing"

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
