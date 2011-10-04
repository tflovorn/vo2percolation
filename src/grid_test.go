package percolation

import (
	"testing"
	/* for performance test
	"math"
	"time"
	"fmt"
	*/
)

var defaultData [][]bool = [][]bool{[]bool{true, false, true}, []bool{false, false, true}}

// Is a grid created with the proper dimensions and values?
func TestGridCreation(t *testing.T) {
	data := defaultData
	grid, err := NewGrid(data)
	if err != nil {
		t.Fatal(err)
	}
	if grid.Lx() != len(data) || grid.Ly() != len(data[0]) {
		t.Fatalf("grid does not have the correct dimensions")
	}
	for x, row := range data {
		for y, val := range row {
			if grid.Get(Point{x, y}) != val {
				t.Fatalf("grid holds incorrect value")
			}
		}
	}
}

// When a grid point is set, does it take that value?
func TestGridSet(t *testing.T) {
	defaultValue := false
	data := [][]bool{[]bool{defaultValue}}
	grid, err := NewGrid(data)
	if err != nil {
		t.Fatal(err)
	}
	p := Point{0, 0}
	if grid.Get(p) != defaultValue {
		t.Fatalf("grid holds incorrect value")
	}
	grid.Set(p, !defaultValue)
	if grid.Get(p) == defaultValue {
		t.Fatalf("set failed to change value")
	}
}

// Does Grid correctly count the number of active sites and dimers?
func TestGridSiteCounting(t *testing.T) {
	grid, err := NewGrid(defaultData)
	if err != nil {
		t.Fatal(err)
	}
	if grid.ActiveSiteCount() != 3 {
		t.Fatalf("grid reports incorrect number of active sites")
	}
	if grid.DimerCount() != 1 {
		t.Fatalf("grid reports incorrect number of dimers")
	}
}

// Does AllClusters return the correct clusters?
// Does LargestCluster pick the right cluster?
func TestGridClusters(t *testing.T) {
	grid, err := NewGrid(defaultData)
	if err != nil {
		t.Fatal(err)
	}
	knownCluster1, knownCluster2 := grid.PointSet(), grid.PointSet()
	knownCluster1.Add(Point{0, 0})
	knownCluster2.Add(Point{0, 2})
	knownCluster2.Add(Point{1, 2})
	clusters := grid.AllClusters()
	for _, ps := range clusters {
		if !ps.Equals(knownCluster1) && !ps.Equals(knownCluster2) {
			t.Fatalf("unexpected cluster")
		}
	}
	largest := grid.LargestCluster()
	if !largest.Equals(knownCluster2) {
		t.Fatalf("incorrect largest cluster")
	}
}

// A RandomConstrainedGrid should start with the number of active sites we
// tell it to have.
func TestRandomConstrainedGridCreation(t *testing.T) {
	activeSites := 128
	L := 64
	grid, err := RandomConstrainedGrid(L, L, activeSites)
	if err != nil {
		t.Fatal(err)
	}
	if grid.ActiveSiteCount() != activeSites {
		t.Fatalf("RandomConstrainedGrid did not start with the requested number of active sites")
	}
}

/*
// The performance of AllClusters should scale linearly with the number of
// sites in the grid.  Does it?
// (lots of code stolen from PointSetPerformance test, may want a common
// framework for grid benchmarks.)
func TestAllClustersPerformance(t *testing.T) {
	failureRatio := 1.5
	startN := 1024
	scaleN := 4
	N, firstTimeRatio := startN, 0.0
	minIters := 5
	minStopN := startN * int(math.Ceil(math.Pow(float64(scaleN), float64(minIters-1))))
	// Percolation threshold is ~0.59 on square lattice.
	// What is the threshold for our restricted rhombic lattice?
	activeFraction := 0.5
	for totalElapsed := int64(0); totalElapsed < 5e8 || N <= minStopN; N *= scaleN {
		initTime := time.Nanoseconds()
		// Initialize the random grid with N sites.
		// A constant fraction of these sites are active.
		activeN := int(activeFraction * float64(N))
		L := int(math.Ceil(math.Sqrt(float64(N))))
		grid, err := RandomConstrainedGrid(L, L, activeN)
		if err != nil {
			t.Fatal(err)
		}
		// How long does it take to find all the clusters?
		// Repeat the process a constant number of times.
		repeatCount := 1
		retrieveStartTime := time.Nanoseconds()
		var clusters []*PointSet
		for i := 0; i < repeatCount; i++ {
			clusters = grid.AllClusters()
		}
		endTime := time.Nanoseconds()
		retrieveTime := endTime - retrieveStartTime
		totalElapsed = endTime - initTime
		if N == startN {
			firstTimeRatio = float64(retrieveTime) / float64(N)
			fmt.Println(N, firstTimeRatio, len(clusters))
		} else {
			retrieveTimeRatio := float64(retrieveTime) / float64(N)
			fmt.Println(N, retrieveTimeRatio, len(clusters))
			if retrieveTimeRatio > failureRatio*firstTimeRatio {
				t.Errorf("KNOWN BUG: time to run Grid.AllClusters() is scaling faster than O(n) wrt Grid size")
			}
		}
	}
}
*/
