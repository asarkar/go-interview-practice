 # Challenge 30: Context Management Implementation

## Overview

Implement a context manager that demonstrates essential Go `context` package patterns. The `context` package is fundamental for managing cancellation signals, timeouts, and request-scoped values in Go applications.

## Your Task

Implement a `ContextManager` interface with **6 core methods** and **2 helper functions**:

### ContextManager Interface

```go
type ContextManager interface {
    // Create a cancellable context from a parent context
    CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)
    
    // Create a context with timeout
    CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
    
    // Add a value to context
    AddValue(parent context.Context, key, value interface{}) context.Context
    
    // Get a value from context
    GetValue(ctx context.Context, key interface{}) (interface{}, bool)
    
    // Execute a task with context cancellation support
    ExecuteWithContext(ctx context.Context, task func() error) error
    
    // Wait for a duration or until context is cancelled
    WaitForCompletion(ctx context.Context, duration time.Duration) error
}
```

### Helper Functions

```go
// Simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error

// Process multiple items with context awareness
func ProcessItems(ctx context.Context, items []string) ([]string, error)
```

## Requirements

### Core Functionality
1. **Context Cancellation**: Handle manual cancellation via `context.WithCancel`
2. **Context Timeouts**: Implement timeout behavior via `context.WithTimeout`
3. **Value Storage**: Store and retrieve values via `context.WithValue`
4. **Task Execution**: Execute functions with cancellation support
5. **Wait Operations**: Wait for durations while respecting cancellation

### Implementation Details

- Use Go's standard `context` package functions
- Handle both `context.Canceled` and `context.DeadlineExceeded` errors
- Return appropriate boolean flags for value existence
- Support goroutine-based task execution with proper synchronization
- Process items in batches with cancellation checks between items

## Test Coverage

Your implementation will be tested with **13 test cases** covering:

- Context creation and cancellation
- Timeout behavior
- Value storage and retrieval
- Task execution scenarios (success, error, cancellation)
- Waiting operations (completion and cancellation)
- Helper function behavior
- Integration scenarios


## Getting Started

1. Examine the solution template and test file
2. Start with simple methods like `AddValue` and `GetValue`
3. Progress to cancellation and timeout contexts
4. Implement task execution with proper goroutine handling
5. Run tests frequently with `go test -v`

**Tip**: Check the `learning.md` file for comprehensive context patterns and examples!