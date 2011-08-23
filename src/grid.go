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
const DimerPartnerError = "Site has no dimer partner"

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

// Flip the grid value at (x, y).  Panic if (x, y) is out of bounds.
func (g *Grid) Toggle(x, y int) {
	if !g.InBounds(x, y) {
		panic(GridBoundsError)
	}
	g.data[x][y] = !g.data[x][y]
}

// Iterate over g, calling f at each site.
func (g *Grid) Iterate(f GridCallback) {
	for x := 0; x < g.Lx(); x++ {
		for y := 0; y < g.Ly(); y++ {
			f(x, y, g.Get(x, y))
		}
	}
}

// Return a pointer to a copy of g.
func (g *Grid) Copy() *Grid {
	copyData := [][]bool{}
	for x := 0; x < g.Lx(); x++ {
		nextColumn := []bool{}
		for y := 0; y < g.Ly(); y++ {
			nextColumn = append(nextColumn, g.Get(x, y))
		}
		copyData = append(copyData, nextColumn)
	}
	cg, err := NewGrid(copyData)
	// this shouldn't happen, g should have had the correct shape
	if err != nil {
		panic(err)
	}
	return cg
}

// Return the number of active sites in g.
func (g *Grid) ActiveSiteCount() int {
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
func (g *Grid) DimerCount() int {
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

// Return the site which can form a dimer with the given site at (x, y).
func (g *Grid) DimerPartner(x, y int) (int, int, os.Error) {
	if !g.InBounds(x, y) {
		panic(GridBoundsError)
	}
	// even site: parter is to the right (if it exists)
	if x%2 == 0 {
		// if Lx is odd, last site doesn't have a partner
		if x+1 == g.Lx() {
			return -1, -1, fmt.Errorf(DimerPartnerError)
		}
		return x + 1, y, nil
	}
	// odd site: partner is to the left
	return x - 1, y, nil
}

// How would flipping the site at (x, y) affect the number of dimers?
func (g *Grid) DimerChange(x, y int) int {
	thisSiteValue := g.Get(x, y)
	xPartner, yPartner, err := g.DimerPartner(x, y)
	if err != nil {
		if err.String() == DimerPartnerError {
			return 0
		}
		panic("unexpected error from DimerPartner")
	}
	partnerSiteValue := g.Get(xPartner, yPartner)
	if thisSiteValue && partnerSiteValue {
		return -1
	}
	if !thisSiteValue && partnerSiteValue {
		return 1
	}
	return 0
}

// Return a PointSet containing all active sites in the grid.
func (g *Grid) ActiveSites() *PointSet {
	ps := NewPointSet(g.Lx(), g.Ly())
	addSite := func(x, y int, value bool) {
		// only add active sites
		if value {
			ps.Add(x, y)
		}
	}
	g.Iterate(addSite)
	return ps
}

// Return a slice containing each cluster of active sites on the grid.
func (g *Grid) AllClusters() []*PointSet {
	clusters := []*PointSet{}
	unexplored := g.ActiveSites()
	for unexplored.Size() > 0 {
		x, y := unexplored.Point()
		thisCluster := g.Cluster(x, y)
		for _, point := range thisCluster.Elements() {
			xr, yr := point[0], point[1]
			unexplored.Remove(xr, yr)
		}
		clusters = append(clusters, thisCluster)
	}
	return clusters
}

// Return the cluster at (x, y).
func (g *Grid) Cluster(x, y int) *PointSet {
	ps := NewPointSet(g.Lx(), g.Ly())
	g.clusterHelper(x, y, ps)
	return ps
}

// Add all members of the cluster at (x, y) to ps.
func (g *Grid) clusterHelper(x, y int, ps *PointSet) {
	if ps == nil {
		panic("must initialize ps in clusterHelper")
	}
	// base case: don't go to inactive or already-seen sites
	if ps.Contains(x, y) || !g.Get(x, y) {
		return
	}
	// haven't seen this active site yet: add it and try its neighbors
	ps.Add(x, y)
	ns := g.Neighbors(x, y)
	for _, point := range ns {
		xn, yn := point[0], point[1]
		if !ps.Contains(xn, yn) {
			g.clusterHelper(xn, yn, ps)
		}
	}
}

// Return a slice containing all neighbors of the given point.
// Only counts nearest neighbors for now - could also count next nearest.
func (g *Grid) Neighbors(x, y int) [][]int {
	ns := [][]int{}
	// left
	if x > 0 {
		ns = append(ns, []int{x - 1, y})
	}
	// right
	if x < g.Lx()-1 {
		ns = append(ns, []int{x + 1, y})
	}
	// down
	if y > 0 {
		ns = append(ns, []int{x, y - 1})
	}
	// up
	if y < g.Ly()-1 {
		ns = append(ns, []int{x, y + 1})
	}
	return ns
}
