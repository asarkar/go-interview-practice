[View the Scoreboard](SCOREBOARD.md)

# Challenge 19: Slice Operations

## Problem Statement

Write functions to perform common operations on slices (Go's dynamic arrays). You'll implement the following functions:

1. `FindMax` - Find the maximum value in a slice of integers.
2. `RemoveDuplicates` - Remove duplicate values from a slice while preserving order.
3. `ReverseSlice` - Reverse the order of elements in a slice.
4. `FilterEven` - Create a new slice containing only even numbers from the original slice.

## Function Signatures

```go
func FindMax(numbers []int) int
func RemoveDuplicates(numbers []int) []int
func ReverseSlice(slice []int) []int
func FilterEven(numbers []int) []int
```

## Input Format

- For all functions, a slice of integers.

## Output Format

- `FindMax` - A single integer representing the maximum value.
- `RemoveDuplicates` - A slice of integers with duplicates removed.
- `ReverseSlice` - A slice of integers in reverse order.
- `FilterEven` - A slice containing only even integers.

## Requirements

1. `FindMax` should return the maximum value from the slice. If the slice is empty, return 0.
2. `RemoveDuplicates` should preserve the original order of elements while removing duplicates.
3. `ReverseSlice` should create a new slice with elements in reverse order.
4. `FilterEven` should return a new slice containing only even numbers.

## Sample Input and Output

### Sample Input 1

```
FindMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
```

### Sample Output 1

```
9
```

### Sample Input 2

```
RemoveDuplicates([]int{3, 1, 4, 1, 5, 9, 2, 6})
```

### Sample Output 2

```
[3 1 4 5 9 2 6]
```

### Sample Input 3

```
ReverseSlice([]int{1, 2, 3, 4, 5})
```

### Sample Output 3

```
[5 4 3 2 1]
```

### Sample Input 4

```
FilterEven([]int{1, 2, 3, 4, 5, 6})
```

### Sample Output 4

```
[2 4 6]
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-19/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required functions.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-19/` directory:

```bash
go test -v
``` 