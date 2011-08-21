// Functions for calculating the energy of a given grid of atoms.
package percolation

// Give the energy corresponding to the given grid of atoms, without including
// the effects of electron motion.
func AtomicHamiltonian(grid *Grid, env *Environment) float64 {
	return env.Delta*float64(grid.ActiveSites()) - env.V*float64(grid.Dimers())
}
