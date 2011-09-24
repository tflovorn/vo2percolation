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

// Return a zeroed LxL symmetric matric.
func NewSymmetricMatrix(L int) *SymmetricMatrix {
	data := make(map[int]map[int]float64)
	sym := new(SymmetricMatrix)
	sym.length = L
	sym.data = data
	return sym
}

// Return the value at row i, column j in sym.
func (sym *SymmetricMatrix) Get(i, j int) float64 {
	if i > sym.length || j > sym.length {
		panic("matrix access out of bounds")
	}
	row, ok := sym.data[i]
	if !ok {
		return 0
	}
	val, ok := row[j]
	if !ok {
		return 0
	}
	return val
}

// Set the value at row i, column j in sym to val.
func (sym *SymmetricMatrix) Set(i, j int, val float64) {
	if i > sym.length || j > sym.length {
		panic("matrix access out of bounds")
	}
	row, ok := sym.data[i]
	if !ok {
		// row doesn't exist yet, need to create it
		row = make(map[int]float64)
		sym.data[i] = row
	}
	row[j] = val
}

// Return a new SymmetricMatrix without the empty rows (and columns) in sym.
// The map returned converts from row indices in the returned matrix to row
// indices in the original matrix.
func (sym *SymmetricMatrix) RemoveEmptyRows() (*SymmetricMatrix, map[int]int) {
	return nil, nil
}

// Return an ordered slice of the eigenvalues of sym, and a slice of the
// eigenvectors in the same order.
func (sym *SymmetricMatrix) Eigensystem() ([]float64, [][]float64) {
	return nil, nil
}

// Return the GSL matrix representation of sym.
func (sym *SymmetricMatrix) toMatrix() *C.gsl_matrix {
	// start with a zeroed matrix (m)

	// iterate over rows in sym (row = i)

	// for each column j in the row i: val = sym(i, j)
	// if diagonal: m(i,j) = val
	// if not diagonal: m(i,j) = val and m(j, i) = val
	return nil
}

func vectorToSlice(v *C.gsl_vector) []float64 {
	xs := []float64{}
	var i C.size_t
	for i = 0; i < v.size; i++ {
		xs = append(xs, float64(C.gsl_vector_get(v, i)))
	}
	return xs
}

func matrixToSlices(m *C.gsl_matrix) [][]float64 {
	vectors := [][]float64{}
	var i, j C.size_t
	for i = 0; i < m.size1; i++ {
		xs := []float64{}
		for j = 0; j < m.size2; j++ {
			xs = append(xs, float64(C.gsl_matrix_get(m, i, j)))
		}
		vectors = append(vectors, xs)
	}
	return vectors
}
