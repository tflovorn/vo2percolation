// 2D lattice where each site is characterized by a boolean value.
// Sites where the boolean is true are 'active' and can form clusters.
package percolation

import (
	"os"
	"fmt"
	"rand"
	"time"
)

const GridShapeError = "Grid data must be rectangular and contain at least one point"
const GridBoundsError = "Grid point out of bounds"

// 2D lattice with bounds-checked access functions and cluster statistics
type Grid struct {
	data [][]bool // CheckDimensions(data) must be true
}

// Function type for iterating over a Grid
type GridCallback func(x, y int, value bool)

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

// Construct a Grid object from initData.  
// Ensure that initData has the correct shape to be a Grid.
// Returns nil and an error if the shape is not correct.
func NewGrid(initData [][]bool) (*Grid, os.Error) {
	if !CheckDimensions(initData) {
		return nil, fmt.Errorf(GridShapeError)
	}
	grid := new(Grid)
	grid.data = initData
	return grid, nil
}

// Generate a random grid of dimensions Lx and Ly.
func RandomGrid(Lx, Ly int) (*Grid, os.Error) {
	// must have at least one site
	if Lx <= 0 || Ly <= 0 {
		return nil, fmt.Errorf("invalid grid dimensions")
	}
	// build the random grid
	data := [][]bool{}
	for x := 0; x < Lx; x++ {
		data := append(data, []bool{})
		for y := 0; y < Ly; y++ {
			data[x] = append(data[x], RandomBool())
		}
	}
	return NewGrid(data)
}

// Generate a random grid of dimensions Lx and Ly with N active sites.
func RandomConstrainedGrid(Lx, Ly, N int) (*Grid, os.Error) {
	// must have at least one site
	if Lx <= 0 || Ly <= 0 {
		return nil, fmt.Errorf("invalid grid dimensions")
	}
	// silently deal with N < 0 and N > Lx * Ly
	if N < 0 {
		N = 0
	} else if N > Lx*Ly {
		N = Lx * Ly
	}
	// build the empty grid
	data := [][]bool{}
	for x := 0; x < Lx; x++ {
		data := append(data, []bool{})
		for y := 0; y < Ly; y++ {
			data[x] = append(data[x], true)
		}
	}
	// activate random sites
	random := rand.New(rand.NewSource(time.Nanoseconds()))
	for activeCount := 0; activeCount < N; {
		x, y := random.Intn(Lx), random.Intn(Ly)
		if !data[x][y] {
			data[x][y] = true
			activeCount += 1
		}
	}
	return NewGrid(data)
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

// Get the grid value at (x, y). Panic if (x, y) is out of bounds.
func (g *Grid) Get(x, y int) bool {
	if !g.InBounds(x, y) {
		panic(GridBoundsError)
	}
	return g.data[x][y]
}

// Set the grid value at (x, y). Panic if (x, y) is out of bounds.
func (g *Grid) Set(x, y int, value bool) {
	if !g.InBounds(x, y) {
		panic(GridBoundsError)
	}
	g.data[x][y] = value
}

// Iterate over g, calling f at each site.
func (g *Grid) Iterate(f GridCallback) {
	for x := 0; x < g.Lx(); x++ {
		for y := 0; y < g.Ly(); y++ {
			f(x, y, g.Get(x, y))
		}
	}
}

// Return the number of active sites in g.
func (g *Grid) ActiveSites() int {
	count := 0
	checkSite := func(x, y int, value bool) {
		if value {
			count++
		}
	}
	g.Iterate(checkSite)
	return count
}

// Return the number of dimers in g, assuming pairing happens in x only.
// Assume that [0, 0] and [0, 1] are paired (this defines the pairing of all
// other dimers).
func (g *Grid) Dimers() int {
	count := 0
	checkSite := func(x, y int, value bool) {
		// ignore odd sites to avoid double-counting
		if x%2 == 1 {
			return
		}
		// if Lx is odd, the final site can't be in a dimer
		if x+1 == g.Lx() {
			return
		}
		// look at the site to the right: count if both are active
		if value && g.Get(x+1, y) {
			count++
		}
	}
	g.Iterate(checkSite)
	return count
}
