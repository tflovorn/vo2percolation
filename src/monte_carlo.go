// Monte Carlo simulation of VO2 system
package percolation

import (
	"os"
	"math"
	"fmt"
)

const MonteCarloValidateError = "Monte Carlo input parameters are invalid"

// Hold parameters necesarry to perform the simulation
type MonteCarlo struct {
	// Minimum value of random flip acceptance probability (must be > 0).
	etaMinimum float64
	// How many total steps should be taken in the simulation?
	totalSteps int
	// Number of steps to take between snapshots of the grid.
	// If recordInterval is <= 0, only snapshot the final grid.
	recordInterval int
}

// Data reported for each time step in the simulation
type MonteCarloOutput struct {
	ActiveSites, Dimers, LargestClusterSize int
	Grid *Grid // may be nil
}

// Create a new (input-validated) MonteCarlo with the given parameters.
func NewMonteCarlo(etaMinimum float64, totalSteps, recordInterval int) (*MonteCarlo, os.Error) {
	mc := new(MonteCarlo)
	mc.etaMinimum = etaMinimum
	mc.totalSteps = totalSteps
	mc.recordInterval = recordInterval
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


// Run a simulation, starting from a random grid with dimensions (Lx, Ly) and
// taking steps equal to mc.totalSteps.  Return a slice containg each recorded
// grid.  May also want to return a slice of the times when each grid was
// recorded.
func (mc *MonteCarlo) Simulate(e *Energetics, Lx, Ly int) ([]*MonteCarloOutput, os.Error) {
	outputList := []*MonteCarloOutput{}
	// estimate starting number of active sites
	expectedActive := int(float64(Lx * Ly) * e.Boltzmann(e.Delta()))
	// generate the initial grid
	grid, err := RandomConstrainedGrid(Lx, Ly, expectedActive)
	if err != nil {
		return nil, err
	}
	// Monte Carlo loop
	for time := 0; time < mc.totalSteps; time++ {
		thisOutput := new(MonteCarloOutput)
		// log grid if it's the right time to
		if time % mc.recordInterval == 0 {
			thisOutput.Grid = grid.Copy()
		}
		// record the quantities we want to know for each configuration
		thisOutput.ActiveSites = grid.ActiveSiteCount()
		thisOutput.Dimers = grid.DimerCount()
		thisOutput.LargestClusterSize = grid.LargestCluster().Size()
		outputList = append(outputList, thisOutput)
		// try to perturb the grid
		mc.Step(e, grid) // could record failure/success here
	}
	return outputList, nil
}
