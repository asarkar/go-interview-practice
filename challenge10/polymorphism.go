package polymorphism

import (
	"errors"
	"fmt"
	"math"
	"slices"
	// Add any necessary imports here
)

// Shape interface defines methods that all shapes must implement
type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer // Includes String() string method
}

// Rectangle represents a four-sided shape with perpendicular sides
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if err := ensurePositive(width, height); err != nil {
		return nil, err
	}
	return &Rectangle{width, height}, nil
}

func ensurePositive(measurements ...float64) error {
	for _, m := range measurements {
		if m <= 0 {
			return fmt.Errorf("measurement less than or equal to zero %f", m)
		}
	}
	return nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle{Width:%f Height:%f}", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if err := ensurePositive(radius); err != nil {
		return nil, err
	}
	return &Circle{radius}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return math.Pi * 2 * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle{Radius:%f}", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if err := ensurePositive(a, b, c); err != nil {
		return nil, err
	}
	if (a+b <= c) || (b+c <= a) || (a+c <= b) {
		return nil, errors.New("triangle sides do not satisfy the triangle inequality theorem")
	}
	return &Triangle{a, b, c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle{Sides:%f %f %f}", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(_ Shape) {
	panic("not implemented")
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	total := 0.0
	for i := range shapes {
		total += shapes[i].Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	largestArea := shapes[0].Area()
	shapeWithLargestArea := shapes[0]

	for i := 1; i < len(shapes); i++ {
		area := shapes[i].Area()
		if area > largestArea {
			largestArea = area
			shapeWithLargestArea = shapes[i]
		}
	}
	return shapeWithLargestArea
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	return slices.SortedFunc(slices.Values(shapes), func(s1, s2 Shape) int {
		order := 0
		diff := s1.Area() - s2.Area()
		if diff < 0 {
			order = -1
		} else if diff > 0 {
			order = 1
		}
		if !ascending {
			order = -order
		}
		return order
	})
}
