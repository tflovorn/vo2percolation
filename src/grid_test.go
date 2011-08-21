package percolation

import (
	"testing"
)

// Is a grid created with the proper dimensions and values?
func TestGridCreation(t *testing.T) {
	data := [][]bool{[]bool{true, false, false}, []bool{false, true, false}}
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
