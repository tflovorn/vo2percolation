package vo2percolation

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var mc_scaling *bool = flag.Bool("mc_scaling", false, "Run Monte Carlo scaling test")

// Want to know how the Monte Carlo simulation scales with grid size.
// Procedure: use an LxL grid, with L steadily increasing. Run the simulation
// for a constant number of steps. Record the time to execute mc.Simulate().
// Relevant outputs: L, execution time, execution time / L^2.
// Only run if the command-line flag mcbench is present.
func TestMonteCarloScaling(t *testing.T) {
	// should we run this?
	flag.Parse()
	if !*mc_scaling {
		return
	}
	// going to run the test; start setup
	etaMinimum := 1e-12
	totalSteps := 10000
	recordInterval := 0
	maxL := 64
	env, err := EnvironmentFromFile("default.json")
	if err != nil {
		t.Fatal(err)
	}
	energ := NewEnergetics(*env)
	mc, err := NewMonteCarlo(etaMinimum, totalSteps, recordInterval)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("L, execTime (ms), execTime / L^2")
	// setup is finished
	for L := 8; L <= maxL; L += 8 {
		execTime, err := TimeMonteCarloStep(energ, mc, L)
		if err != nil {
			t.Fatal(err)
		}
		msTime := execTime / 1.0e6
		fmt.Println(L, msTime, msTime/float64(L*L))
	}
}

func TimeMonteCarloStep(energ *Energetics, mc *MonteCarlo, L int) (float64, error) {
	initTime := time.Now()
	_, err := mc.Simulate(energ, L, L)
	if err != nil {
		return -1, err
	}
	elapsedTime := time.Now().Sub(initTime).Seconds()
	return elapsedTime, nil
}
