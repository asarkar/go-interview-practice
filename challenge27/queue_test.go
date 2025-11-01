package generics

import "testing"

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
