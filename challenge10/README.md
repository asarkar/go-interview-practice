[View the Scoreboard](SCOREBOARD.md)

# Challenge 10: Polymorphic Shape Calculator

## Problem Statement

Implement a system to calculate properties of various geometric shapes using Go interfaces. This challenge focuses on understanding and correctly implementing Go's interface system to enable polymorphism.

## Requirements

1. Implement a `Shape` interface with the following methods:
   - `Area() float64`: Calculates the area of the shape
   - `Perimeter() float64`: Calculates the perimeter (or circumference) of the shape
   - `String() string`: Returns a string representation of the shape (implementing fmt.Stringer)

2. Implement the following concrete shapes:
   - `Rectangle`: Defined by width and height
   - `Circle`: Defined by radius
   - `Triangle`: Defined by three sides (use Heron's formula for area)

3. Implement a `ShapeCalculator` that can:
   - Take any shape and return its properties
   - Calculate the total area of multiple shapes
   - Find the shape with the largest area from a collection
   - Sort shapes by area in ascending or descending order

## Function Signatures

```go
// Shape interface
type Shape interface {
    Area() float64
    Perimeter() float64
    fmt.Stringer // Includes String() string method
}

// Concrete types
type Rectangle struct {
    Width, Height float64
}

type Circle struct {
    Radius float64
}

type Triangle struct {
    SideA, SideB, SideC float64
}

// Constructor functions
func NewRectangle(width, height float64) (*Rectangle, error)
func NewCircle(radius float64) (*Circle, error)
func NewTriangle(a, b, c float64) (*Triangle, error)

// ShapeCalculator
type ShapeCalculator struct{}

func NewShapeCalculator() *ShapeCalculator
func (sc *ShapeCalculator) PrintProperties(s Shape)
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape
```

## Constraints

- All measurements must be positive values
- Triangle sides must satisfy the triangle inequality theorem (sum of lengths of any two sides must exceed the length of the remaining side)
- Implement proper validation in constructors and return appropriate errors
- Use constants for π (pi) when calculating circle properties
- The `String()` method should return a formatted string with shape type and dimensions

## Sample Usage

```go
// Create shapes
rect, _ := NewRectangle(5, 3)
circle, _ := NewCircle(4)
triangle, _ := NewTriangle(3, 4, 5)

// Use shapes polymorphically
calculator := NewShapeCalculator()
shapes := []Shape{rect, circle, triangle}

// Calculate total area
totalArea := calculator.TotalArea(shapes)
fmt.Printf("Total area: %.2f\n", totalArea)

// Sort shapes by area
sortedShapes := calculator.SortByArea(shapes, true)
for _, s := range sortedShapes {
    calculator.PrintProperties(s)
}

// Find largest shape
largest := calculator.LargestShape(shapes)
fmt.Printf("Largest shape: %s with area %.2f\n", largest, largest.Area())
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-10/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required interfaces, types, and methods.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-10/` directory:

```bash
go test -v
``` 