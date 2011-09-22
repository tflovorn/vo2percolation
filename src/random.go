// Functions for producing random values.
package percolation

import (
	"rand"
	"time"
)

// Return a random boolean value.
// --- does this produce a sufficiently random sample? ---
func RandomBool() bool {
	// seed to current time
	rand.Seed(time.Nanoseconds())
	// pick 0 or 1 randomly
	i := rand.Intn(2)
	// sanity check
	if i < 0 || i > 1 {
		panic("random integer outside expected bounds")
	}
	// convert to bool
	if i == 0 {
		return false
	}
	return true
}

// Return a pair of random integer values whose maxima are given by topX, topY.
func RandomPoint(g *Grid) Point {
	// one seed for both rand.Intn calls
	rand.Seed(time.Nanoseconds())
	randX := rand.Intn(g.Lx())
	randY := rand.Intn(g.Ly())
	return Point{randX, randY}
}

// Return a random float64 in the range [0.0, 1.0).
func RandomFloat() float64 {
	rand.Seed(time.Nanoseconds())
	return rand.Float64()
}
