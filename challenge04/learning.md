# Learning Materials for Concurrent Graph BFS Queries

## Concurrency and Goroutines in Go

Go was designed with concurrency as a core feature, making it easy to write programs that efficiently use multiple CPU cores and handle asynchronous tasks. This challenge focuses on implementing concurrent breadth-first search (BFS) in graph traversal.

### Goroutines

Goroutines are lightweight threads managed by the Go runtime. They allow you to run functions concurrently with minimal overhead:

```go
// Basic goroutine
go functionName()  // Runs the function in a separate goroutine

// Anonymous function as a goroutine
go func() {
    // Do work here
    fmt.Println("Running in a goroutine")
}()
```

Compared to traditional threads, goroutines:
- Are much cheaper (a few KB of memory vs MB for threads)
- Are managed by Go's runtime scheduler instead of the OS
- Can scale to hundreds of thousands or millions on a single machine

### Channels

Channels are the primary mechanism for communication between goroutines. They provide a way to send and receive values with synchronization built in:

```go
// Create a channel
ch := make(chan int)  // Unbuffered channel
bufferedCh := make(chan string, 10)  // Buffered channel with capacity 10

// Send values (blocks if channel is full)
ch <- 42
bufferedCh <- "hello"

// Receive values (blocks if channel is empty)
value := <-ch
message := <-bufferedCh

// Close a channel when done (optional)
close(ch)

// Check if channel is closed
value, ok := <-ch  // ok is false if channel is closed
```

### Channel Patterns

Several common patterns for using channels effectively:

#### Fan-out / Fan-in

Use multiple goroutines to process data in parallel, then combine results:

```go
func fanOut(input []int) <-chan int {
    // Distribute work to multiple goroutines
    out := make(chan int)
    
    go func() {
        defer close(out)
        for _, n := range input {
            out <- process(n)
        }
    }()
    
    return out
}

func fanIn(channels ...<-chan int) <-chan int {
    // Combine results from multiple channels
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

#### Worker Pools

Create a fixed number of workers that process tasks from a queue:

```go
func workerPool(numWorkers int, tasks <-chan Task, results chan<- Result) {
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for task := range tasks {
                results <- processTask(task)
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
}
```

### sync Package

The `sync` package provides synchronization primitives:

```go
// WaitGroup: wait for a group of goroutines to finish
var wg sync.WaitGroup
wg.Add(n)  // Add n goroutines to wait for
wg.Done()  // Mark one goroutine as complete
wg.Wait()  // Block until all goroutines are done

// Mutex: protect access to shared data
var mu sync.Mutex
mu.Lock()
// Critical section (only one goroutine at a time)
mu.Unlock()

// RWMutex: allows multiple readers or one writer
var rwMu sync.RWMutex
rwMu.RLock() // Read lock (multiple allowed)
// Read shared data
rwMu.RUnlock()

rwMu.Lock() // Write lock (exclusive)
// Modify shared data
rwMu.Unlock()
```

### Context Package

The `context` package helps manage cancellation and timeouts in concurrent operations:

```go
// Create a context with cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // Call cancel when done

// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Checking if context is done
select {
case <-ctx.Done():
    // Context was cancelled or timed out
    return ctx.Err()
case result := <-resultChan:
    return result
}
```

### Graph Traversal with BFS

Breadth-First Search visits all vertices at the current depth before moving to the next level. Key concepts:

- **Queue-based approach**: Use a queue data structure to track nodes to visit
- **Visited tracking**: Keep track of visited nodes to avoid cycles
- **Level-by-level processing**: Process all nodes at current distance before moving to next

### Concurrent BFS Considerations

When implementing concurrent BFS, consider:

1. **Goroutine coordination**: Using a goroutine pool for processing nodes at each level
2. **Communication patterns**: Using channels to communicate between workers
3. **Synchronization**: Using sync.WaitGroup to wait for each level to complete
4. **Shared state protection**: Using mutex to protect the visited map if shared across goroutines
5. **Work distribution**: How to divide the graph traversal work among goroutines
6. **Result aggregation**: How to collect results from multiple goroutines safely

### Important Concurrency Patterns for BFS

- **Worker pools**: Fixed number of workers processing BFS tasks
- **Level synchronization**: Ensuring all nodes at one level are processed before moving to next
- **Shared state management**: Protecting visited nodes map across goroutines
- **Channel communication**: Using channels for distributing work and collecting results

### Concurrency Gotchas

Common pitfalls to avoid:

1. **Race Conditions**: Always protect shared data with mutex or channels
2. **Deadlocks**: Avoid situations where goroutines wait for each other indefinitely
3. **Goroutine Leaks**: Ensure goroutines can exit when their work is done
4. **Channel Misuse**: Be careful about channel closing - only the sender should close

## Further Reading

- [Go Concurrency Patterns](https://blog.golang.org/pipelines)
- [Effective Go: Concurrency](https://golang.org/doc/effective_go#concurrency)
- [Visualizing Concurrency in Go](https://divan.dev/posts/go_concurrency_visualize/)
- [The Context Package](https://blog.golang.org/context) 