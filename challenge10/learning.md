# Learning Materials for Polymorphic Shape Calculator

## Interfaces and Polymorphism in Go

This challenge focuses on using Go's interfaces to implement polymorphism for geometric shape calculations.

### Understanding Interfaces in Go

In Go, interfaces define behavior without specifying implementation. An interface is a collection of method signatures:

```go
// Define an interface
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

A type implements an interface implicitly by implementing its methods:

```go
// Rectangle implements the Shape interface
type Rectangle struct {
    Width  float64
    Height float64
}

// Implement the Area method
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Implement the Perimeter method
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle also implements the Shape interface
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}
```

### Using Interfaces for Polymorphism

Interfaces allow for polymorphic behavior—different types can be treated uniformly based on their behavior:

```go
// Function that works with any Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

// Usage
rect := Rectangle{Width: 5, Height: 3}
circ := Circle{Radius: 2}

PrintShapeInfo(rect)  // Works with Rectangle
PrintShapeInfo(circ)  // Works with Circle
```

### Interface Values

An interface value consists of two components:
1. The dynamic type: The concrete type stored in the interface
2. The dynamic value: The actual value of that type

```go
var s Shape                // nil interface value (nil type, nil value)
s = Rectangle{5, 3}        // s has type Rectangle, value Rectangle{5, 3}
s = Circle{2.5}            // s now has type Circle, value Circle{2.5}
```

### Empty Interface

The empty interface `interface{}` or `any` (Go 1.18+) has no methods and can hold any value:

```go
func PrintAny(a interface{}) {
    fmt.Println(a)
}

PrintAny(42)              // Works with int
PrintAny("Hello")         // Works with string
PrintAny(Rectangle{5, 3}) // Works with Rectangle
```

### Type Assertions

Type assertions extract the underlying value from an interface:

```go
// Type assertion with single return value
rect := s.(Rectangle) // Panics if s is not a Rectangle

// Type assertion with check
rect, ok := s.(Rectangle)
if ok {
    fmt.Println("It's a rectangle with width:", rect.Width)
} else {
    fmt.Println("It's not a rectangle")
}
```

### Type Switches

Type switches handle multiple types:

```go
func Describe(s Shape) string {
    switch v := s.(type) {
    case Rectangle:
        return fmt.Sprintf("Rectangle with width %.2f and height %.2f", v.Width, v.Height)
    case Circle:
        return fmt.Sprintf("Circle with radius %.2f", v.Radius)
    case nil:
        return "nil shape"
    default:
        return fmt.Sprintf("Unknown shape of type %T", v)
    }
}
```

### Interface Composition

Interfaces can be composed of other interfaces:

```go
type Sizer interface {
    Area() float64
}

type Perimeterer interface {
    Perimeter() float64
}

// Composed interface
type Shape interface {
    Sizer
    Perimeterer
    String() string  // Additional method
}
```

### Embedding Interfaces

Go allows embedding one interface into another:

```go
type Stringer interface {
    String() string
}

type Shape interface {
    Area() float64
    Perimeter() float64
}

// CompleteShape embeds Shape and Stringer
type CompleteShape interface {
    Shape
    Stringer
}
```

### Interface Implementation with Pointer Receivers

Method receiver types matter for interface implementation:

```go
type Modifier interface {
    Scale(factor float64)
}

// Value receiver - doesn't modify original
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver - modifies original
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

var m Modifier
r := Rectangle{5, 3}

// This works - r is addressable
m = &r
m.Scale(2)

// This doesn't work - interface expects pointer receiver
// m = r // Compile error
```

### Interface Best Practices

1. **Keep interfaces small**: Prefer interfaces with few methods (often just one)
2. **Define interfaces at the point of use**: Define them in the package that uses them, not where they're implemented
3. **Interfaces as behavior, not types**: Focus on what something does, not what it is

```go
// Good - defines behavior
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Less good - defines a type
type Car interface {
    Drive()
    Stop()
    Refuel()
}
```

### The Liskov Substitution Principle

The Liskov Substitution Principle states that objects of a superclass should be replaceable with objects of a subclass without affecting program correctness:

```go
// A common violation is adding requirements in subtypes
type Parallelogram interface {
    SetWidth(w float64)
    SetHeight(h float64)
    Area() float64
}

type Rectangle struct {
    width, height float64
}

func (r *Rectangle) SetWidth(w float64) { r.width = w }
func (r *Rectangle) SetHeight(h float64) { r.height = h }
func (r Rectangle) Area() float64 { return r.width * r.height }

type Square struct {
    side float64
}

// This implementation breaks expectations!
func (s *Square) SetWidth(w float64) {
    s.side = w
    // Square changes both dimensions when one is set
}

func (s *Square) SetHeight(h float64) {
    s.side = h
}

func (s Square) Area() float64 { return s.side * s.side }
```

### Practical Example: Shape Calculator

Let's implement a complete shape calculator:

```go
package shape

import (
    "fmt"
    "math"
)

// Shape is the basic interface
type Shape interface {
    Area() float64
    Perimeter() float64
    String() string
}

// Circle implementation
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
    return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Rectangle implementation
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
    return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Triangle implementation
type Triangle struct {
    SideA float64
    SideB float64
    SideC float64
}

func (t Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

func (t Triangle) Area() float64 {
    // Heron's formula
    s := t.Perimeter() / 2
    return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t Triangle) String() string {
    return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator handles multiple shapes
type ShapeCalculator struct {
    shapes []Shape
}

func NewCalculator() *ShapeCalculator {
    return &ShapeCalculator{shapes: make([]Shape, 0)}
}

func (c *ShapeCalculator) AddShape(s Shape) {
    c.shapes = append(c.shapes, s)
}

func (c *ShapeCalculator) TotalArea() float64 {
    total := 0.0
    for _, s := range c.shapes {
        total += s.Area()
    }
    return total
}

func (c *ShapeCalculator) TotalPerimeter() float64 {
    total := 0.0
    for _, s := range c.shapes {
        total += s.Perimeter()
    }
    return total
}

func (c *ShapeCalculator) ListShapes() []string {
    result := make([]string, len(c.shapes))
    for i, s := range c.shapes {
        result[i] = s.String()
    }
    return result
}
```

### Extending with New Shapes

One advantage of interfaces is the ability to add new types without changing existing code:

```go
// Add a new shape: Regular Polygon
type RegularPolygon struct {
    Sides     int
    SideLength float64
}

func (p RegularPolygon) Perimeter() float64 {
    return float64(p.Sides) * p.SideLength
}

func (p RegularPolygon) Area() float64 {
    return (float64(p.Sides) * p.SideLength * p.SideLength) / (4 * math.Tan(math.Pi/float64(p.Sides)))
}

func (p RegularPolygon) String() string {
    return fmt.Sprintf("RegularPolygon(sides=%d, length=%.2f)", p.Sides, p.SideLength)
}

// Works with the existing calculator without changes
calculator.AddShape(RegularPolygon{Sides: 6, SideLength: 5})
```

### Testing with Interfaces

Interfaces facilitate testing by allowing mock implementations:

```go
// Interface definition
type AreaCalculator interface {
    Area() float64
}

// Function that uses the interface
func IsLargeShape(s AreaCalculator) bool {
    return s.Area() > 100
}

// Test with a mock
type MockShape struct{
    MockArea float64
}

func (m MockShape) Area() float64 {
    return m.MockArea
}

func TestIsLargeShape(t *testing.T) {
    small := MockShape{50}
    large := MockShape{150}
    
    if IsLargeShape(small) {
        t.Error("Expected small shape to not be large")
    }
    
    if !IsLargeShape(large) {
        t.Error("Expected large shape to be large")
    }
}
```

## Further Reading

- [Go Interfaces Tutorial](https://tour.golang.org/methods/9)
- [Effective Go: Interfaces](https://golang.org/doc/effective_go#interfaces)
- [SOLID Design in Go](https://dave.cheney.net/2016/08/20/solid-go-design)
- [The Laws of Reflection](https://blog.golang.org/laws-of-reflection) (for understanding interfaces at a deeper level) 