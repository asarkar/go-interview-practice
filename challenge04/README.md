[View the Scoreboard](SCOREBOARD.md)

# Challenge 4: Concurrent Graph BFS Queries

You are required to concurrently process multiple breadth-first search (BFS) queries on a single graph. Each query specifies a starting node, and you must compute the BFS order from that node. Unlike a simple single-threaded BFS, your solution should utilize goroutines and channels (or concurrency-safe data structures) to handle multiple queries efficiently and in parallel.

## Function Signature

```go
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int
```

Parameters:
- graph: A representation of the graph as an adjacency list. 
  - The key is a node (an integer).
  - The value is a slice of adjacent nodes.
- queries: A slice of starting nodes for which BFS must be performed.
- numWorkers: The number of goroutines (workers) that concurrently handle these BFS queries.

Returns:
- A map from the query node to the BFS order starting from that node.

## Requirements

1. You must use concurrency (goroutines + channels, or concurrency-safe data structures) to process BFS queries in parallel.
2. A naive or purely sequential approach may be too slow, especially for large graphs and many queries.
3. The BFS algorithm itself can be standard (using a queue), but each BFS query should run concurrently if workers are available.

## Example Usage (Not Tested by the Official Tests)

```go
func main() {
    graph := map[int][]int{
        0: {1, 2},
        1: {2, 3},
        2: {3},
        3: {4},
        4: {},
    }
    queries := []int{0, 1, 2}
    numWorkers := 2

    results := ConcurrentBFSQueries(graph, queries, numWorkers)
    /*
       Possible output:
       results[0] = [0 1 2 3 4]
       results[1] = [1 2 3 4]
       results[2] = [2 3 4]
    */
}
```

## Instructions

1. Fork this repository and clone your fork.  
2. Create a directory for your submission: `challenge-4/submissions/<yourgithubusername>/`.  
3. Copy `solution-template.go` into your submission directory.  
4. Implement the function `ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int`.  
5. Use goroutines to handle BFS queries in parallel, respecting the number of workers.  
6. Open a pull request with your solution.  

## Testing Locally

From inside `challenge-4/`, run:

```bash
go test -v
```