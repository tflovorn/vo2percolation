// Monte Carlo simulation of VO2 system
package percolation

import (
	"os"
	"math"
	"fmt"
)

const MonteCarloValidateError = "Monte Carlo input parameters are invalid"

type MonteCarlo struct {
	// Minimum value of random flip acceptance probability (must be > 0).
	etaMinimum float64
	// How many total steps should be taken in the simulation?
	totalSteps int
	// Number of steps to take between snapshots of the grid.
	// If recordInterval is <= 0, only snapshot the final grid.
	recordInterval int
}

// Build a MonteCarlo from the JSON file at filePath.
func MonteCarloFromFile(filePath string) (*MonteCarlo, os.Error) {
	mc := new(MonteCarlo)
	err := CopyFromFile(filePath, mc)
	if err != nil {
		return nil, err
	}
	if !mc.validate() {
		return nil, fmt.Errorf(MonteCarloValidateError)
	}
	return mc, nil
}

// Build a MonteCarlo from the given JSON string.
func MonteCarloFromString(jsonData string) (*MonteCarlo, os.Error) {
	mc := new(MonteCarlo)
	err := CopyFromString(jsonData, mc)
	if err != nil {
		return nil, err
	}
	if !mc.validate() {
		return nil, fmt.Errorf(MonteCarloValidateError)
	}
	return mc, nil
}

// Do the fields of mc have acceptable values?
func (mc *MonteCarlo) validate() bool {
	return mc.etaMinimum > 0 && mc.totalSteps > 0
}

// Make a random perturbation on the Grid g.  If this perturbation leads to a
// negative energy change, accept it.  If it leads to a positive energy change,
// accept it with a random probability.  Return true if and only if the
// perturbation is accepted.
func (mc *MonteCarlo) Step(e *Energetics, g *Grid) bool {
	// choose a random site
	xf, yf := RandomIntPair(g.Lx(), g.Ly())
	// calculate the energy change due to flipping (xf, yf)
	energyChange := e.SiteFlipEnergy(g, xf, yf)
	// going to lower energy: accept it
	if energyChange < 0 {
		g.Toggle(xf, yf)
		return true
	}
	// gaining energy: accept if eta + etaMinimum <= e^(-beta*energyChange)
	log_eta := math.Log(RandomFloat() + mc.etaMinimum)
	acceptFactor := e.LogBoltzmann(energyChange)
	if log_eta <= acceptFactor {
		g.Toggle(xf, yf)
		return true
	}
	return false
}

// Run a simulation, starting from a random grid and taking steps equal to
// mc.totalSteps.  Return a slice containg each recorded grid.
// May also want to return a slice of the times when each grid was recorded.
func (mc *MonteCarlo) Simulate(e *Energetics) []*Grid {
	return nil
}
