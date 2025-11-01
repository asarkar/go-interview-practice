package ctx

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Test Basic Context Manager Functionality

func TestCreateCancellableContext(t *testing.T) {
	cm := NewContextManager()
	parent := context.Background()

	ctx, cancel := cm.CreateCancellableContext(parent)
	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}
	if cancel == nil {
		t.Fatal("Expected non-nil cancel function")
	}

	// Context should not be cancelled initially
	select {
	case <-ctx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
		// Good
	}

	// Cancel the context
	cancel()

	// Context should be cancelled now
	select {
	case <-ctx.Done():
		// Good
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be cancelled after calling cancel")
	}

	if ctx.Err() != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", ctx.Err())
	}
}

func TestCreateTimeoutContext(t *testing.T) {
	cm := NewContextManager()
	parent := context.Background()
	timeout := 50 * time.Millisecond

	ctx, cancel := cm.CreateTimeoutContext(parent, timeout)
	defer cancel()

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	// Wait for timeout
	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected context.DeadlineExceeded, got %v", ctx.Err())
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should timeout after specified duration")
	}
}

func TestAddAndGetValue(t *testing.T) {
	cm := NewContextManager()
	parent := context.Background()

	// Add values to context
	ctx := cm.AddValue(parent, "user", "alice")
	ctx = cm.AddValue(ctx, "requestID", "12345")
	ctx = cm.AddValue(ctx, "count", 42)

	// Test getting existing values
	value, exists := cm.GetValue(ctx, "user")
	if !exists {
		t.Fatal("Expected user value to exist")
	}
	if value != "alice" {
		t.Errorf("Expected 'alice', got %v", value)
	}

	value, exists = cm.GetValue(ctx, "requestID")
	if !exists {
		t.Fatal("Expected requestID value to exist")
	}
	if value != "12345" {
		t.Errorf("Expected '12345', got %v", value)
	}

	value, exists = cm.GetValue(ctx, "count")
	if !exists {
		t.Fatal("Expected count value to exist")
	}
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Test getting non-existent value
	value, exists = cm.GetValue(ctx, "nonexistent")
	if exists {
		t.Error("Expected nonexistent value to not exist")
	}
	if value != nil {
		t.Errorf("Expected nil value for nonexistent key, got %v", value)
	}
}

func TestExecuteWithContext_Success(t *testing.T) {
	cm := NewContextManager()
	ctx := context.Background()

	// Task that completes successfully
	task := func() error {
		time.Sleep(20 * time.Millisecond)
		return nil
	}

	err := cm.ExecuteWithContext(ctx, task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExecuteWithContext_TaskError(t *testing.T) {
	cm := NewContextManager()
	ctx := context.Background()

	expectedErr := errors.New("task failed")
	task := func() error {
		return expectedErr
	}

	err := cm.ExecuteWithContext(ctx, task)
	if err != expectedErr {
		t.Errorf("Expected task error %v, got %v", expectedErr, err)
	}
}

func TestExecuteWithContext_Cancellation(t *testing.T) {
	cm := NewContextManager()
	ctx, cancel := context.WithCancel(context.Background())

	// Long running task
	task := func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	}

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := cm.ExecuteWithContext(ctx, task)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestWaitForCompletion_Success(t *testing.T) {
	cm := NewContextManager()
	ctx := context.Background()
	duration := 50 * time.Millisecond

	start := time.Now()
	err := cm.WaitForCompletion(ctx, duration)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if elapsed < duration {
		t.Errorf("Expected to wait at least %v, waited %v", duration, elapsed)
	}
}

func TestWaitForCompletion_Cancellation(t *testing.T) {
	cm := NewContextManager()
	ctx, cancel := context.WithCancel(context.Background())
	duration := 200 * time.Millisecond

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := cm.WaitForCompletion(ctx, duration)
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	if elapsed >= duration {
		t.Errorf("Expected to be cancelled before %v, waited %v", duration, elapsed)
	}
}

func TestSimulateWork_Success(t *testing.T) {
	ctx := context.Background()
	duration := 30 * time.Millisecond

	start := time.Now()
	err := SimulateWork(ctx, duration, "test work")
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if elapsed < duration {
		t.Errorf("Expected work to take at least %v, took %v", duration, elapsed)
	}
}

func TestSimulateWork_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	duration := 200 * time.Millisecond

	// Cancel after short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := SimulateWork(ctx, duration, "test work")
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	if elapsed >= duration {
		t.Errorf("Expected work to be cancelled before %v, took %v", duration, elapsed)
	}
}

func TestProcessItems_Success(t *testing.T) {
	ctx := context.Background()
	items := []string{"item1", "item2", "item3"}

	results, err := ProcessItems(ctx, items)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	for i, result := range results {
		expected := "processed_" + items[i]
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	}
}

func TestProcessItems_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	// Cancel after processing a couple items
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	results, err := ProcessItems(ctx, items)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Should have partial results
	if len(results) == 0 {
		t.Error("Expected some results before cancellation")
	}

	if len(results) >= len(items) {
		t.Error("Expected cancellation to prevent all items from being processed")
	}
}

// Integration test combining multiple features
func TestContextIntegration(t *testing.T) {
	cm := NewContextManager()

	// Create context with timeout and values
	ctx, cancel := cm.CreateTimeoutContext(context.Background(), 200*time.Millisecond)
	defer cancel()

	ctx = cm.AddValue(ctx, "user", "bob")
	ctx = cm.AddValue(ctx, "session", "abc123")

	// Verify values are accessible
	user, exists := cm.GetValue(ctx, "user")
	if !exists || user != "bob" {
		t.Error("Expected user value to be accessible")
	}

	session, exists := cm.GetValue(ctx, "session")
	if !exists || session != "abc123" {
		t.Error("Expected session value to be accessible")
	}

	// Execute quick task (should succeed)
	quickTask := func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	err := cm.ExecuteWithContext(ctx, quickTask)
	if err != nil {
		t.Errorf("Quick task should succeed, got %v", err)
	}

	// Execute slow task (should timeout)
	slowTask := func() error {
		time.Sleep(300 * time.Millisecond)
		return nil
	}

	err = cm.ExecuteWithContext(ctx, slowTask)
	if err != context.DeadlineExceeded {
		t.Errorf("Slow task should timeout, got %v", err)
	}
}
