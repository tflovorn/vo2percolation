package vo2percolation

import (
	"math"
	"testing"
)

// Examine the statistics of RandomBool.
// There is a nonzero (but small) chance for this to fail, controlled by the
// value of epsilon.
func TestRandomBoolIsRandom(t *testing.T) {
	epsilon := 2e-2 // greatest allowed relError
	repeatCount := int(math.Pow(2.0, 16.0))
	trueCount, falseCount := 0, 0
	for i := 0; i < repeatCount; i++ {
		val := RandomBool()
		if val {
			trueCount += 1
		} else {
			falseCount += 1
		}
	}
	difference := math.Abs(float64(trueCount - falseCount))
	relError := difference / float64(repeatCount)
	if relError > epsilon {
		t.Fatalf("RandomBool() produced large excess of true or false")
	}
}
