package perf

import (
	"fmt"
	"strings"
	"time"
)

// SlowSort sorts a slice of integers using a very inefficient algorithm (bubble sort)
// TODO: Optimize this function to be more efficient
func SlowSort(data []int) []int {
	// Make a copy to avoid modifying the original
	result := make([]int, len(data))
	copy(result, data)

	// Bubble sort implementation
	for range result {
		for j := range len(result) - 1 {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	return result
}

// OptimizedSort is your optimized version of SlowSort.
// It should produce identical results but perform better.
//
// Benchmark:
// +--------------------------------+------------+---------------+
// | Benchmark                      | Iterations | Time per Op   |
// +--------------------------------+------------+---------------+
// | BenchmarkSlowSort/10-10        | 14,492,658 |   82.83 ns/op |
// | BenchmarkSlowSort/100-10       |    176,055 |   6,785 ns/op |
// | BenchmarkSlowSort/1000-10      |      1,651 | 721,573 ns/op |
// | BenchmarkOptimizedSort/10-10   | 25,138,983 |   48.82 ns/op |
// | BenchmarkOptimizedSort/100-10  |  1,551,513 |   704.4 ns/op |
// | BenchmarkOptimizedSort/1000-10 |    130,231 |   9,158 ns/op |
// +--------------------------------+------------+---------------+

func OptimizedSort(data []int) []int {
	// Hint: Consider using sort package or a more efficient algorithm
	res := make([]int, len(data))
	copy(res, data)
	SortDualPivot(res, 0, len(res)-1)
	return res
}

// https://www.youtube.com/watch?v=XYVbjQXkmiI
// https://www.youtube.com/watch?v=r3a25XPf2A8
func SortDualPivot(nums []int, lo, hi int) {
	if lo >= hi {
		return
	}

	leftPivotIndex, rightPivotIndex := partition(nums, lo, hi)

	// Recurse on the three partitions:
	// 1) elements < left pivot
	// 2) elements between pivots
	// 3) elements > right pivot
	SortDualPivot(nums, lo, leftPivotIndex-1)
	SortDualPivot(nums, leftPivotIndex+1, rightPivotIndex-1)
	SortDualPivot(nums, rightPivotIndex+1, hi)
}

// partition rearranges nums[lo:hi+1] around two pivots (initially at lo and hi).
// After partitioning:
//   - all elements < leftPivotValue are to the left of leftPivotIndex,
//   - all elements between the pivots are between leftPivotIndex and rightPivotIndex,
//   - all elements > rightPivotValue are to the right of rightPivotIndex.
//
// It returns the final indices of the left and right pivots.
func partition(nums []int, lo, hi int) (int, int) {
	// Ensure left pivot <= right pivot, swap if necessary.
	if nums[lo] > nums[hi] {
		nums[lo], nums[hi] = nums[hi], nums[lo]
	}

	leftPivotValue := nums[lo]
	rightPivotValue := nums[hi]

	// leftBoundary: the right boundary (inclusive) for the region where all
	// elements < leftPivotValue.
	// rightBoundary: the left boundary (inclusive) for the region where all
	// elements > rightPivotValue.
	leftBoundary := lo + 1
	rightBoundary := hi - 1

	// cursor scans the slice between leftBoundary and rightBoundary (inclusive).
	cursor := leftBoundary

	for cursor <= rightBoundary {
		switch {
		case nums[cursor] < leftPivotValue:
			// Place nums[cursor] into the left region.
			nums[cursor], nums[leftBoundary] = nums[leftBoundary], nums[cursor]
			leftBoundary++
			cursor++
		case nums[cursor] > rightPivotValue:
			// Place nums[cursor] into the right region.
			nums[cursor], nums[rightBoundary] = nums[rightBoundary], nums[cursor]
			rightBoundary--
			// Note: Do not advance cursor here because the swapped-in element needs evaluation.
		default:
			// Between the two pivots.
			cursor++
		}
	}

	// Put pivots into their final places.
	leftBoundary--
	nums[lo], nums[leftBoundary] = nums[leftBoundary], nums[lo]

	rightBoundary++
	nums[hi], nums[rightBoundary] = nums[rightBoundary], nums[hi]

	return leftBoundary, rightBoundary
}

// InefficientStringBuilder builds a string by repeatedly concatenating
// TODO: Optimize this function to be more efficient
func InefficientStringBuilder(parts []string, repeatCount int) string {
	result := ""

	for range repeatCount {
		for _, part := range parts {
			result += part
		}
	}

	return result
}

// OptimizedStringBuilder is your optimized version of InefficientStringBuilder.
// It should produce identical results but perform better.
//
// Benchmark 1, with preallocated size.
// +---------------------------------------------+------------+--------------+
// | Benchmark                                   | Iterations | Time per Op  |
// +---------------------------------------------+------------+--------------+
// | BenchmarkInefficientStringBuilder/Small-10  | 1,584,055  |     759.4 ns |
// | BenchmarkInefficientStringBuilder/Medium-10 |    16,587  |    71,380 ns |
// | BenchmarkInefficientStringBuilder/Large-10  |       188  | 6,354,634 ns |
// | BenchmarkOptimizedStringBuilder/Small-10    | 6,453,026  |     185.5 ns |
// | BenchmarkOptimizedStringBuilder/Medium-10   |   419,534  |     2,826 ns |
// | BenchmarkOptimizedStringBuilder/Large-10    |    37,254  |    32,119 ns |
// +---------------------------------------------+------------+--------------+
//
// Benchmark 2, default size.
// +---------------------------------------------+------------+--------------+
// | Benchmark                                   | Iterations | Time per Op  |
// +---------------------------------------------+------------+--------------+
// | BenchmarkInefficientStringBuilder/Small-10  | 1,606,806  |     756.1 ns |
// | BenchmarkInefficientStringBuilder/Medium-10 |    16,682  |    71,648 ns |
// | BenchmarkInefficientStringBuilder/Large-10  |       189  | 6,327,952 ns |
// | BenchmarkOptimizedStringBuilder/Small-10    | 5,830,648  |     205.6 ns |
// | BenchmarkOptimizedStringBuilder/Medium-10   |   390,987  |     3,065 ns |
// | BenchmarkOptimizedStringBuilder/Large-10    |    36,820  |    32,509 ns |
// +---------------------------------------------+------------+--------------+
//
// Clearly, preallocating is faster, even if it takes an additional iteration over the `parts`.

func OptimizedStringBuilder(parts []string, repeatCount int) string {
	size := 0
	for i := range parts {
		size += len(parts[i])
	}
	// Hint: Consider using strings.Builder or bytes.Buffer
	var sb strings.Builder
	sb.Grow(size)

	for range repeatCount {
		for i := range parts {
			sb.WriteString(parts[i])
		}
	}

	return sb.String()
}

// ExpensiveCalculation performs a computation with redundant work
// It computes the sum of all fibonacci numbers up to n
// TODO: Optimize this function to be more efficient
func ExpensiveCalculation(n int) int {
	if n <= 0 {
		return 0
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += fibonacci(i)
	}

	return sum
}

// Helper function that computes the fibonacci number at position n
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
//
// Beanchmark:
// +-----------------------------------------+------------+--------------+
// | Benchmark                               | Iterations | Time per Op  |
// +-----------------------------------------+------------+--------------+
// | BenchmarkExpensiveCalculation/Small-10  |  2,496,655 |     478.6 ns |
// | BenchmarkExpensiveCalculation/Medium-10 |     20,018 |    59,857 ns |
// | BenchmarkExpensiveCalculation/Large-10  |        162 | 7,545,767 ns |
// | BenchmarkOptimizedCalculation/Small-10  | 63,297,398 |     18.51 ns |
// | BenchmarkOptimizedCalculation/Medium-10 | 29,563,537 |     40.66 ns |
// | BenchmarkOptimizedCalculation/Large-10  | 19,654,690 |     60.79 ns |
// +-----------------------------------------+------------+--------------+
func OptimizedCalculation(n int) int {
	if n < 0 {
		// Negative Fibonacci sequence formula is F(n-2) = F(n) - F(n-1).
		// Starting with F(0)=0 and F(1)=1, we can calculate:
		// F(-1) = F(1) - F(0) 	 = 1 - 0    =  1
		// F(-2) = F(0) - F(-1)  = 0 - 1    = -1
		// F(-3) = F(-1) - F(-2) = 1 - (-1) =  2
		// The members of the sequence have alternating signs.
		//
		// We don't want to deal with negative `n` for this exercise.
		panic(fmt.Sprintf("expected non-negative `n`, got %d", n))
	}
	if n < 2 {
		return n
	}
	res := 1
	fibs := [2]int{0, 1}
	for i := 2; i <= n; i++ {
		fibs[i%2] = fibs[0] + fibs[1]
		res += fibs[i%2]
	}
	return res
}

// HighAllocationSearch searches for all occurrences of a substring and creates a map with their
// positions.
// TODO: Optimize this function to reduce allocations
//
// Beanchmark:
// +----------------------------------------------+------------+-------------+
// | Benchmark                                    | Iterations | Time per Op |
// +----------------------------------------------+------------+-------------+
// | BenchmarkHighAllocationSearch/Short_Text-10  | 4,070,913  |    292.3 ns |
// | BenchmarkHighAllocationSearch/Medium_Text-10 |   502,993  |    2,314 ns |
// | BenchmarkHighAllocationSearch/Long_Text-10   |    51,723  |   22,605 ns |
// | BenchmarkOptimizedSearch/Short_Text-10       | 3,863,580  |    309.3 ns |
// | BenchmarkOptimizedSearch/Medium_Text-10      |   529,297  |    2,258 ns |
// | BenchmarkOptimizedSearch/Long_Text-10        |    62,539  |   19,194 ns |
// +----------------------------------------------+------------+-------------+
//
// +----------------------------------------+------------+-------------+----------+-----------+
// | Benchmark                              | Iterations | Time per Op | B/op     | Allocs/op |
// +----------------------------------------+------------+-------------+----------+-----------+
// | BenchmarkMemoryHighAllocationSearch-10 |    53,097  | 22,455 ns   | 11,816 B | 12        |
// | BenchmarkMemoryOptimizedSearch-10.     |    62,920  | 20,411 ns   | 28,880 B | 14        |
// +--------------------------------------------+------------+---------+----------+-----------+
//
// B/op — Bytes per Operation: The average number of bytes allocated on the heap for a single
// iteration of the benchmarked function.
// Allocs/op — Allocations per Operation: The average number of heap allocations that occurred
// during a single iteration.
//
// Even though KMP is asymptotically faster for large inputs, for short to medium texts with small
// patterns and case-insensitive comparison, the simpler linear scan with pre-lowercased strings
// (HighAllocationSearch) can outperform KMP in real-world benchmarks.
func HighAllocationSearch(text, substr string) map[int]string {
	result := make(map[int]string)

	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)

	for i := range len(lowerText) {
		// Check if we can fit the substring starting at position i
		if i+len(lowerSubstr) <= len(lowerText) {
			// Extract the potential match
			potentialMatch := lowerText[i : i+len(lowerSubstr)]

			// Check if it matches
			if potentialMatch == lowerSubstr {
				// Store the original case version
				result[i] = text[i : i+len(substr)]
			}
		}
	}

	return result
}

// OptimizedSearch is your optimized version of HighAllocationSearch
// It should produce identical results but perform better with fewer allocations
func OptimizedSearch(text, substr string) map[int]string {
	// Hint: Consider avoiding temporary string allocations and reusing memory
	haystack := []rune(strings.ToLower(text))
	needle := []rune(strings.ToLower(substr))
	n := len(needle)
	lps := make([]int, n)

	// Build LPS.
	kmpLoop(needle, needle, 1, lps, true)

	// Search.
	starts := kmpLoop(haystack, needle, 0, lps, false)

	result := make(map[int]string, len(starts))
	for i := range starts {
		result[starts[i]] = text[starts[i] : starts[i]+n]
	}

	return result
}

// kmpLoop either builds the LPS array or finds matches, depending on buildLPS flag.
func kmpLoop(haystack, needle []rune, start int, lps []int, buildLPS bool) []int {
	i, j := 0, start
	var matches []int

	for j < len(haystack) {
		switch {
		case needle[i] == haystack[j]:
			i++
			j++
			if buildLPS && i < len(lps) {
				lps[i] = i
			} else if i == len(needle) {
				matches = append(matches, j-i)
				i = lps[i-1]
			}
		case i > 0:
			i = lps[i-1]
		default:
			if buildLPS {
				lps[i] = 0
			}
			j++
		}
	}

	return matches
}

// A function to simulate CPU-intensive work for benchmarking
// You don't need to optimize this; it's just used for testing
func SimulateCPUWork(duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		// Just waste CPU cycles
		for i := range 1_000_000 {
			_ = i
		}
	}
}
