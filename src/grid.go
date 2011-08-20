// 2D lattice where each site is characterized by a boolean value.
// Sites where the boolean is true are 'active' and can form clusters.
package percolation

import (
	"os"
	"fmt"
)

const ShapeError = "Grid data must be rectangular and have at least one point"
const BoundsError = "Grid point out of bounds"

// Return true if and only if gridData is a M x N rectangle
// (i.e. all sub-arrays are the same length), where M > 0 and N > 0.
func CheckDimensions(gridData [][]bool) bool {
	// Must have at least one point.
	// (short-circuit evaluation protects the second check)
	if len(gridData) == 0 || len(gridData[0]) == 0 {
		return false
	}
	// All sub-arrays must have the same size as the first.
	size := len(gridData[0])
	for _, subArray := range gridData {
		if len(subArray) != size {
			return false
		}
	}
	return true
}

type Grid struct {
	data [][]bool
}

// Construct a Grid object from initData.  
// Ensure that initData has the correct shape to be a Grid.
// Returns nil and an error if the shape is not correct.
func NewGrid(initData [][]bool) (*Grid, os.Error) {
	if !CheckDimensions(initData) {
		return nil, fmt.Errorf(ShapeError)
	}
	grid := new(Grid)
	grid.data = initData
	return grid, nil
}

// Width of the grid.
func (g *Grid) Lx() int {
	return len(g.data)
}

// Height of the grid.
func (g *Grid) Ly() int {
	return len(g.data[0])
}

// Is the point (x, y) within the grid boundaries?
func (g *Grid) InBounds(x, y int) bool {
	if x < 0 || y < 0 || x >= g.Lx() || y >= g.Ly() {
		return false
	}
	return true
}

// Get the grid value at (x, y). Return an error if (x, y) is out of bounds.
func (g *Grid) Get(x, y int) (bool, os.Error) {
	if !g.InBounds(x, y) {
		return false, fmt.Errorf(BoundsError)
	}
	return g.data[x][y], nil
}

// Set the grid value at (x, y). Return an error if (x, y) is out of bounds.
func (g *Grid) Set(x, y int, value bool) os.Error {
	if !g.InBounds(x, y) {
		return fmt.Errorf(BoundsError)
	}
	g.data[x][y] = value
	return nil
}
