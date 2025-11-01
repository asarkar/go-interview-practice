# Learning Materials for String Pattern Matching

## Introduction to String Pattern Matching

String pattern matching is the process of finding occurrences of a pattern string within a larger text string. This fundamental algorithm is used in many applications such as:

- Text editors (find and replace)
- Bioinformatics (DNA sequence matching)
- Internet search engines
- Data mining and analysis
- Spam filters and security systems

In this learning material, we'll explore three different algorithms for string pattern matching:

1. Naive (Brute Force) Algorithm
2. Knuth-Morris-Pratt (KMP) Algorithm
3. Rabin-Karp Algorithm

## 1. Naive Pattern Matching Algorithm

The naive approach is the most straightforward method for string matching. It checks for a match starting at each possible position in the text.

### How the Naive Algorithm Works

1. Align the pattern at the beginning of the text
2. Compare each character of the pattern with the corresponding character in the text
3. If all characters match, record the starting position
4. Shift the pattern by one position to the right
5. Repeat steps 2-4 until the end of the text is reached

### Implementing the Naive Algorithm in Go

```go
func NaivePatternMatch(text, pattern string) []int {
    matches := []int{}
    
    // Handle edge cases
    if len(pattern) == 0 || len(text) < len(pattern) {
        return matches
    }
    
    // Check each possible position in the text
    for i := 0; i <= len(text)-len(pattern); i++ {
        j := 0
        
        // Check if the pattern matches at this position
        for j < len(pattern) && text[i+j] == pattern[j] {
            j++
        }
        
        // If j reached the end of the pattern, we found a match
        if j == len(pattern) {
            matches = append(matches, i)
        }
    }
    
    return matches
}
```

### Complexity of the Naive Algorithm

- **Time Complexity**: O(n*m) where n is the length of the text and m is the length of the pattern
- **Space Complexity**: O(k) where k is the number of matches found

The naive algorithm is simple to implement but can be inefficient for large texts or patterns with many partial matches.

## 2. Knuth-Morris-Pratt (KMP) Algorithm

The KMP algorithm improves upon the naive approach by avoiding redundant comparisons when a mismatch occurs. It uses a preprocessed table (called the "failure function" or "longest proper prefix which is also suffix" array) to skip characters that we know will match based on previously matched characters.

### How the KMP Algorithm Works

1. Preprocess the pattern to build a partial match table (also called the "LPS" or "π" table)
2. Use this table to determine how far to shift the pattern when a mismatch occurs
3. Never backtrack in the text - each character in the text is examined exactly once

### Creating the LPS (Longest Prefix Suffix) Array

The LPS array helps determine the longest proper prefix of the pattern that is also a suffix of the pattern up to each position. This information is used to avoid redundant comparisons.

```go
func computeLPSArray(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    
    // Length of the previous longest prefix suffix
    length := 0
    i := 1
    
    // The loop calculates lps[i] for i = 1 to m-1
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            // This is the tricky part
            if length != 0 {
                length = lps[length-1]
                // Note: We do not increment i here
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    
    return lps
}
```

### Implementing the KMP Algorithm in Go

```go
func KMPSearch(text, pattern string) []int {
    matches := []int{}
    
    // Handle edge cases
    if len(pattern) == 0 || len(text) < len(pattern) {
        return matches
    }
    
    n := len(text)
    m := len(pattern)
    
    // Preprocess the pattern
    lps := computeLPSArray(pattern)
    
    i := 0 // Index for text
    j := 0 // Index for pattern
    
    for i < n {
        // Current characters match, move both pointers forward
        if pattern[j] == text[i] {
            i++
            j++
        }
        
        // Found a complete match
        if j == m {
            matches = append(matches, i-j)
            // Use lps to shift pattern for next match
            j = lps[j-1]
        } else if i < n && pattern[j] != text[i] {
            // Mismatch after j matches
            if j != 0 {
                // Use lps to shift pattern
                j = lps[j-1]
            } else {
                // No match found, move to next character in text
                i++
            }
        }
    }
    
    return matches
}
```

### Complexity of the KMP Algorithm

- **Time Complexity**: O(n+m) where n is the length of the text and m is the length of the pattern
- **Space Complexity**: O(m) for the LPS array plus O(k) for storing the matches

The KMP algorithm is much more efficient than the naive approach for texts with many potential matches, especially for longer patterns.

## 3. Rabin-Karp Algorithm

The Rabin-Karp algorithm uses hashing to find pattern matches more efficiently. Instead of comparing each character, it compares hash values of the pattern and substrings of the text.

### How the Rabin-Karp Algorithm Works

1. Compute the hash value of the pattern
2. Compute hash values for all possible m-length substrings of the text using a rolling hash function
3. Compare the hash value of the pattern with the hash value of each substring
4. If the hash values match, verify the actual strings character by character

### Implementing a Rolling Hash Function

A rolling hash function allows us to compute the hash value of the next substring in constant time by using the hash value of the current substring:

```go
// Remove the leftmost character and add the rightmost character
newHash = (oldHash - oldChar * pow) * base + newChar
```

### Implementing the Rabin-Karp Algorithm in Go

```go
func RabinKarpSearch(text, pattern string) []int {
    matches := []int{}
    
    // Handle edge cases
    if len(pattern) == 0 || len(text) < len(pattern) {
        return matches
    }
    
    n := len(text)
    m := len(pattern)
    
    // Large prime number to avoid hash collisions
    prime := 101
    
    // Base value for the hash function
    base := 256
    
    // Hash value for pattern and initial window
    patternHash := 0
    windowHash := 0
    
    // Highest power of base that we need
    h := 1
    for i := 0; i < m-1; i++ {
        h = (h * base) % prime
    }
    
    // Calculate initial hash values
    for i := 0; i < m; i++ {
        patternHash = (base*patternHash + int(pattern[i])) % prime
        windowHash = (base*windowHash + int(text[i])) % prime
    }
    
    // Slide the pattern over text one by one
    for i := 0; i <= n-m; i++ {
        // Check if hash values match
        if patternHash == windowHash {
            // Verify the match character by character
            match := true
            for j := 0; j < m; j++ {
                if text[i+j] != pattern[j] {
                    match = false
                    break
                }
            }
            if match {
                matches = append(matches, i)
            }
        }
        
        // Calculate hash value for next window
        if i < n-m {
            windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime
            
            // Ensure we only have positive hash values
            if windowHash < 0 {
                windowHash += prime
            }
        }
    }
    
    return matches
}
```

### Complexity of the Rabin-Karp Algorithm

- **Average Case Time Complexity**: O(n+m) where n is the length of the text and m is the length of the pattern
- **Worst Case Time Complexity**: O(n*m) if there are many hash collisions
- **Space Complexity**: O(k) where k is the number of matches found

The Rabin-Karp algorithm is particularly efficient for multiple pattern searching and plagiarism detection.

## Comparing the Algorithms

| Algorithm | Average Time Complexity | Worst Case Time Complexity | Space Complexity | Strengths | Weaknesses |
|-----------|-------------------------|----------------------------|------------------|-----------|------------|
| Naive     | O(n*m)                 | O(n*m)                    | O(1)             | Simple, low overhead | Inefficient for large strings |
| KMP       | O(n+m)                 | O(n+m)                    | O(m)             | Very efficient, no backtracking | More complex, requires preprocessing |
| Rabin-Karp | O(n+m)                | O(n*m)                    | O(1)             | Good for multiple patterns | Hash collisions can occur |

## Practical Applications

### 1. Text Editors

Pattern matching is essential for "find" and "replace" operations in text editors.

### 2. Bioinformatics

DNA sequence matching uses pattern matching to find gene sequences within the genome.

### 3. Intrusion Detection

Network security systems use pattern matching to identify suspicious patterns in network traffic.

### 4. Plagiarism Detection

Document similarity checks use pattern matching to find copied content.

## Further Reading

1. [Knuth-Morris-Pratt Algorithm (GeeksforGeeks)](https://www.geeksforgeeks.org/kmp-algorithm-for-pattern-searching/)
2. [Rabin-Karp Algorithm (GeeksforGeeks)](https://www.geeksforgeeks.org/rabin-karp-algorithm-for-pattern-searching/)
3. [String Matching Algorithms (Wikipedia)](https://en.wikipedia.org/wiki/String-searching_algorithm)
4. [Advanced String Searching (Stanford CS166)](https://web.stanford.edu/class/cs166/) 