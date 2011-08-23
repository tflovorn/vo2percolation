package percolation

import (
	"testing"
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
			if grid.Get(x, y) != val {
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
	if grid.Get(0, 0) != defaultValue {
		t.Fatalf("grid holds incorrect value")
	}
	grid.Set(0, 0, !defaultValue)
	if grid.Get(0, 0) == defaultValue {
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
	knownCluster1, knownCluster2 := NewPointSet(grid.Lx(), grid.Ly()), NewPointSet(grid.Lx(), grid.Ly())
	knownCluster1.Add(0, 0)
	knownCluster2.Add(0, 2)
	knownCluster2.Add(1, 2)
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
