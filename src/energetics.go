// Functions for calculating energies based on a Grid and Environment.
package percolation

import "math"

type Energetics struct {
	env Environment
}

func NewEnergetics(env Environment) {
	e := new(Energetics)
	e.env = env
	return e
}

// Environment access functions
func (e *Energetics) Beta() float64 {
	return env.Beta
}

func (e *Energetics) Delta() float64 {
	return env.Delta
}

func (e *Energetics) V() float64 {
	return env.V
}

// Give the energy corresponding to the given grid of atoms without including
// the effect of the electrons.
func (e *Energetics) AtomicHamiltonian(grid *Grid) float64 {
	activeSites := float64(grid.ActiveSiteCount())
	dimers := float64(grid.DimerCount())
	return e.Delta()*activeSites - e.V()*dimers
}

// Probability to accept a Monte Carlo state when energy changed by the given
// positive value.
func (e *Energetics) BoltzmannFactor(energy float64) float64 {
	return math.Exp(e.Beta() * energy)
}
