package percolation

import (
	"testing"
	"fmt"
)

func TestVectorSort(t *testing.T) {
	values := []float64{-5, 10, 0, -10, 5}
	vectors := [][]float64{[]float64{1.0}, []float64{2.0}, []float64{3.0}, []float64{4.0}, []float64{5.0}}
	sortedVals, sortedVecs, err := VectorSort(values, vectors)
	if err != nil {
		t.Fatal(err)
	}
	expectedVals := []float64{-10, -5, 0, 5, 10}
	expectedVecs := [][]float64{[]float64{4.0}, []float64{1.0}, []float64{3.0}, []float64{5.0}, []float64{2.0}}
	for i := 0; i < len(values); i++ {
		if sortedVals[i] != expectedVals[i] {
			fmt.Println("incorrect sorted value")
		}
		for j := 0; j < len(expectedVecs[i]); j++ {
			if sortedVecs[i][j] != expectedVecs[i][j] {
				fmt.Println("incorrect sorted vector")
			}
		}
	}
}
