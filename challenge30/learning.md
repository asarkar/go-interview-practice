# Learning Guide: Go Context Package

## What is Context?

The `context` package in Go provides a way to carry **cancellation signals**, **timeouts**, and **request-scoped values** across API boundaries and goroutines. It's one of the most important packages in Go for building robust, production-ready applications.

### Why Context Matters

```go
// Without context - no way to cancel or timeout
func fetchData() ([]byte, error) {
    resp, err := http.Get("https://api.example.com/data")
    // What if this takes 5 minutes? No way to cancel!
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

// With context - cancellable and time-bounded
func fetchDataWithContext(ctx context.Context) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/data", nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err // Could be context.DeadlineExceeded or context.Canceled
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

## Core Context Types

### 1. Background Context
The root of all contexts - never cancelled, has no deadline, carries no values.

```go
ctx := context.Background()
// Use this as the top-level context in main(), tests, or initialization
```

### 2. TODO Context
A placeholder when you're unsure which context to use.

```go
ctx := context.TODO()
// Use this during development when context isn't clear yet
```

### 3. Cancellation Context
Can be manually cancelled to signal goroutines to stop work.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // Always call cancel to prevent memory leaks

go func() {
    select {
    case <-ctx.Done():
        fmt.Println("Work cancelled:", ctx.Err())
        return
    case <-time.After(5 * time.Second):
        fmt.Println("Work completed")
    }
}()

// Cancel after 2 seconds
time.Sleep(2 * time.Second)
cancel() // This triggers ctx.Done()
```

### 4. Deadline/Timeout Context
Automatically cancelled after a specific time.

```go
// WithDeadline - cancel at specific time
deadline := time.Now().Add(30 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// WithTimeout - cancel after duration
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### 5. Value Context
Carries request-scoped data across function calls.

```go
ctx := context.WithValue(context.Background(), "userID", "12345")
ctx = context.WithValue(ctx, "requestID", "req-abc-123")

// Retrieve values
userID := ctx.Value("userID").(string)
requestID := ctx.Value("requestID").(string)
```

## Context Patterns

### Pattern 1: Checking for Cancellation

```go
func doWork(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check for cancellation periodically
        select {
        case <-ctx.Done():
            return ctx.Err() // context.Canceled or context.DeadlineExceeded
        default:
            // Continue work
        }
        
        // Simulate work
        time.Sleep(10 * time.Millisecond)
        fmt.Printf("Processed item %d\n", i)
    }
    return nil
}
```

### Pattern 2: Context with Goroutines

```go
func processInParallel(ctx context.Context, items []string) error {
    errChan := make(chan error, len(items))
    
    for _, item := range items {
        go func(item string) {
            select {
            case <-ctx.Done():
                errChan <- ctx.Err()
                return
            case errChan <- processItem(item):
                return
            }
        }(item)
    }
    
    // Wait for all goroutines
    for i := 0; i < len(items); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    
    return nil
}
```

### Pattern 3: Context Racing

```go
func executeWithTimeout(ctx context.Context, task func() error) error {
    done := make(chan error, 1)
    
    go func() {
        done <- task()
    }()
    
    select {
    case err := <-done:
        return err // Task completed first
    case <-ctx.Done():
        return ctx.Err() // Context cancelled/timeout first
    }
}
```

## Real-World Examples

### Web Server with Request Timeouts

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Create context with timeout for this request
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    // Add request-specific values
    ctx = context.WithValue(ctx, "requestID", generateRequestID())
    ctx = context.WithValue(ctx, "userID", getUserID(r))
    
    // Process request with context
    result, err := processRequest(ctx, r)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusRequestTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

func processRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    // Extract values from context
    requestID := ctx.Value("requestID").(string)
    userID := ctx.Value("userID").(string)
    
    log.Printf("Processing request %s for user %s", requestID, userID)
    
    // Make database call with context
    data, err := fetchFromDatabase(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Make external API call with context
    enriched, err := enrichData(ctx, data)
    if err != nil {
        return nil, err
    }
    
    return enriched, nil
}
```

### Worker Pool with Graceful Shutdown

```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, 100),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    log.Printf("Worker %d started", id)
    defer log.Printf("Worker %d stopped", id)
    
    for {
        select {
        case <-wp.ctx.Done():
            log.Printf("Worker %d shutting down: %v", id, wp.ctx.Err())
            return
        case job := <-wp.jobs:
            wp.processJob(job)
        }
    }
}

func (wp *WorkerPool) processJob(job Job) {
    // Create context with timeout for this job
    ctx, cancel := context.WithTimeout(wp.ctx, job.Timeout)
    defer cancel()
    
    err := job.Execute(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            log.Printf("Job %s timed out", job.ID)
        } else {
            log.Printf("Job %s failed: %v", job.ID, err)
        }
    }
}

func (wp *WorkerPool) Shutdown() {
    wp.cancel() // This will cause all workers to stop
}
```

### Database Operations with Context

```go
func getUserOrders(ctx context.Context, db *sql.DB, userID string) ([]Order, error) {
    // Create query with context
    query := `
        SELECT id, user_id, product_name, amount, created_at 
        FROM orders 
        WHERE user_id = $1 
        ORDER BY created_at DESC
    `
    
    // Execute query with context (will be cancelled if context is cancelled)
    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    var orders []Order
    for rows.Next() {
        // Check for cancellation while processing rows
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        var order Order
        err := rows.Scan(&order.ID, &order.UserID, &order.ProductName, &order.Amount, &order.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        orders = append(orders, order)
    }
    
    return orders, nil
}
```

## Context Best Practices

### ✅ DO:

1. **Pass context as first parameter**
```go
func ProcessData(ctx context.Context, data []byte) error // ✅ Good
func ProcessData(data []byte, ctx context.Context) error // ❌ Bad
```

2. **Always call cancel() to prevent memory leaks**
```go
ctx, cancel := context.WithTimeout(parent, 30*time.Second)
defer cancel() // ✅ Always do this
```

3. **Check ctx.Done() in loops and long operations**
```go
for i := 0; i < len(items); i++ {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    processItem(items[i])
}
```

4. **Use context for request-scoped values**
```go
ctx = context.WithValue(ctx, "traceID", "abc123")
ctx = context.WithValue(ctx, "userID", "user456")
```

5. **Derive child contexts from parent contexts**
```go
childCtx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
```

### ❌ DON'T:

1. **Don't store contexts in structs** (with rare exceptions)
```go
// ❌ Bad - context stored in struct
type Server struct {
    ctx context.Context
}

// ✅ Good - context passed as parameter
func (s *Server) ProcessRequest(ctx context.Context) error
```

2. **Don't pass nil context**
```go
ProcessData(nil, data) // ❌ Bad
ProcessData(context.Background(), data) // ✅ Good
```

3. **Don't use context for optional parameters**
```go
// ❌ Bad - using context for config
ctx = context.WithValue(ctx, "retryCount", 3)

// ✅ Good - use struct for config
type Config struct {
    RetryCount int
}
func ProcessData(ctx context.Context, cfg Config) error
```

4. **Don't ignore context cancellation**
```go
// ❌ Bad - ignoring context
func doWork(ctx context.Context) {
    for i := 0; i < 1000; i++ {
        // No context checking
        time.Sleep(100 * time.Millisecond)
    }
}

// ✅ Good - respecting context
func doWork(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        time.Sleep(100 * time.Millisecond)
    }
    return nil
}
```

## Common Errors and Solutions

### Error 1: Memory Leaks from Not Calling Cancel

```go
// ❌ Memory leak - cancel not called
func badExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // cancel() never called - goroutine and timer leak!
    doWork(ctx)
}

// ✅ Fixed - always call cancel
func goodExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // Always call cancel
    doWork(ctx)
}
```

### Error 2: Race Conditions with Context Values

```go
// ❌ Race condition - value might change
func badExample(ctx context.Context) {
    go func() {
        userID := ctx.Value("userID").(string) // Might panic if nil
        processUser(userID)
    }()
}

// ✅ Safe value extraction
func goodExample(ctx context.Context) {
    userIDValue := ctx.Value("userID")
    if userIDValue == nil {
        return // Handle missing value
    }
    userID, ok := userIDValue.(string)
    if !ok {
        return // Handle wrong type
    }
    
    go func() {
        processUser(userID)
    }()
}
```

### Error 3: Context Inheritance Issues

```go
// ❌ Bad - creating independent contexts
func badChain() {
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second) // Independent!
    defer cancel2()
    
    doWork(ctx2) // Won't inherit ctx1's cancellation
}

// ✅ Good - proper context chaining
func goodChain() {
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    ctx2, cancel2 := context.WithTimeout(ctx1, 5*time.Second) // Inherits from ctx1
    defer cancel2()
    
    doWork(ctx2) // Will be cancelled when ctx1 OR ctx2 times out
}
```

## Testing with Context

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    err := doSlowWork(ctx)
    if err != context.DeadlineExceeded {
        t.Errorf("Expected timeout, got %v", err)
    }
}

func TestWithCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(50 * time.Millisecond)
        cancel()
    }()
    
    err := doWork(ctx)
    if err != context.Canceled {
        t.Errorf("Expected cancellation, got %v", err)
    }
}
```

## Advanced Context Patterns

### Custom Context Types (Advanced)

```go
type contextKey string

const (
    RequestIDKey contextKey = "requestID"
    UserIDKey   contextKey = "userID"
)

// Type-safe context helpers
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
    requestID, ok := ctx.Value(RequestIDKey).(string)
    return requestID, ok
}
```

### Context Middleware

```go
func contextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add request ID to context
        requestID := generateRequestID()
        ctx := WithRequestID(r.Context(), requestID)
        
        // Add user info to context
        if userID := getUserFromAuth(r); userID != "" {
            ctx = context.WithValue(ctx, UserIDKey, userID)
        }
        
        // Set response header
        w.Header().Set("X-Request-ID", requestID)
        
        // Call next handler with enriched context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Performance Considerations

1. **Context overhead is minimal** - Don't worry about performance for normal use
2. **Avoid excessive context chaining** - Each WithValue creates a new context
3. **Use context values sparingly** - They're not meant for large data
4. **Be careful with context in hot paths** - Profile if you suspect issues

## Resources for Further Learning

- [Go Context Package Documentation](https://pkg.go.dev/context)
- [Go Blog: Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Wiki: Context](https://github.com/golang/go/wiki/Context)

## Summary

The context package is essential for:
- **Cancellation**: Stop work when it's no longer needed
- **Timeouts**: Prevent operations from running too long
- **Request-scoped values**: Pass data across function boundaries
- **Graceful shutdowns**: Coordinate cleanup across goroutines

Master these patterns and you'll write more robust, maintainable Go applications! 