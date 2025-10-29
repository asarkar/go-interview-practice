[View the Scoreboard](SCOREBOARD.md)

# Challenge 20: Circuit Breaker Pattern

## Problem Statement

Implement the **Circuit Breaker Pattern** to build resilient systems that can handle failures gracefully. A circuit breaker monitors calls to external services and prevents cascading failures when those services become unavailable.

The circuit breaker has three states:
- **Closed**: Normal operation, requests pass through
- **Open**: Service is failing, requests are blocked and fail fast
- **Half-Open**: Testing if service has recovered

You'll implement a flexible circuit breaker that can wrap any function call and provide automatic failure detection and recovery.

## Function Signatures

```go
type CircuitBreaker interface {
    Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
    GetState() State
    GetMetrics() Metrics
}

type State int
const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type Metrics struct {
    Requests          int64
    Successes         int64
    Failures          int64
    ConsecutiveFailures int64
    LastFailureTime     time.Time
}

func NewCircuitBreaker(config Config) CircuitBreaker
```

## Configuration

```go
type Config struct {
    MaxRequests      uint32        // Max requests allowed in half-open state
    Interval         time.Duration // Statistical window for closed state
    Timeout          time.Duration // Time to wait before half-open
    ReadyToTrip      func(Metrics) bool // Function to determine when to trip
    OnStateChange    func(name string, from State, to State) // State change callback
}
```

## Requirements

### 1. State Management
- **Closed → Open**: When `ReadyToTrip` returns true
- **Open → Half-Open**: After `Timeout` duration
- **Half-Open → Closed**: ~~When operation succeeds~~ When `Config.MaxRequest` number of operations succeed consecutively
- **Half-Open → Open**: When operation fails

### 2. Request Handling
- **Closed**: Allow all requests, track metrics
- **Open**: Reject requests immediately with `ErrCircuitBreakerOpen`
- **Half-Open**: Allow up to `MaxRequests`, then decide state

### 3. Metrics Tracking
- Count total requests, successes, failures
- Track consecutive failures
- Record last failure time
- Reset metrics when transitioning to closed state

## Sample Usage

```go
// Create circuit breaker for external API calls
cb := NewCircuitBreaker(Config{
    MaxRequests: 3,
    Interval:    time.Minute,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(m Metrics) bool {
        return m.ConsecutiveFailures >= 5
    },
})

// Use circuit breaker to wrap API calls
result, err := cb.Call(ctx, func() (interface{}, error) {
    return httpClient.Get("https://api.example.com/data")
})
```

## Test Scenarios

Your implementation will be tested with:

1. **Normal Operation**: Circuit remains closed for successful calls
2. **Failure Detection**: Circuit opens after consecutive failures
3. **Fast Fail**: Requests fail immediately when circuit is open
4. **Recovery Testing**: Circuit transitions to half-open after timeout
5. **Full Recovery**: Circuit closes after successful half-open requests
6. **Concurrent Safety**: Multiple goroutines using the same circuit breaker

## Error Types

```go
var (
    ErrCircuitBreakerOpen    = errors.New("circuit breaker is open")
    ErrTooManyRequests      = errors.New("too many requests in half-open state")
)
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-20/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the Circuit Breaker pattern with all required functionality.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-20/` directory:

```bash
go test -v -race
```

## Difficulty: 🔶 Intermediate

This challenge tests your understanding of:
- Design patterns for resilience
- State management and concurrency
- Error handling strategies
- Metrics collection
- Thread-safe programming 