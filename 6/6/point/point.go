package point

import "fmt"

// Point is a simple struct type
// to represent a single point
type Point struct {
	x int
	y int
}

// Set is a method for
// setting the coordinates
func (p *Point) Set(x, y int) {
	p.x, p.y = x, y
}

// String is a method for user-friandly
// content output
func (p Point) String() string {
	return fmt.Sprintf("x: %d;\ty: %d", p.x, p.y)
}

// Difference is the structs' method
// for yielding a difference
func (p Point) Difference(q Point) Point {
	return Point{p.x - q.x, p.y - q.y}
}

// Difference an exported function
// for yielding a difference
func Difference(p, q Point) Point {
	return Point{p.x - q.x, p.y - q.y}
}
