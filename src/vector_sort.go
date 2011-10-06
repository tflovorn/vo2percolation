// Sort a group of vectors by keys given in a separate list by implementing
// sort.Interface and calling sort.Sort().
package percolation

import (
	"sort"
	"fmt"
	"os"
)

type VectorSorter struct {
	// values and vectors must have the same length
	values  []float64
	vectors [][]float64
}

func NewVectorSorter(values []float64, vectors [][]float64) (*VectorSorter, os.Error) {
	if len(values) != len(vectors) {
		return nil, fmt.Errorf("arguments to NewVectorSorter must have the same length")
	}
	return &VectorSorter{values, vectors}, nil
}

func (vs *VectorSorter) Len() int {
	return len(vs.values)
}

func (vs *VectorSorter) Less(i, j int) bool {
	return vs.values[i] < vs.values[j]
}

func (vs *VectorSorter) Swap(i, j int) {
	vs.values[i], vs.values[j] = vs.values[j], vs.values[i]
	vs.vectors[i], vs.vectors[j] = vs.vectors[j], vs.vectors[i]
}

func VectorSort(values []float64, vectors [][]float64) ([]float64, [][]float64, os.Error) {
	vs, err := NewVectorSorter(values, vectors)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(vs)
	return vs.values, vs.vectors, nil
}
