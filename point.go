// A point on the 2D grid.
package vo2percolation

type Point struct {
	// Coordinates of the point.
	x, y int
}

func (p Point) X() int {
	return p.x
}

func (p Point) Y() int {
	return p.y
}
