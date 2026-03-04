# Learning Materials for Concurrent Web Content Aggregator

## Advanced Concurrency Patterns in Go

This challenge focuses on implementing a system that concurrently fetches and processes web content, employing advanced concurrency patterns, context handling, and rate limiting.

### Concurrency vs. Parallelism

- **Concurrency**: Structuring a program as independently executing components
- **Parallelism**: Executing multiple computations simultaneously

Go enables both through goroutines and the runtime scheduler:

```go
// Run multiple tasks concurrently
go task1()
go task2()
go task3()
```

### Web Scraping Basics

Fetching web content in Go:

```go
import (
    "io"
    "net/http"
)

func fetchURL(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    return string(body), nil
}
```

### Advanced Context Usage

The `context` package helps manage cancellation, deadlines, and request values:

```go
import (
    "context"
    "net/http"
    "time"
)

// With timeout
func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    return string(body), err
}

// With cancellation
func fetchMultipleURLs(urls []string) <-chan string {
    ctx, cancel := context.WithCancel(context.Background())
    results := make(chan string)
    
    // Launch goroutine for each URL
    for _, url := range urls {
        go func(u string) {
            resp, err := fetchWithContext(ctx, u)
            if err != nil {
                cancel() // Cancel all other requests if one fails
                return
            }
            results <- resp
        }(url)
    }
    
    return results
}
```

### Context Values

Pass request-scoped values through the call chain:

```go
// Define custom key types to avoid collisions
type contextKey string

const (
    userKey   contextKey = "user"
    requestID contextKey = "request-id"
)

// Store values in context
func enrichContext(ctx context.Context) context.Context {
    ctx = context.WithValue(ctx, userKey, "admin")
    ctx = context.WithValue(ctx, requestID, uuid.New().String())
    return ctx
}

// Retrieve values from context
func processWithContext(ctx context.Context, url string) {
    user, ok := ctx.Value(userKey).(string)
    if !ok {
        user = "anonymous"
    }
    
    id := ctx.Value(requestID)
    
    // Use the values
    log.Printf("User %s (request %v) processing URL: %s", user, id, url)
}
```

### Rate Limiting

Control the rate of requests to avoid overwhelming servers or hitting API limits:

```go
import (
    "context"
    "golang.org/x/time/rate"
    "net/http"
)

// Client with rate limiting
type RateLimitedClient struct {
    client  *http.Client
    limiter *rate.Limiter
}

func NewRateLimitedClient(rps float64, burst int) *RateLimitedClient {
    return &RateLimitedClient{
        client:  &http.Client{},
        limiter: rate.NewLimiter(rate.Limit(rps), burst),
    }
}

func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
    // Wait for rate limiter
    err := c.limiter.Wait(req.Context())
    if err != nil {
        return nil, err
    }
    
    // Perform the request
    return c.client.Do(req)
}

// Usage
client := NewRateLimitedClient(1.0, 5) // 1 request per second, bursts of 5
```

### Worker Pools

Limit the number of concurrent operations:

```go
func WorkerPool(urls []string, numWorkers int) <-chan string {
    var wg sync.WaitGroup
    results := make(chan string)
    
    // Create job channel
    jobs := make(chan string, len(urls))
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range jobs {
                content, err := fetchURL(url)
                if err != nil {
                    log.Printf("Error fetching %s: %v", url, err)
                    continue
                }
                results <- content
            }
        }()
    }
    
    // Send jobs to workers
    for _, url := range urls {
        jobs <- url
    }
    close(jobs)
    
    // Close results channel when all workers are done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    return results
}
```

### Fan-out, Fan-in Pattern

Process data in multiple stages, distributing work and collecting results:

```go
func fetch(urls <-chan string) <-chan *Result {
    results := make(chan *Result)
    
    go func() {
        defer close(results)
        for url := range urls {
            res, err := fetchURL(url)
            results <- &Result{URL: url, Content: res, Error: err}
        }
    }()
    
    return results
}

func process(results <-chan *Result) <-chan *ProcessedResult {
    processed := make(chan *ProcessedResult)
    
    go func() {
        defer close(processed)
        for res := range results {
            if res.Error != nil {
                continue
            }
            
            // Process the content
            data := extractData(res.Content)
            processed <- &ProcessedResult{URL: res.URL, Data: data}
        }
    }()
    
    return processed
}

func merge(channels ...<-chan *ProcessedResult) <-chan *ProcessedResult {
    var wg sync.WaitGroup
    merged := make(chan *ProcessedResult)
    
    // Function to copy from a channel to the merged channel
    output := func(c <-chan *ProcessedResult) {
        defer wg.Done()
        for val := range c {
            merged <- val
        }
    }
    
    // Start an output goroutine for each input channel
    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }
    
    // Close the merged channel when all output goroutines are done
    go func() {
        wg.Wait()
        close(merged)
    }()
    
    return merged
}

// Usage
func main() {
    urls := make(chan string)
    
    // Distribute work to multiple fetchers (fan-out)
    var fetchers []<-chan *Result
    for i := 0; i < 5; i++ {
        fetchers = append(fetchers, fetch(urls))
    }
    
    // Merge results (fan-in)
    results := mergeFetchResults(fetchers...)
    
    // Process results with multiple processors (fan-out)
    var processors []<-chan *ProcessedResult
    for i := 0; i < 3; i++ {
        processors = append(processors, process(results))
    }
    
    // Merge processed results (fan-in)
    processed := merge(processors...)
    
    // Send URLs to process
    go func() {
        for _, url := range targetURLs {
            urls <- url
        }
        close(urls)
    }()
    
    // Collect results
    for p := range processed {
        fmt.Printf("URL: %s, Data: %v\n", p.URL, p.Data)
    }
}
```

### Error Handling in Concurrent Code

Several strategies for handling errors in concurrent operations:

#### 1. Return errors through channels

```go
type Result struct {
    Value string
    Error error
}

func fetchAsync(url string) <-chan Result {
    result := make(chan Result, 1)
    
    go func() {
        resp, err := http.Get(url)
        if err != nil {
            result <- Result{Error: err}
            close(result)
            return
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        result <- Result{Value: string(body), Error: err}
        close(result)
    }()
    
    return result
}
```

#### 2. Use errgroup for coordinated error handling

```go
import "golang.org/x/sync/errgroup"

func fetchAll(urls []string) ([]string, error) {
    var g errgroup.Group
    results := make([]string, len(urls))
    
    for i, url := range urls {
        i, url := i, url // Create local variables for the closure
        
        g.Go(func() error {
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                return err
            }
            
            results[i] = string(body)
            return nil
        })
    }
    
    // Wait for all HTTP fetches to complete
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return results, nil
}
```

### Retry Logic

Implement retries with backoff to handle transient failures:

```go
func fetchWithRetry(url string, maxRetries int) (string, error) {
    var (
        body  string
        err   error
        sleep time.Duration = 100 * time.Millisecond
    )
    
    for i := 0; i <= maxRetries; i++ {
        if i > 0 {
            log.Printf("Retry #%d for %s after %v", i, url, sleep)
            time.Sleep(sleep)
            sleep *= 2 // Exponential backoff
        }
        
        body, err = fetchURL(url)
        if err == nil {
            return body, nil
        }
        
        // Check if we should retry
        if !isRetryable(err) {
            return "", err
        }
    }
    
    return "", fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

func isRetryable(err error) bool {
    // Check for network errors, 429 Too Many Requests, 5xx Server Errors
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Temporary() {
        return true
    }
    
    var httpErr *url.Error
    if errors.As(err, &httpErr) {
        return httpErr.Timeout() || isRetryableStatusCode(httpErr)
    }
    
    return false
}
```

### Circuit Breaker Pattern

Prevent cascading failures by "breaking the circuit" after too many errors:

```go
type CircuitBreaker struct {
    maxFailures     int
    failureCount    int
    resetTimeout    time.Duration
    lastFailureTime time.Time
    mu              sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
    }
}

func (cb *CircuitBreaker) Execute(op func() error) error {
    cb.mu.Lock()
    
    // Check if circuit is open (too many failures recently)
    if cb.failureCount >= cb.maxFailures {
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            // Reset after timeout
            cb.failureCount = 0
        } else {
            cb.mu.Unlock()
            return fmt.Errorf("circuit open: too many failures")
        }
    }
    
    cb.mu.Unlock()
    
    // Execute the operation
    err := op()
    
    if err != nil {
        cb.mu.Lock()
        cb.failureCount++
        cb.lastFailureTime = time.Now()
        cb.mu.Unlock()
    }
    
    return err
}
```

### Handling External API Dependencies

When aggregating content from external APIs, consider these best practices:

1. **Timeouts**: Set appropriate timeouts for all requests
2. **Caching**: Cache responses to reduce load and improve performance
3. **Fallbacks**: Provide fallback content when services are unavailable
4. **Retries**: Implement retries with backoff for transient failures
5. **Circuit Breakers**: Prevent cascading failures

```go
// Simplified content fetcher with these patterns
type ContentFetcher struct {
    client          *http.Client
    cache           map[string]CachedResponse
    cacheMu         sync.RWMutex
    circuitBreakers map[string]*CircuitBreaker
    rateLimiters    map[string]*rate.Limiter
}

func (f *ContentFetcher) FetchContent(ctx context.Context, url string) (string, error) {
    // Check cache first
    f.cacheMu.RLock()
    if cached, ok := f.cache[url]; ok && !cached.Expired() {
        f.cacheMu.RUnlock()
        return cached.Content, nil
    }
    f.cacheMu.RUnlock()
    
    // Check if domain is rate limited
    domain := extractDomain(url)
    if limiter, ok := f.rateLimiters[domain]; ok {
        if err := limiter.Wait(ctx); err != nil {
            return "", err
        }
    }
    
    // Check circuit breaker
    if breaker, ok := f.circuitBreakers[domain]; ok {
        var content string
        err := breaker.Execute(func() error {
            var fetchErr error
            content, fetchErr = f.fetchWithRetry(ctx, url, 3)
            return fetchErr
        })
        
        if err != nil {
            return f.getFallbackContent(url), nil
        }
        
        // Cache successful response
        f.cacheResponse(url, content)
        return content, nil
    }
    
    // Regular fetch with retry
    content, err := f.fetchWithRetry(ctx, url, 3)
    if err != nil {
        return f.getFallbackContent(url), nil
    }
    
    // Cache successful response
    f.cacheResponse(url, content)
    return content, nil
}
```

## Further Reading

- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)
- [Advanced Go Concurrency Patterns](https://talks.golang.org/2013/advconc.slide)
- [Context Package Documentation](https://pkg.go.dev/context)
- [Rate Limiting in Go](https://pkg.go.dev/golang.org/x/time/rate)
- [Errgroup Package](https://pkg.go.dev/golang.org/x/sync/errgroup)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html) 