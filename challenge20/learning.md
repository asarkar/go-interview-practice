# Learning Materials for Circuit Breaker Pattern

## Introduction to Circuit Breaker Pattern

The Circuit Breaker pattern is a design pattern used in software development to detect failures and encapsulate the logic of preventing a failure from constantly recurring. It's particularly useful in distributed systems where services may become temporarily unavailable.

Named after electrical circuit breakers that protect electrical circuits from damage, the software circuit breaker pattern prevents an application from repeatedly trying to execute operations that are likely to fail.

## The Problem It Solves

In distributed systems, services often depend on external resources:
- External APIs
- Databases
- File systems
- Network services

When these resources become unavailable or slow, your application can:
- Keep retrying and waste resources
- Create cascading failures
- Degrade user experience
- Overwhelm already struggling services

## How Circuit Breaker Works

### Three States

1. **Closed State** (Normal Operation)
   - Requests pass through to the service
   - Failures are monitored and counted
   - If failure threshold is reached, circuit trips to Open

2. **Open State** (Failing Fast)
   - All requests fail immediately without calling the service
   - Prevents wasting resources on operations likely to fail
   - After a timeout period, circuit moves to Half-Open

3. **Half-Open State** (Testing Recovery)
   - Limited number of requests are allowed through
   - If requests succeed, circuit closes
   - If requests fail, circuit opens again

### State Transition Diagram

```
    [Closed] --failure threshold--> [Open]
        ^                             |
        |                             |
    success                      timeout elapsed
        |                             |
        v                             v
    [Half-Open] --failure--> [Open]
```

## Go Implementation Concepts

### 1. Thread Safety

Circuit breakers must be thread-safe since they're typically shared across goroutines:

```go
type CircuitBreaker struct {
    state   State
    metrics Metrics
    mutex   sync.RWMutex  // Protects shared state
}

func (cb *CircuitBreaker) GetState() State {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    return cb.state
}
```

### 2. Metrics Collection

Track essential metrics for decision making:

```go
type Metrics struct {
    Requests            int64
    Successes           int64
    Failures            int64
    ConsecutiveFailures int64
    LastFailureTime     time.Time
}
```

### 3. Configurable Behavior

Make the circuit breaker configurable for different use cases:

```go
type Config struct {
    MaxRequests uint32                      // Half-open request limit
    Interval    time.Duration               // Metrics window
    Timeout     time.Duration               // Open -> Half-open timeout
    ReadyToTrip func(Metrics) bool          // Custom failure condition
    OnStateChange func(string, State, State) // State change callback
}
```

### 4. Context Support

Respect Go's context cancellation:

```go
func (cb *CircuitBreaker) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
    // Check context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Circuit breaker logic...
}
```

## Common Implementation Patterns

### 1. Functional Options Pattern

```go
type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.Timeout = timeout
    }
}

func NewCircuitBreaker(options ...Option) CircuitBreaker {
    config := &Config{/* defaults */}
    for _, option := range options {
        option(config)
    }
    return &circuitBreakerImpl{config: *config}
}
```

### 2. Error Wrapping

Distinguish between circuit breaker errors and operation errors:

```go
var (
    ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
    ErrTooManyRequests   = errors.New("too many requests in half-open state")
)

func (cb *CircuitBreaker) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
    if err := cb.canExecute(); err != nil {
        return nil, err  // Circuit breaker error
    }
    
    result, err := operation()  // Original operation error
    cb.recordResult(err == nil)
    return result, err
}
```

### 3. State Management

Clean state transitions with proper cleanup:

```go
func (cb *CircuitBreaker) setState(newState State) {
    if cb.state == newState {
        return
    }
    
    oldState := cb.state
    cb.state = newState
    cb.lastStateChange = time.Now()
    
    // Reset state-specific data
    switch newState {
    case StateClosed:
        cb.metrics = Metrics{}  // Reset metrics
    case StateHalfOpen:
        cb.halfOpenRequests = 0  // Reset request counter
    }
    
    // Trigger callback
    if cb.config.OnStateChange != nil {
        cb.config.OnStateChange(cb.name, oldState, newState)
    }
}
```

## Best Practices

### 1. Choose Appropriate Thresholds

- **Failure Threshold**: Too low = unnecessary tripping, too high = delayed protection
- **Timeout Duration**: Balance between service recovery time and user experience
- **Half-Open Requests**: Enough to test recovery but not overwhelm

### 2. Implement Proper Monitoring

```go
func (cb *CircuitBreaker) GetMetrics() Metrics {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    
    // Return copy to prevent data races
    return Metrics{
        Requests:            cb.metrics.Requests,
        Successes:           cb.metrics.Successes,
        Failures:            cb.metrics.Failures,
        ConsecutiveFailures: cb.metrics.ConsecutiveFailures,
        LastFailureTime:     cb.metrics.LastFailureTime,
    }
}
```

### 3. Handle Different Failure Types

Not all errors should trip the circuit:

```go
func (cb *CircuitBreaker) shouldCountAsFailure(err error) bool {
    // Don't count client errors (4xx) as circuit breaker failures
    if httpErr, ok := err.(*HTTPError); ok {
        return httpErr.StatusCode >= 500
    }
    
    // Don't count context cancellation as failure
    if errors.Is(err, context.Canceled) {
        return false
    }
    
    return true
}
```

### 4. Graceful Degradation

Provide fallback mechanisms:

```go
func CallWithFallback(cb CircuitBreaker, primary, fallback func() (interface{}, error)) (interface{}, error) {
    result, err := cb.Call(context.Background(), primary)
    if err != nil && errors.Is(err, ErrCircuitBreakerOpen) {
        return fallback()
    }
    return result, err
}
```

## Testing Strategies

### 1. State Transition Testing

Verify correct state changes under various conditions:

```go
func TestStateTransitions(t *testing.T) {
    cb := NewCircuitBreaker(Config{
        ReadyToTrip: func(m Metrics) bool {
            return m.ConsecutiveFailures >= 3
        },
        Timeout: 100 * time.Millisecond,
    })
    
    // Test Closed -> Open
    for i := 0; i < 3; i++ {
        cb.Call(ctx, failingOperation)
    }
    assert.Equal(t, StateOpen, cb.GetState())
    
    // Test Open -> Half-Open
    time.Sleep(150 * time.Millisecond)
    cb.Call(ctx, successOperation)
    assert.Equal(t, StateClosed, cb.GetState())
}
```

### 2. Concurrency Testing

Ensure thread safety under load:

```go
func TestConcurrentAccess(t *testing.T) {
    cb := NewCircuitBreaker(config)
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cb.Call(ctx, someOperation)
        }()
    }
    
    wg.Wait()
    // Verify metrics consistency
}
```

### 3. Mock Operations

Create controllable operations for testing:

```go
type MockOperation struct {
    shouldFail bool
    delay      time.Duration
    callCount  int32
}

func (m *MockOperation) Execute() (interface{}, error) {
    atomic.AddInt32(&m.callCount, 1)
    
    if m.delay > 0 {
        time.Sleep(m.delay)
    }
    
    if m.shouldFail {
        return nil, errors.New("operation failed")
    }
    return "success", nil
}
```

## Real-World Applications

### 1. HTTP Client Wrapper

```go
type ResilientHTTPClient struct {
    client  *http.Client
    breaker CircuitBreaker
}

func (c *ResilientHTTPClient) Get(url string) (*http.Response, error) {
    result, err := c.breaker.Call(context.Background(), func() (interface{}, error) {
        return c.client.Get(url)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.(*http.Response), nil
}
```

### 2. Database Connection Pool

```go
type ResilientDB struct {
    db      *sql.DB
    breaker CircuitBreaker
}

func (rdb *ResilientDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    result, err := rdb.breaker.Call(context.Background(), func() (interface{}, error) {
        return rdb.db.Query(query, args...)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.(*sql.Rows), nil
}
```

### 3. Microservice Communication

```go
type ServiceClient struct {
    baseURL string
    breaker CircuitBreaker
    client  *http.Client
}

func (sc *ServiceClient) CallService(endpoint string, data interface{}) (interface{}, error) {
    return sc.breaker.Call(context.Background(), func() (interface{}, error) {
        // Implement HTTP call to microservice
        resp, err := sc.client.Post(sc.baseURL+endpoint, "application/json", data)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        
        if resp.StatusCode >= 500 {
            return nil, fmt.Errorf("server error: %d", resp.StatusCode)
        }
        
        // Parse response...
        return response, nil
    })
}
```

## Advanced Features

### 1. Multiple Circuit Breakers

Different services may need different configurations:

```go
type CircuitBreakerRegistry struct {
    breakers map[string]CircuitBreaker
    configs  map[string]Config
}

func (r *CircuitBreakerRegistry) GetBreaker(serviceName string) CircuitBreaker {
    if breaker, exists := r.breakers[serviceName]; exists {
        return breaker
    }
    
    config := r.configs[serviceName]
    breaker := NewCircuitBreaker(config)
    r.breakers[serviceName] = breaker
    return breaker
}
```

### 2. Health Check Integration

```go
func (cb *CircuitBreaker) HealthCheck() error {
    state := cb.GetState()
    metrics := cb.GetMetrics()
    
    if state == StateOpen {
        return fmt.Errorf("circuit breaker is open: %d consecutive failures", 
            metrics.ConsecutiveFailures)
    }
    
    if metrics.Failures > 0 && float64(metrics.Failures)/float64(metrics.Requests) > 0.5 {
        return fmt.Errorf("high failure rate: %.2f%%", 
            float64(metrics.Failures)/float64(metrics.Requests)*100)
    }
    
    return nil
}
```

### 3. Metrics Export

```go
func (cb *CircuitBreaker) ExportMetrics() map[string]interface{} {
    metrics := cb.GetMetrics()
    state := cb.GetState()
    
    return map[string]interface{}{
        "state":                state.String(),
        "total_requests":       metrics.Requests,
        "successful_requests":  metrics.Successes,
        "failed_requests":      metrics.Failures,
        "consecutive_failures": metrics.ConsecutiveFailures,
        "last_failure_time":    metrics.LastFailureTime,
        "failure_rate":         float64(metrics.Failures) / float64(metrics.Requests),
    }
}
```

## Summary

The Circuit Breaker pattern is essential for building resilient distributed systems. It provides:

- **Failure Detection**: Automatically detects when services are failing
- **Fast Failure**: Prevents resource waste by failing quickly
- **Automatic Recovery**: Tests service recovery and automatically resumes normal operation
- **System Protection**: Prevents cascading failures across services

Key implementation considerations:
- Thread safety for concurrent access
- Configurable thresholds and timeouts
- Proper error handling and classification
- Comprehensive testing including race conditions
- Integration with monitoring and alerting systems

This pattern is widely used in production systems at companies like Netflix, Amazon, and Google to ensure system reliability and user experience. 