# Learning Materials for Performance Optimization with Benchmarking

## Benchmarking in Go

Go provides excellent built-in support for benchmarking through the `testing` package. Benchmarks are functions that start with the word `Benchmark` followed by a name, take a parameter of type `*testing.B`, and are run by the `go test` command with the `-bench` flag.

```go
func BenchmarkMyFunction(b *testing.B) {
    // Run the target function b.N times
    for i := 0; i < b.N; i++ {
        MyFunction()
    }
}
```

### Running Benchmarks

To run benchmarks, use the `go test` command with the `-bench` flag:

```bash
go test -bench=.                 # Run all benchmarks
go test -bench=BenchmarkMyFunction  # Run a specific benchmark
```

Add the `-benchmem` flag to also measure memory allocations:

```bash
go test -bench=. -benchmem
```

### Understanding Benchmark Results

The output of a benchmark looks like this:

```
BenchmarkMyFunction-8   	10000000	       118 ns/op	      16 B/op	       1 allocs/op
```

This means:
- The benchmark ran on 8 CPU cores (`-8`)
- It ran 10,000,000 iterations
- Each operation took about 118 nanoseconds
- Each operation allocated 16 bytes
- Each operation performed 1 allocation

### Comparing Benchmarks

To compare benchmark results, you can use the `benchstat` tool:

```bash
go test -bench=. -count=5 > old.txt
# Make changes to the code
go test -bench=. -count=5 > new.txt
benchstat old.txt new.txt
```

## Common Performance Issues and Solutions

### 1. Inefficient String Concatenation

**Problem**: Using the `+` operator for string concatenation creates a new string each time, leading to quadratic complexity.

**Inefficient:**
```go
// O(n²) complexity
func ConcatenateStrings(strings []string) string {
    result := ""
    for _, s := range strings {
        result += s // Creates a new string each time
    }
    return result
}
```

**Efficient:**
```go
// O(n) complexity
func ConcatenateStrings(strings []string) string {
    var builder strings.Builder
    for _, s := range strings {
        builder.WriteString(s)
    }
    return builder.String()
}
```

### 2. Unnecessary Memory Allocations

**Problem**: Creating new objects or slices inside loops can lead to excessive GC pressure.

**Inefficient:**
```go
func ProcessItems(items []Item) []Result {
    var results []Result
    for _, item := range items {
        // Allocates a new slice for each item
        data := make([]byte, len(item.Data))
        copy(data, item.Data)
        
        // Process the data
        result := ProcessData(data)
        results = append(results, result)
    }
    return results
}
```

**Efficient:**
```go
func ProcessItems(items []Item) []Result {
    // Pre-allocate the slice with the expected capacity
    results := make([]Result, 0, len(items))
    
    // Reuse a buffer across iterations
    buffer := make([]byte, 0, 1024) // Reasonable starting size
    
    for _, item := range items {
        // Reuse the buffer
        buffer = buffer[:0]
        buffer = append(buffer, item.Data...)
        
        // Process the data
        result := ProcessData(buffer)
        results = append(results, result)
    }
    return results
}
```

### 3. Inefficient Algorithms

**Problem**: Using algorithms with suboptimal complexity for the problem at hand.

**Inefficient (Bubble Sort):**
```go
// O(n²) complexity
func BubbleSort(items []int) {
    for i := 0; i < len(items); i++ {
        for j := 0; j < len(items)-1; j++ {
            if items[j] > items[j+1] {
                items[j], items[j+1] = items[j+1], items[j]
            }
        }
    }
}
```

**Efficient (QuickSort):**
```go
// O(n log n) average case complexity
func QuickSort(items []int) {
    sort.Ints(items) // Uses an efficient sorting algorithm
}
```

### 4. Redundant Calculations

**Problem**: Recalculating values that could be cached or computed once.

**Inefficient:**
```go
// Recursive Fibonacci with exponential complexity
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}
```

**Efficient (Memoization):**
```go
// Linear complexity with memoization
func Fibonacci(n int) int {
    memo := make([]int, n+1)
    return fibMemo(n, memo)
}

func fibMemo(n int, memo []int) int {
    if n <= 1 {
        return n
    }
    
    if memo[n] != 0 {
        return memo[n]
    }
    
    memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
    return memo[n]
}
```

**Even More Efficient (Iterative):**
```go
// Linear complexity with iteration
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

## Profiling in Go

For more detailed performance analysis, Go provides profiling tools in the `runtime/pprof` and `net/http/pprof` packages.

### CPU Profiling

```go
import "runtime/pprof"

func main() {
    // Create a CPU profile file
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    
    // Start CPU profiling
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Run your code
    ExpensiveOperation()
}
```

### Memory Profiling

```go
import "runtime/pprof"

func main() {
    // Run your code first to generate allocations
    ExpensiveOperation()
    
    // Create a memory profile file
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    // Write memory profile
    pprof.WriteHeapProfile(f)
}
```

### Analyzing Profiles

Use the `go tool pprof` command to analyze profiles:

```bash
go tool pprof cpu.prof      # Interactive mode
go tool pprof -http=:8080 cpu.prof  # Web UI
```

## Memory Management Optimization

### Slice Capacity

Pre-allocate slices when you know the approximate size:

```go
// Inefficient - may cause multiple allocations and copies
data := []int{}
for i := 0; i < 10000; i++ {
    data = append(data, i)
}

// Efficient - allocates once with the right capacity
data := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    data = append(data, i)
}
```

### Reducing Pointer Indirection

Prefer value types over pointer types when dealing with small objects:

```go
// More allocations, more GC pressure
type Point struct {
    X, Y float64
}

points := make([]*Point, 1000)
for i := 0; i < 1000; i++ {
    points[i] = &Point{X: float64(i), Y: float64(i)}
}

// Fewer allocations, better cache locality
type Point struct {
    X, Y float64
}

points := make([]Point, 1000)
for i := 0; i < 1000; i++ {
    points[i] = Point{X: float64(i), Y: float64(i)}
}
```

### Sync.Pool for Frequent Allocations

Use `sync.Pool` to reuse temporary objects:

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func ProcessData(data []byte) []byte {
    // Get a buffer from the pool
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer)
    
    // Use the buffer for processing
    buffer = buffer[:0]
    // ... processing logic ...
    
    return result
}
```

## Concurrency Optimizations

### Utilizing Multiple Cores

Use goroutines and channels to parallelize independent work:

```go
func ProcessItems(items []Item) []Result {
    numCPU := runtime.NumCPU()
    numWorkers := numCPU
    
    // Create channels
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                result := ProcessItem(item)
                results <- result
            }
        }()
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Wait for workers and close results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var finalResults []Result
    for result := range results {
        finalResults = append(finalResults, result)
    }
    
    return finalResults
}
```

### Avoiding Goroutine Overhead

Beware of creating too many goroutines for small tasks:

```go
// Inefficient for small items
for _, item := range items {
    go ProcessItem(item) // Goroutine overhead may exceed benefits
}

// More efficient for small items - process in batches
batchSize := 1000
for i := 0; i < len(items); i += batchSize {
    end := i + batchSize
    if end > len(items) {
        end = len(items)
    }
    
    batch := items[i:end]
    go func(batch []Item) {
        for _, item := range batch {
            ProcessItem(item)
        }
    }(batch)
}
```

## Additional Optimization Techniques

### Loop Unrolling

For tight loops with simple operations, unrolling can improve performance:

```go
// Before unrolling
sum := 0
for i := 0; i < len(data); i++ {
    sum += data[i]
}

// After unrolling
sum := 0
remainder := len(data) % 4
for i := 0; i < remainder; i++ {
    sum += data[i]
}
for i := remainder; i < len(data); i += 4 {
    sum += data[i] + data[i+1] + data[i+2] + data[i+3]
}
```

### Reducing Interface Conversions

Avoid frequent type assertions or interface conversions in hot paths:

```go
// Inefficient - type assertion in loop
func ProcessItems(items []interface{}) int {
    sum := 0
    for _, item := range items {
        if val, ok := item.(int); ok {
            sum += val
        }
    }
    return sum
}

// More efficient - use concrete types when possible
func ProcessItems(items []int) int {
    sum := 0
    for _, item := range items {
        sum += item
    }
    return sum
}
```

### Function Inlining

Small, frequently called functions may be inlined by the compiler, but you can help:

```go
// May be inlined automatically if small enough
func add(a, b int) int {
    return a + b
}

// Suggest inlining with the "go:inline" directive
//go:inline
func add(a, b int) int {
    return a + b
}

// Prevent inlining of large functions with "go:noinline"
//go:noinline
func complexFunction(data []int) int {
    // Complex logic...
}
```

## Best Practices for Performance Optimization

1. **Measure First**: Always benchmark before and after optimization to confirm improvements
2. **80/20 Rule**: Focus on the 20% of the code that causes 80% of the performance issues
3. **Start Simple**: Use efficient algorithms and data structures before optimizing at a lower level
4. **Readability vs. Performance**: Balance performance gains against code maintainability
5. **Profile in Production-like Environments**: Performance characteristics can vary between environments
6. **Test Case Sizes**: Test with different input sizes to understand algorithmic complexity 