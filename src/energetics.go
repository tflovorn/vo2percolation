// Functions for calculating energies based on a Grid and Environment.
package percolation

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
func (e *Energetics) SiteFlipEnergy(g *Grid, xf, yf int) float64 {
	siteValue := g.Get(xf, yf)
	dimerChange := g.DimerChange(xf, yf)
	energyChange := 0.0
	if siteValue {
		energyChange -= e.Delta()
	} else {
		energyChange += e.Delta()
	}
	energyChange -= float64(dimerChange) * e.V()
	return energyChange
}

// Probability to accept a Monte Carlo state when energy changed by the given
// positive value.
func (e *Energetics) LogBoltzmann(energy float64) float64 {
	return -e.Beta() * energy
}
