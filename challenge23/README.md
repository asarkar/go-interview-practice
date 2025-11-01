[View the Scoreboard](SCOREBOARD.md)

# Challenge 23: String Pattern Matching

## Problem Statement

Implement efficient string pattern matching algorithms to find all occurrences of a pattern in a text. In this challenge, you'll implement three different pattern matching algorithms:

1. `NaivePatternMatch` - A simple brute force approach that checks every possible position.
2. `KMPSearch` - The Knuth-Morris-Pratt algorithm that avoids unnecessary comparisons by using a preprocessed prefix table.
3. `RabinKarpSearch` - The Rabin-Karp algorithm that uses hashing to find patterns efficiently.

## Function Signatures

```go
func NaivePatternMatch(text, pattern string) []int
func KMPSearch(text, pattern string) []int
func RabinKarpSearch(text, pattern string) []int
```

## Input Format

- `text` - The main text string in which to search for the pattern.
- `pattern` - The pattern string to search for.

## Output Format

- All functions should return a slice of integers containing the starting indices of all occurrences of the pattern in the text.
- If no matches are found, return an empty slice.
- Indices should be 0-based (the first character is at position 0).

## Requirements

1. `NaivePatternMatch` should implement a straightforward brute force algorithm.
2. `KMPSearch` should implement the Knuth-Morris-Pratt algorithm.
3. `RabinKarpSearch` should implement the Rabin-Karp algorithm.
4. All three functions should return the same correct results.
5. Pay attention to edge cases like empty strings, patterns longer than the text, etc.

## Sample Input and Output

### Sample Input 1

```
NaivePatternMatch("ABABDABACDABABCABAB", "ABABCABAB")
```

### Sample Output 1

```
[10]
```

### Sample Input 2

```
KMPSearch("AABAACAADAABAABA", "AABA")
```

### Sample Output 2

```
[0, 9, 12]
```

### Sample Input 3

```
RabinKarpSearch("GEEKSFORGEEKS", "GEEK")
```

### Sample Output 3

```
[0, 8]
```

### Sample Input 4

```
NaivePatternMatch("AAAAAA", "AA")
```

### Sample Output 4

```
[0, 1, 2, 3, 4]
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-23/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required functions.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-23/` directory:

```bash
go test -v
```

## Performance Expectations

- **Naive Algorithm**: O(n*m) time complexity where n is the length of the text and m is the length of the pattern.
- **KMP Algorithm**: O(n+m) time complexity.
- **Rabin-Karp Algorithm**: Average case O(n+m) time complexity, worst case O(n*m). 