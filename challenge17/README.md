[View the Scoreboard](SCOREBOARD.md)

# Challenge 17: Palindrome Checker

## Problem Statement

Write a function `IsPalindrome` that checks if a given string is a palindrome. A palindrome is a word, phrase, number, or other sequence of characters that reads the same forward and backward (ignoring spaces, punctuation, and capitalization).

## Function Signature

```go
func IsPalindrome(s string) bool
```

## Input Format

- A string `s` containing letters, numbers, spaces and punctuation.

## Output Format

- A boolean value: `true` if the string is a palindrome, `false` otherwise.

## Requirements

1. The function should be case-insensitive ("A" is the same as "a").
2. The function should ignore spaces and punctuation marks.
3. The function should handle alphanumeric strings.

## Sample Input and Output

### Sample Input 1

```
"racecar"
```

### Sample Output 1

```
true
```

### Sample Input 2

```
"A man, a plan, a canal: Panama"
```

### Sample Output 2

```
true
```

### Sample Input 3

```
"hello"
```

### Sample Output 3

```
false
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-17/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the `IsPalindrome` function.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-17/` directory:

```bash
go test -v
``` 