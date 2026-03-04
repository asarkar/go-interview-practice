# Challenge 11: Concurrent Web Content Aggregator

## Problem Statement

Implement a concurrent web content aggregator that fetches, processes, and aggregates data from multiple sources with proper concurrency control and context handling.

## Requirements

1. Implement a `ContentAggregator` that:
   - Concurrently fetches content from multiple URLs
   - Processes the content (extract specific information)
   - Aggregates results with proper error handling
   - Uses proper context management for cancellation and timeouts
   - Implements rate limiting to avoid overwhelming sources

2. You must implement the following concurrency patterns:
   - **Worker Pool**: Process fetched content using a fixed number of worker goroutines
   - **Fan-Out, Fan-In**: Distribute processing tasks and collect results
   - **Context handling**: Proper propagation of cancellation and timeout signals
   - **Rate Limiting**: Limit the rate of requests using a token bucket or similar approach
   - **Concurrent data structures**: Safe access to shared data

3. The solution should demonstrate understanding of:
   - Goroutines and channel management
   - Proper error handling in concurrent code
   - Synchronization primitives (Mutex, RWMutex, WaitGroup)
   - Context package for managing request lifecycles
   - Graceful shutdown

## Function Signatures

```go
// Core types
type ContentFetcher interface {
    Fetch(ctx context.Context, url string) ([]byte, error)
}

type ContentProcessor interface {
    Process(ctx context.Context, content []byte) (ProcessedData, error)
}

type ProcessedData struct {
    Title       string
    Description string
    Keywords    []string
    Timestamp   time.Time
    Source      string
}

type ContentAggregator struct {
    // Add fields as needed
}

// Constructor function
func NewContentAggregator(
    fetcher ContentFetcher, 
    processor ContentProcessor, 
    workerCount int, 
    requestsPerSecond int,
) *ContentAggregator

// Methods
func (ca *ContentAggregator) FetchAndProcess(
    ctx context.Context, 
    urls []string,
) ([]ProcessedData, error)

func (ca *ContentAggregator) Shutdown() error

// Helper functions for different concurrency patterns
func (ca *ContentAggregator) workerPool(
    ctx context.Context, 
    jobs <-chan string, 
    results chan<- ProcessedData,
    errors chan<- error,
)

func (ca *ContentAggregator) fanOut(
    ctx context.Context, 
    urls []string,
) ([]ProcessedData, []error)
```

## Constraints

- The solution must handle errors gracefully and never lose error information
- Implement proper resource cleanup (close channels, release locks, etc.)
- The number of concurrent requests should be configurable
- The request rate limiting must be implemented
- Timeout and cancellation must be properly handled
- The code should guard against goroutine leaks

## Sample Usage

```go
// Create content fetcher and processor
fetcher := &HTTPFetcher{
    Client: &http.Client{Timeout: 5 * time.Second},
}
processor := &HTMLProcessor{}

// Create aggregator with 5 workers and 10 requests per second limit
aggregator := NewContentAggregator(fetcher, processor, 5, 10)

// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// URLs to fetch and process
urls := []string{
    "https://example.com",
    "https://example.org",
    "https://example.net",
    // Add more URLs as needed
}

// Fetch and process in parallel with rate limiting
results, err := aggregator.FetchAndProcess(ctx, urls)
if err != nil {
    log.Fatalf("Error in aggregate operation: %v", err)
}

// Process results
for _, data := range results {
    fmt.Printf("Title: %s\nSource: %s\nKeywords: %v\n\n", 
        data.Title, data.Source, data.Keywords)
}

// Clean up
aggregator.Shutdown()
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-11/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required interfaces and types.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-11/` directory:

```bash
go test -v
``` 