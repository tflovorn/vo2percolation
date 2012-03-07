// Functions for producing random values.
package vo2percolation

import (
	"rand"
	"time"
)

// To avoid calling rand.NewSource many times, r is a module-level variable.
// This strategy is unlikely to scale over multiple goroutines!
var r *rand.Rand = newRandom()

// Create a new random value generator.
func newRandom() *rand.Rand {
	src := rand.NewSource(time.Nanoseconds())
	return rand.New(src)
}

// Return a random boolean value.
// --- does this produce a sufficiently random sample? ---
func RandomBool() bool {
	// pick 0 or 1 randomly
	i := r.Intn(2)
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
	randX := r.Intn(g.Lx())
	randY := r.Intn(g.Ly())
	return Point{randX, randY}
}

// Return a random float64 in the range [0.0, 1.0).
func RandomFloat() float64 {
	return r.Float64()
}
