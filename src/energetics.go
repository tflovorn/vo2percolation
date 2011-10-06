// Functions for calculating energies based on a Grid and Environment.
package percolation

import (
	"math"
)

type Energetics struct {
	env Environment
}

func NewEnergetics(env Environment) *Energetics {
	e := new(Energetics)
	e.env = env
	return e
}

// Environment access functions
func (e *Energetics) Beta() float64 {
	return e.env.Beta
}

func (e *Energetics) Delta() float64 {
	return e.env.Delta
}

func (e *Energetics) V() float64 {
	return e.env.V
}

// Give the energy corresponding to the given grid of atoms without including
// the effect of the electrons.
func (e *Energetics) AtomicHamiltonian(g *Grid) float64 {
	activeSites := float64(g.ActiveSiteCount())
	dimers := float64(g.DimerCount())
	return e.Delta()*activeSites - e.V()*dimers
}

// Energy change corresponding to an active/inactive flip on g at (xf, yf).
// (only includes change due to the atomic Hamiltonian for now - should also
// include change in electron energy)
func (e *Energetics) SiteFlipEnergy(g *Grid, p Point) float64 {
	siteValue := g.Get(p)
	dimerChange := g.DimerChange(p)
	energyChange := 0.0
	if siteValue {
		energyChange -= e.Delta()
	} else {
		energyChange += e.Delta()
	}
	energyChange -= float64(dimerChange) * e.V()
	return energyChange
}

// Hamiltonian for the electrons on g. Assumes two orbitals, where electrons
// in one orbital move only in the dimer and diagonal directions. Electrons
// on the other orbital move in both directions. Neither orbital allows the
// electrons to move in the direction perpindicular to the dimer direction.
func (e *Energetics) ElectronHamiltonian(g *Grid) []*SymmetricMatrix {
	alpha := NewSymmetricMatrix(g.Lx() * g.Ly())
	beta := NewSymmetricMatrix(g.Lx() * g.Ly())
	activePoints := g.ActiveSites().Elements()
	convert := g.ConvertTo1D()
	for _, p := range activePoints {
		id := convert(p)
		// on-site energy
		alpha.Add(id, id, e.env.Epsilon_alpha)
		beta.Add(id, id, e.env.Epsilon_beta)
		// dimer-direction hopping
		nsDim := g.DimerNeighbors(p)
		for _, n := range nsDim {
			if g.Get(n) {
				nId := convert(n)
				// factors of 1/2 due to double counting
				alpha.Add(id, nId, -e.env.T_alpha/2)
				beta.Add(id, nId, -e.env.T_beta_dimer/2)
			}
		}
		// diagonal-direction hopping
		nsDiag := g.DiagNeighbors(p)
		for _, n := range nsDiag {
			if g.Get(n) {
				nId := convert(n)
				beta.Add(id, nId, -e.env.T_beta_diag/2)
			}
		}
	}
	return []*SymmetricMatrix{alpha, beta}
}

// Boltzmann factor.
func (e *Energetics) Boltzmann(energy float64) float64 {
	return math.Exp(-e.Beta() * energy)
}

func (e *Energetics) LogBoltzmann(energy float64) float64 {
	return -e.Beta() * energy
}
