package percolation

import (
	"testing"
	"fmt"
)

var energyDefaultGrid [][]bool = [][]bool{[]bool{true, false, true}, []bool{false, false, true}}

func TestSiteFlipEnergyKnown(t *testing.T) {
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
	noDimerDeactivate := e.SiteFlipEnergy(grid, 0, 0)
	if noDimerDeactivate != -e.Delta() {
		t.Fatalf("incorrect result from SiteFlipEnergy (noDimerDeactivate)")
	}
	withDimerDeactivate := e.SiteFlipEnergy(grid, 0, 2)
	if withDimerDeactivate != -e.Delta()+e.V() {
		t.Fatalf("incorrect result from SiteFlipEnergy (withDimerDeactivate)")
	}
	noDimerActivate := e.SiteFlipEnergy(grid, 0, 1)
	if noDimerActivate != e.Delta() {
		t.Fatalf("incorrect result from SiteFlipEnergy (noDimerActivate)")
	}
	withDimerActivate := e.SiteFlipEnergy(grid, 1, 0)
	if withDimerActivate != e.Delta()-e.V() {
		t.Fatalf("incorrect result from SiteFlipEnergy (withDimerActivate)")
	}
}
