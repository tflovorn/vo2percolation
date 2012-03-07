package vo2percolation

import (
	"testing"
	"math"
	"fmt"
)

func TestSiteFlipEnergyKnown(t *testing.T) {
	var energyDefaultGrid [][]bool = [][]bool{[]bool{true, false, true}, []bool{false, false, true}}
	grid, err := NewGrid(energyDefaultGrid)
	if err != nil {
		t.Fatal(err)
	}
	Delta, V, Beta := 1.0, 0.5, 1.0
	data := fmt.Sprintf("{\"Delta\":%f, \"V\":%f, \"Beta\":%f}", Delta, V, Beta)
	env, err := EnvironmentFromString(data)
	if err != nil {
		t.Fatal(err)
	}
	e := NewEnergetics(*env)
	noDimerDeactivate := e.SiteFlipEnergy(grid, Point{0, 0})
	if noDimerDeactivate != -e.Delta() {
		t.Fatalf("incorrect result from SiteFlipEnergy (noDimerDeactivate)")
	}
	withDimerDeactivate := e.SiteFlipEnergy(grid, Point{0, 2})
	if withDimerDeactivate != -e.Delta()+e.V() {
		t.Fatalf("incorrect result from SiteFlipEnergy (withDimerDeactivate)")
	}
	noDimerActivate := e.SiteFlipEnergy(grid, Point{0, 1})
	if noDimerActivate != e.Delta() {
		t.Fatalf("incorrect result from SiteFlipEnergy (noDimerActivate)")
	}
	withDimerActivate := e.SiteFlipEnergy(grid, Point{1, 0})
	if withDimerActivate != e.Delta()-e.V() {
		t.Fatalf("incorrect result from SiteFlipEnergy (withDimerActivate)")
	}
}

func TestElectronHamiltonian4x4(t *testing.T) {
	// initialize config data
	env, err := EnvironmentFromFile("default.json")
	if err != nil {
		t.Fatal(err)
	}
	e := NewEnergetics(*env)
	// initialize the grid (2x2, all active)
	grid := NewGridWithDims(2, 2)
	activate := func(p Point, val bool) {
		grid.Set(p, true)
	}
	grid.Iterate(activate)
	H_el := e.ElectronHamiltonian(grid)
	alpha_evals, _ := H_el[0].Eigensystem()
	beta_evals, _ := H_el[1].Eigensystem()
	eps := 1e-12
	neq := func(x float64, y float64) bool {
		return math.Fabs(x-y) > eps
	}
	expected_alpha_evals := []float64{2, 0, 2, 0}
	sq17 := math.Sqrt(17)
	expected_beta_evals := []float64{0.5 * (1 - sq17), 1, 0.5 * (1 + sq17), 2.0}
	for i := 0; i < 4; i++ {
		if neq(alpha_evals[i], expected_alpha_evals[i]) || neq(beta_evals[i], expected_beta_evals[i]) {
			t.Fatalf("encountered unexpected eigenvalue")
		}
	}
}

func TestFermiEnergy2x2(t *testing.T) {
	eps := 1e-12
	neq := func(x float64, y float64) bool {
		return math.Fabs(x-y) > eps
	}
	// initialize config data
	env, err := EnvironmentFromFile("default.json")
	if err != nil {
		t.Fatal(err)
	}
	e := NewEnergetics(*env)
	// initialize the grid (2x2, all active)
	grid := NewGridWithDims(2, 2)
	activate := func(p Point, val bool) {
		grid.Set(p, true)
	}
	grid.Iterate(activate)
	fermi, err := e.FermiEnergy(grid, 4)
	if err != nil {
		t.Fatal(err)
	}
	if neq(0.0, fermi) {
		t.Fatalf("unexpected value for Fermi energy")
	}
}

func TestFindMu(t *testing.T) {
	eps := 1e-8
	neq := func(x float64, y float64) bool {
		return math.Fabs(x-y) > eps
	}
	// initialize config data
	env, err := EnvironmentFromFile("default.json")
	if err != nil {
		t.Fatal(err)
	}
	e := NewEnergetics(*env)
	// initialize the grid (2x2, all active)
	grid := NewGridWithDims(2, 2)
	activate := func(p Point, val bool) {
		grid.Set(p, true)
	}
	grid.Iterate(activate)
	mu, err := e.FindMu(grid, 4)
	if err != nil {
		t.Fatal(err)
	}
	expectedMu := -0.45589897859312906
	if neq(mu, expectedMu) {
		t.Fatalf("unexpected value of mu")
	}
}
