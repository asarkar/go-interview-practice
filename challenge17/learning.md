# Learning Materials for Palindrome Checker

## String Manipulation in Go

In this challenge, you'll work with string manipulation in Go. Go provides excellent support for working with strings through the `strings` package and other built-in functions.

### Basic String Operations

#### String Length

You can get the length of a string using the built-in `len()` function:

```go
str := "Hello, World!"
length := len(str)  // length is 13
```

#### String Comparison

Strings can be compared using standard comparison operators:

```go
str1 := "apple"
str2 := "banana"
if str1 == str2 {
    // Strings are equal
}
```

### String Manipulation

#### Converting Case

The `strings` package provides functions to convert between uppercase and lowercase:

```go
import "strings"

str := "Hello"
lower := strings.ToLower(str)  // "hello"
upper := strings.ToUpper(str)  // "HELLO"
```

#### Removing Characters

To remove characters from a string, you can:

1. Use `strings.ReplaceAll()` to replace characters with an empty string:

```go
import "strings"

str := "Hello, World!"
withoutCommas := strings.ReplaceAll(str, ",", "")  // "Hello World!"
```

2. Build a new string, keeping only the characters you want:

```go
func removeNonAlphanumeric(s string) string {
    var result strings.Builder
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
            result.WriteRune(r)
        }
    }
    return result.String()
}
```

### Working with Runes

Go represents Unicode characters as runes (alias for int32). When iterating over a string with a range loop, you get the runes, not bytes:

```go
str := "Hello, 世界"
for i, r := range str {
    fmt.Printf("%d: %c (%d)\n", i, r, r)
}
```

### Reversing a String

Go doesn't have a built-in function to reverse a string, but you can write one:

```go
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

### Checking Palindromes

To check if a string is a palindrome:

1. Clean the string (remove spaces, punctuation, and convert to lowercase)
2. Check if the cleaned string is the same forwards and backwards

```go
func isPalindrome(s string) bool {
    // Convert to lowercase
    s = strings.ToLower(s)
    
    // Remove non-alphanumeric characters
    var cleaned strings.Builder
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
            cleaned.WriteRune(r)
        }
    }
    
    // Get the cleaned string
    cleanedStr := cleaned.String()
    
    // Check if it's a palindrome
    for i := 0; i < len(cleanedStr)/2; i++ {
        if cleanedStr[i] != cleanedStr[len(cleanedStr)-1-i] {
            return false
        }
    }
    
    return true
}
```

## Performance Considerations

### strings.Builder vs String Concatenation

When building strings, use `strings.Builder` instead of concatenating with `+` to avoid creating multiple temporary strings:

```go
// Inefficient
result := ""
for i := 0; i < 1000; i++ {
    result += "a"  // Creates a new string each time
}

// Efficient
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.WriteString("a")
}
result := builder.String()
```

### Unicode Considerations

Remember that strings in Go are byte sequences, not character sequences. When dealing with Unicode characters, convert to runes:

```go
str := "Hello, 世界"
runes := []rune(str)
```

## Further Reading

- [Go by Example: Strings](https://gobyexample.com/strings)
- [strings package documentation](https://pkg.go.dev/strings)
- [unicode package documentation](https://pkg.go.dev/unicode)
- [Go Blog: Strings, bytes, runes and characters in Go](https://blog.golang.org/strings) 