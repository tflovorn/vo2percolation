// 2D centered rectangular (rhombic) lattice where each site is characterized 
// by a boolean value.  Site labeling is shown at:
// http://dl.dropbox.com/u/41859823/vo2/lattice.pdf
// Sites where the boolean is true are 'active' and can form clusters.
// Clusters are connected by neighbors in the dimer-forming (x) direction and
// in the diagonal directions.
package percolation

import (
	"os"
	"fmt"
	"rand"
	"time"
	"math"
)

const GridShapeError = "Grid data must be rectangular and contain at least one point"
const GridBoundsError = "Grid point (%d, %d) out of bounds"
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

// Construct an empty grid with the given dimensions.
func NewGridWithDims(Lx, Ly int) *Grid {
	data := [][]bool{}
	for x := 0; x < Lx; x++ {
		newData := []bool{}
		for y := 0; y < Ly; y++ {
			newData = append(newData, false)
		}
		data = append(data, newData)
	}
	g, err := NewGrid(data)
	if err != nil {
		panic("NewGridWithDims failed")
	}
	return g
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

// Return true if the point (x, y) is within the grid boundaries; return false
// otherwise.
func (g *Grid) InBounds(x, y int) bool {
	if x < 0 || y < 0 || x >= g.Lx() || y >= g.Ly() {
		return false
	}
	return true
}

// If the point (x, y) is not in the grid boundaries, panic. Return true if 
// the point is within bounds.
func (g *Grid) CheckBounds(x, y int) bool {
	if !g.InBounds(x, y) {
		panic(fmt.Sprintf(GridBoundsError, x, y))
	}
	return true
}

// Get the grid value at (x, y). Panic if (x, y) is out of bounds.
func (g *Grid) Get(x, y int) bool {
	g.CheckBounds(x, y)
	return g.data[x][y]
}

// Set the grid value at (x, y). Panic if (x, y) is out of bounds.
func (g *Grid) Set(x, y int, value bool) {
	g.CheckBounds(x, y)
	g.data[x][y] = value
}

// Flip the grid value at (x, y).  Panic if (x, y) is out of bounds.
func (g *Grid) Toggle(x, y int) {
	g.CheckBounds(x, y)
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
	g.CheckBounds(x, y)
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

// Convert from 1D map keys to x-y grid coordinates.  Panics if the key is not
// on the grid.
func (g *Grid) convertFrom1D() func(int) Point {
	return func(key int) Point {
		lx, ly := g.Lx(), g.Ly()
		if key < 0 || key > lx*ly {
			panic("1D point conversion out of bounds")
		}
		y := int(math.Floor(float64(key) / float64(lx)))
		x := key - lx*y
		return Point{x, y}
	}
}

// Convert the 2D x-y coordinates in the Lx by Ly discrete grid to a single
// integer, useful as a map key.  Panics if (x,y) is not on the grid.
func (g *Grid) convertTo1D() func(Point) int {
	return func(p Point) int {
		lx, ly := g.Lx(), g.Ly()
		x, y := p.X(), p.Y()
		if x < 0 || y < 0 || x > lx || y > ly {
			panic("1D point conversion out of bounds")
		}
		return lx*y + x
	}
}

// Return a new PointSet on g.
func (g *Grid) PointSet() *PointSet {
	return NewPointSet(g.convertFrom1D(), g.convertTo1D())
}

// Return a PointSet containing all active sites in the grid.
func (g *Grid) ActiveSites() *PointSet {
	ps := g.PointSet()
	addSite := func(x, y int, value bool) {
		// only add active sites
		if value {
			ps.Add(Point{x, y})
		}
	}
	g.Iterate(addSite)
	return ps
}

// Return a slice containing each cluster of active sites on the grid.
// This is -not- the accepted method for doing this search.
// That method is the 'Hoshen-Kopelman algorithm'. See:
// http://www.ocf.berkeley.edu/~fricke/projects/hoshenkopelman/hoshenkopelman.html
func (g *Grid) AllClusters() []*PointSet {
	clusters := []*PointSet{}
	unexplored := g.ActiveSites()
	for unexplored.Size() > 0 {
		p := unexplored.Point()
		thisCluster := g.Cluster(p)
		for _, clusterPoint := range thisCluster.Elements() {
			unexplored.Remove(clusterPoint)
		}
		clusters = append(clusters, thisCluster)
	}
	return clusters
}

// Return the largest cluster on the grid.
func (g *Grid) LargestCluster() *PointSet {
	clusters := g.AllClusters()
	var max *PointSet = nil
	maxSize := 0
	for _, cluster := range clusters {
		thisSize := cluster.Size()
		if max == nil || thisSize > maxSize {
			max = cluster
			maxSize = thisSize
		}
	}
	return max
}

// Return the cluster at (x, y).
func (g *Grid) Cluster(p Point) *PointSet {
	ps := g.PointSet()
	g.clusterHelper(p, ps)
	return ps
}

// Add all members of the cluster at (x, y) to ps.
func (g *Grid) clusterHelper(start Point, ps *PointSet) {
	if ps == nil {
		panic("must initialize ps in clusterHelper")
	}
	// base case: don't go to inactive or already-seen sites
	if ps.Contains(start) || !g.Get(start.X(), start.Y()) {
		return
	}
	// haven't seen this active site yet: add it and try its neighbors
	ps.Add(start)
	ns := g.Neighbors(start)
	for _, point := range ns {
		if !ps.Contains(point) {
			g.clusterHelper(point, ps)
		}
	}
}

// Return a slice containing all neighbors of the given point.
// A site which isn't on a boundary has 6 neighbors: 2 in the dimer direction
// and 4 in the diagonal directions.
func (g *Grid) Neighbors(p Point) []Point {
	ns := []Point{}
	ns = append(ns, g.DimerNeighbors(p)...)
	ns = append(ns, g.DiagNeighbors(p)...)
	return ns
}

// Return the dimer-direction neighbors of the given point.
func (g *Grid) DimerNeighbors(p Point) []Point {
	ns := []Point{}
	xmax := g.Lx() - 1
	x, y := p.X(), p.Y()
	// left
	if x > 0 {
		ns = append(ns, Point{x - 1, y})
	}
	// right
	if x < xmax {
		ns = append(ns, Point{x + 1, y})
	}
	return ns
}

// Return the diagonal neighbors of the given point.
func (g *Grid) DiagNeighbors(p Point) []Point {
	ns := []Point{}
	xmax, ymax := g.Lx()-1, g.Ly()-1
	x, y := p.X(), p.Y()
	// diagonal neighbor labeling depends on parity of y
	if y%2 == 0 {
		// down-right
		if y > 0 {
			ns = append(ns, Point{x, y - 1})
		}
		// up-right
		if y < ymax {
			ns = append(ns, Point{x, y + 1})
		}
		// down-left
		if x > 0 && y > 0 {
			ns = append(ns, Point{x - 1, y - 1})
		}
		// up-left
		if x > 0 && y < ymax {
			ns = append(ns, Point{x - 1, y + 1})
		}
	} else {
		// odd y ==> know that y > 0
		// down-right
		if x < xmax {
			ns = append(ns, Point{x + 1, y - 1})
		}
		// up-right
		if x < xmax && y < ymax {
			ns = append(ns, Point{x + 1, y + 1})
		}
		// down-left
		ns = append(ns, Point{x, y - 1})
		// up-left
		if y < ymax {
			ns = append(ns, Point{x, y + 1})
		}
	}
	return ns
}
