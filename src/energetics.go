// Functions for calculating energies based on a Grid and Environment.
package percolation

// Give the energy corresponding to the given grid of atoms without including
// the effect of the electrons.
func AtomicHamiltonian(grid *Grid, env *Environment) float64 {
	activeSites := float64(grid.ActiveSites())
	dimers := float64(grid.Dimers())
	return env.Delta*activeSites - env.V*dimers
}
