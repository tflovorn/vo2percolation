// Unordered set of 2D points
package percolation

import "math"

type PointSet struct {
	// The dimensions of the grid.
	Lx, Ly int
	// data's keys are 1D coordinates covering the grid.  When the value
	// associated with a key is true, that key is part of the point set.
	data map[int]bool
}

// Create a new point set with grid dimensions (Lx, Ly).
func NewPointSet(Lx, Ly int) *PointSet {
	ps := new(PointSet)
	ps.Lx = Lx
	ps.Ly = Ly
	ps.data = make(map[int]bool)
	return ps
}

// Return a slice of all (x, y) points in the set.
func (ps *PointSet) Elements() [][]int {
	elements := [][]int{}
	for k, v := range ps.data {
		if v {
			x, y := ps.convertFrom1D(k)
			elements = append(elements, []int{x, y})
		}
	}
	return elements
}

// Return true if and only if (x, y) is in the point set.
func (ps *PointSet) Contains(x, y int) bool {
	value, ok := ps.data[ps.convertTo1D(x, y)]
	if ok {
		return value
	}
	return false
}

// Add a point to the set.
func (ps *PointSet) Add(x, y int) {
	ps.data[ps.convertTo1D(x, y)] = true
}

// Remove a point from the set.
func (ps *PointSet) Remove(x, y int) {
	ps.data[ps.convertTo1D(x, y)] = false
}

// Convert the 2D x-y coordinates in the Lx by Ly discrete grid to a single
// integer, useful as a map key.  Panics if (x,y) is not on the grid.
func (ps *PointSet) convertTo1D(x, y int) int {
	if x < 0 || y < 0 || x > ps.Lx || y > ps.Ly {
		panic("point set access out of bounds")
	}
	return ps.Lx*y + x
}

// Convert from 1D map keys to x-y grid coordinates.  Panics if the key is not
// on the grid.
func (ps *PointSet) convertFrom1D(key int) (int, int) {
	if key < 0 || key > ps.Lx*ps.Ly {
		panic("point set access out of bounds")
	}
	y := int(math.Floor(float64(key) / float64(ps.Lx)))
	x := key - ps.Lx*y
	return x, y
}
