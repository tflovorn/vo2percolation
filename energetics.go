// Functions for calculating energies based on a Grid and Environment.
package vo2percolation

import (
	"fmt"
	"math"
	"sort"
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

// Boltzmann factor.
func (e *Energetics) Boltzmann(energy float64) float64 {
	return math.Exp(-e.Beta() * energy)
}

func (e *Energetics) LogBoltzmann(energy float64) float64 {
	return -e.Beta() * energy
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

// Return a sorted list of the electronic energy levels. Each level has a
// degeneracy of two due to the Hamiltonian's spin invariance.
func (e *Energetics) ElectronEnergies(g *Grid) []float64 {
	H_el := e.ElectronHamiltonian(g)
	alpha_evals, _ := H_el[0].Eigensystem()
	beta_evals, _ := H_el[1].Eigensystem()
	energies := append(alpha_evals, beta_evals...)
	sort.Float64s(energies)
	return energies
}

// Determine the Fermi energy by filling the lowest available states with
// a number of particles equal to particleCount.
func (e *Energetics) FermiEnergy(g *Grid, particleCount int) (float64, error) {
	if particleCount <= 0 {
		return 0.0, fmt.Errorf("Fermi energy not defined for given number of particles")
	}
	energies := e.ElectronEnergies(g)
	// numOccupied is the number of occupied energy levels
	var numOccupied int
	if particleCount%2 == 0 {
		numOccupied = particleCount / 2
	} else {
		numOccupied = (particleCount + 1) / 2
	}
	return energies[numOccupied-1], nil
}

// Fermi distribution function.
func FermiDist(energy float64) float64 {
	return 1.0 / (math.Exp(energy) + 1)
}

// Return the electron number corresponding to the chemical potential mu.
func (e *Energetics) NumElectrons(energies []float64, mu float64) float64 {
	sum := 0.0
	for _, energy := range energies {
		// might want to make this a Kahan summation
		sum += 2.0 * FermiDist(energy-mu) // 2 for spin degeneracy
	}
	return sum
}

// Returns the error in the number of particles calculated from mu.
func (e *Energetics) NumElectronsError(g *Grid, particleCount int, mu float64) float64 {
	energies := e.ElectronEnergies(g)
	return float64(particleCount) - e.NumElectrons(energies, mu)
}

// Find the value of mu appropriate for the given number of particles.
func (e *Energetics) FindMu(g *Grid, particleCount int) (float64, error) {
	error := func(mu float64) float64 {
		return e.NumElectronsError(g, particleCount, mu)
	}
	// arbitrary end points; assume error(mu) is monotonic
	muMin := -100.0 * e.Delta()
	muMax := 100.0 * e.Delta()
	eps := 1e-9
	return Solve1D(error, muMin, muMax, eps, eps)
}
