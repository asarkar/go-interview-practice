# Learning Materials for Go Generics

## Introduction to Generics in Go

Go 1.18 introduced support for generic programming, allowing developers to write code that can work with multiple types while maintaining type safety. This feature enables more flexible and reusable code without sacrificing compile-time type checking.

### Why Generics?

Before generics, Go developers had several approaches to handle multiple types:

1. **Interface{}**: Using the empty interface allowed functions to accept any type, but required type assertions and lost compile-time type checking.
2. **Code generation**: Tools like `go generate` could create type-specific implementations, but added complexity to the build process.
3. **Copy and paste**: Duplicating code for different types led to maintenance issues.

Generics solve these problems by providing a way to write code that is both type-safe and reusable across multiple types.

### Basic Syntax

The basic syntax for defining a generic function in Go:

```go
func MyGenericFunction[T any](param T) T {
    // Function body
    return param
}
```

And for a generic type:

```go
type MyGenericType[T any] struct {
    Value T
}
```

### Type Parameters and Constraints

Type parameters are specified in square brackets `[T any]` where:
- `T` is the type parameter name
- `any` is the type constraint (in this case, any type is allowed)

Go provides several predefined constraints in the `constraints` package:

```go
import "golang.org/x/exp/constraints"

// Function that works with any ordered type
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### Custom Type Constraints

You can define custom constraints using interface types:

```go
// Define a constraint that requires String() method
type Stringer interface {
    String() string
}

// Function that works with any type implementing String()
func PrintValue[T Stringer](value T) {
    fmt.Println(value.String())
}
```

### Union Types in Constraints

Go generics support union types in constraints, allowing a parameter to accept multiple specific types:

```go
// A constraint that accepts either int or float64
type Number interface {
    int | float64
}

// Function that works with either int or float64
func Add[T Number](a, b T) T {
    return a + b
}
```

### Type Sets

The concept of type sets is central to Go's generics implementation. A constraint defines a set of types that satisfy it:

```go
// Constraint for types that can be compared with == and !=
type Comparable[T any] interface {
    comparable
}

// Function that checks if two values are equal
func AreEqual[T comparable](a, b T) bool {
    return a == b
}
```

### Generic Data Structures

Generics are particularly useful for implementing data structures:

```go
// Generic Stack implementation
type Stack[T any] struct {
    elements []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{elements: make([]T, 0)}
}

func (s *Stack[T]) Push(element T) {
    s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() (T, error) {
    var zero T
    if len(s.elements) == 0 {
        return zero, errors.New("stack is empty")
    }
    
    lastIndex := len(s.elements) - 1
    element := s.elements[lastIndex]
    s.elements = s.elements[:lastIndex]
    return element, nil
}

func (s *Stack[T]) Peek() (T, error) {
    var zero T
    if len(s.elements) == 0 {
        return zero, errors.New("stack is empty")
    }
    
    return s.elements[len(s.elements)-1], nil
}

func (s *Stack[T]) Size() int {
    return len(s.elements)
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.elements) == 0
}
```

### Type Inference

Go can often infer the type parameters from the arguments:

```go
func Identity[T any](value T) T {
    return value
}

// Type inference in action
str := Identity("hello")    // T is inferred as string
num := Identity(42)         // T is inferred as int
```

### Multiple Type Parameters

Functions and types can have multiple type parameters:

```go
// Map function that converts a slice of one type to another
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Usage
numbers := []int{1, 2, 3, 4}
squares := Map(numbers, func(x int) int { return x * x })
// squares: [1, 4, 9, 16]

// Convert numbers to strings
strNumbers := Map(numbers, func(x int) string { return strconv.Itoa(x) })
// strNumbers: ["1", "2", "3", "4"]
```

### Combining Generics with Methods

Methods can be defined on generic types:

```go
type Pair[T, U any] struct {
    First  T
    Second U
}

func (p Pair[T, U]) Swap() Pair[U, T] {
    return Pair[U, T]{First: p.Second, Second: p.First}
}

// Usage
pair := Pair[string, int]{First: "answer", Second: 42}
swapped := pair.Swap() // Pair[int, string]{First: 42, Second: "answer"}
```

### Constraints Package

The `golang.org/x/exp/constraints` package provides useful constraints:

```go
import "golang.org/x/exp/constraints"

// Function that works with any integer type
func Sum[T constraints.Integer](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// Function that works with any floating point type
func Average[T constraints.Float](values []T) T {
    sum := T(0)
    for _, v := range values {
        sum += v
    }
    return sum / T(len(values))
}
```

Key constraints include:
- `Integer`: any integer type
- `Float`: any floating-point type
- `Complex`: any complex number type
- `Ordered`: any type that supports the < operator
- `Signed`: any signed integer type
- `Unsigned`: any unsigned integer type

### Generic Algorithms

Generics are ideal for implementing algorithms that work with multiple types:

```go
// Generic binary search function
func BinarySearch[T constraints.Ordered](slice []T, target T) int {
    left, right := 0, len(slice)-1
    
    for left <= right {
        mid := (left + right) / 2
        
        if slice[mid] == target {
            return mid
        } else if slice[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1 // Not found
}
```

### Type Parameters in Methods

Methods themselves cannot have type parameters separate from the receiver type, but you can work around this with generic functions:

```go
// This won't compile - methods can't have their own type parameters
// func (s *Stack[T]) ConvertTo[U any](converter func(T) U) []U { ... }

// Instead, use a regular function
func ConvertStack[T, U any](stack *Stack[T], converter func(T) U) []U {
    result := make([]U, stack.Size())
    for i, v := range stack.elements {
        result[i] = converter(v)
    }
    return result
}
```

### Zero Values and Generic Types

When working with generics, it's often necessary to produce a "zero value" of the type parameter:

```go
func GetZero[T any]() T {
    var zero T
    return zero
}

// Usage
zeroInt := GetZero[int]()       // 0
zeroString := GetZero[string]() // ""
```

### Performance Considerations

Generics in Go are implemented with careful attention to performance:

1. **Compilation approach**: Go uses a hybrid approach, generating specific code for each type instantiation while sharing as much code as possible.
2. **Runtime efficiency**: Generic code is optimized at compile time, so there's minimal runtime overhead compared to manually written type-specific code.
3. **Code size**: Using many type instantiations can increase binary size, but the compiler works to minimize this impact.

### Best Practices for Using Generics

1. **Don't overuse generics**: Use generics when they provide clear benefits in terms of code reuse and type safety.
2. **Be specific with constraints**: Use the most specific constraint possible for your use case.
3. **Provide clear documentation**: Document the expected behavior of generic functions and types clearly.
4. **Consider performance implications**: Be mindful of how generics affect compilation time and binary size.

### Real-World Examples

#### Generic Result Type

A common pattern is to create a generic result type for handling success and error cases:

```go
type Result[T any] struct {
    Value T
    Error error
}

func NewSuccess[T any](value T) Result[T] {
    return Result[T]{Value: value, Error: nil}
}

func NewError[T any](err error) Result[T] {
    var zero T
    return Result[T]{Value: zero, Error: err}
}

// Usage
func DivideInts(a, b int) Result[int] {
    if b == 0 {
        return NewError[int](errors.New("division by zero"))
    }
    return NewSuccess(a / b)
}
```

#### Generic Set Implementation

```go
type Set[T comparable] struct {
    elements map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
    return Set[T]{elements: make(map[T]struct{})}
}

func (s *Set[T]) Add(element T) {
    s.elements[element] = struct{}{}
}

func (s *Set[T]) Remove(element T) {
    delete(s.elements, element)
}

func (s *Set[T]) Contains(element T) bool {
    _, exists := s.elements[element]
    return exists
}

func (s *Set[T]) Size() int {
    return len(s.elements)
}

func (s *Set[T]) Elements() []T {
    result := make([]T, 0, len(s.elements))
    for element := range s.elements {
        result = append(result, element)
    }
    return result
}

// Set operations
func Union[T comparable](s1, s2 Set[T]) Set[T] {
    result := NewSet[T]()
    
    for element := range s1.elements {
        result.Add(element)
    }
    
    for element := range s2.elements {
        result.Add(element)
    }
    
    return result
}

func Intersection[T comparable](s1, s2 Set[T]) Set[T] {
    result := NewSet[T]()
    
    for element := range s1.elements {
        if s2.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}
```

### Further Reading

1. [Go Generics Design Document](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
2. [Go by Example: Generics](https://gobyexample.com/generics)
3. [Go Generics 101](https://go101.org/generics/101.html)
4. [The Go Programming Language Blog: Using Generics in Go](https://go.dev/blog/intro-generics)
5. [Go Generics in Practice](https://bitfieldconsulting.com/golang/generics) 