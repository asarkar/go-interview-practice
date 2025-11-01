package bfs

import "sync"

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	type MapEntry struct {
		start int
		bfs   []int
	}

	starts := make(chan int)
	results := make(chan MapEntry)

	var wg sync.WaitGroup

	// Start workers
	for range numWorkers {
		wg.Go(func() {
			for s := range starts {
				results <- MapEntry{start: s, bfs: bfs(graph, s)}
			}
		})
	}

	// Send queries
	go func() {
		for _, q := range queries {
			starts <- q
		}
		close(starts)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	res := make(map[int][]int, len(queries))
	for r := range results {
		res[r.start] = r.bfs
	}

	return res
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]struct{})
	queue := []int{start}
	result := make([]int, 0)

	visited[start] = struct{}{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node)

		for _, neighbor := range graph[node] {
			if _, ok := visited[neighbor]; !ok {
				visited[neighbor] = struct{}{}
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}
