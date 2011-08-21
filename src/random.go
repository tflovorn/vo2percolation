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
	// generate a random integer in the range [0, 2) (i.e. 0 or 1)
	i := rand.Intn(2)
	if i < 0 || i > 1 {
		panic("random integer outside expected bounds")
	}
	if i == 0 {
		return false
	}
	return true
}
