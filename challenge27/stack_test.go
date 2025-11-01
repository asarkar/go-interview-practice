package generics

import "testing"

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
