# Learning Materials for Slice Operations

## Working with Slices in Go

Slices are one of Go's most versatile and frequently used data structures. They provide a convenient view into an underlying array and are used extensively throughout Go programs.

### Slice Basics

A slice is a reference to a contiguous segment of an array. Unlike arrays, slices are dynamic in size.

#### Creating Slices

```go
// Create a slice with a literal
numbers := []int{1, 2, 3, 4, 5}

// Create an empty slice with make
slice := make([]int, 5)      // slice of length 5, capacity 5
slice := make([]int, 0, 10)  // slice of length 0, capacity 10

// Create a slice from an existing array or slice
array := [5]int{1, 2, 3, 4, 5}
slice := array[1:4]  // [2, 3, 4]
```

#### Slice Length and Capacity

Slices have both a length and a capacity:
- Length: The number of elements the slice contains (`len(slice)`)
- Capacity: The number of elements in the underlying array (`cap(slice)`)

```go
slice := make([]int, 3, 5)
fmt.Println(len(slice))  // 3
fmt.Println(cap(slice))  // 5
```

### Common Slice Operations

#### Appending to a Slice

The `append` function adds elements to the end of a slice and returns a new slice:

```go
slice := []int{1, 2, 3}
slice = append(slice, 4)        // [1, 2, 3, 4]
slice = append(slice, 5, 6, 7)  // [1, 2, 3, 4, 5, 6, 7]

// Append one slice to another
other := []int{8, 9, 10}
slice = append(slice, other...)  // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

#### Slicing a Slice

You can create a slice from a slice using the slicing operator:

```go
slice := []int{1, 2, 3, 4, 5}
sub := slice[1:3]  // [2, 3]
```

#### Copying Slices

The `copy` function copies elements from one slice to another:

```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
n := copy(dst, src)  // n is number of elements copied (3)
```

### Common Slice Algorithms

#### Finding Maximum Value

To find the maximum value in a slice of integers:

```go
func findMax(numbers []int) int {
    if len(numbers) == 0 {
        return 0  // or another default value
    }
    
    max := numbers[0]
    for _, n := range numbers[1:] {
        if n > max {
            max = n
        }
    }
    return max
}
```

#### Removing Duplicates

To remove duplicates while preserving order:

```go
func removeDuplicates(numbers []int) []int {
    if len(numbers) == 0 {
        return []int{}
    }
    
    // Use a map to track seen values
    seen := make(map[int]bool)
    result := make([]int, 0, len(numbers))
    
    for _, n := range numbers {
        if !seen[n] {
            seen[n] = true
            result = append(result, n)
        }
    }
    
    return result
}
```

#### Reversing a Slice

To reverse the order of elements in a slice:

```go
func reverseSlice(slice []int) []int {
    result := make([]int, len(slice))
    for i, v := range slice {
        result[len(slice)-1-i] = v
    }
    return result
}

// Alternative approach that modifies the original slice
func reverseInPlace(slice []int) {
    for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
        slice[i], slice[j] = slice[j], slice[i]
    }
}
```

#### Filtering Elements

To filter elements based on a condition (e.g., keeping only even numbers):

```go
func filterEven(numbers []int) []int {
    result := make([]int, 0)
    for _, n := range numbers {
        if n%2 == 0 {
            result = append(result, n)
        }
    }
    return result
}
```

### Slice Gotchas and Tips

#### Slice Mutations

Slices are references to arrays, so modifying a slice modifies the underlying array:

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3]
sub[0] = 42  // Modifies original too

fmt.Println(original)  // [1, 42, 3, 4, 5]
```

#### Creating Independent Copies

To create an independent copy of a slice:

```go
original := []int{1, 2, 3, 4, 5}
copy := make([]int, len(original))
copy = append(copy, original...)
```

#### Pre-allocating Slices

For efficiency when building slices incrementally, pre-allocate with a capacity:

```go
// Inefficient
var result []int
for i := 0; i < 10000; i++ {
    result = append(result, i)  // Many allocations and copies
}

// Efficient
result := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    result = append(result, i)  // No reallocation needed
}
```

#### Empty vs. Nil Slices

An empty slice has length 0 but is not nil:

```go
var nilSlice []int         // nil, len 0, cap 0
emptySlice := []int{}      // not nil, len 0, cap 0
emptyMake := make([]int, 0) // not nil, len 0, cap 0

fmt.Println(nilSlice == nil)   // true
fmt.Println(emptySlice == nil) // false
```

## Further Reading

- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)
- [Go by Example: Slices](https://gobyexample.com/slices)
- [The Go Programming Language Specification: Slice types](https://go.dev/ref/spec#Slice_types) 