// Matrix functions necessary for eigendecomposition of electron Hamiltonian
// Inspired by matrix implementation in go-gsl
// (https://bitbucket.org/fhs/go-gsl)
package percolation

// some of these included packages may not be necessary

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_math.h>
#include <gsl/gsl_vector.h>
#include <gsl/gsl_matrix.h>
#include <gsl/gsl_eigen.h>
*/
import "C"

type SymmetricMatrix struct {
	// dimensions (symmetric, so lx = ly)
	length int
	// Data is a map since many rows/cols may be empty.
	// first index: row, second index: column
	data map[int]map[int]float64
}

// Return a new SymmetricMatrix without the empty rows (and columns) in sym.
// The map returned maps from row indices in the returned matrix to row
// indices in the original matrix.
func (sym *SymmetricMatrix) RemoveEmptyRows() (*SymmetricMatrix, map[int]int) {
	return nil, nil
}

// Return an ordered slice of the eigenvalues of sym, and a slice of the
// eigenvectors in the same order.
func Eigensystem(sym *SymmetricMatrix) ([]float64, [][]float64) {
	return nil, nil
}

func symToMatrix(sym *SymmetricMatrix) *C.gsl_matrix {
	return nil
}

func vectorToSlice(v *C.gsl_vector) []float64 {
	return nil
}

func matrixToSlices(m *C.gsl_matrix) [][]float64 {
	return nil
}
