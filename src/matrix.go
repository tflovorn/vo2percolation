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

import "fmt"

type SymmetricMatrix struct {
	// dimensions (symmetric, so lx = ly)
	length int
	// Data is a map since many rows/cols may be empty.
	// first index: row, second index: column
	data map[int]map[int]float64
}

// Return a zeroed LxL symmetric matrix.
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
		otherRow, ok := sym.data[j]
		if !ok {
			return 0
		} else {
			val, ok := otherRow[i]
			if !ok {
				return 0
			}
			return val
		}
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
		row[j] = val
		sym.data[i] = row
	} else {
		sym.data[i][j] = val
		// need to keep symmetric data consistent
		if i != j {
			if otherRow, ok := sym.data[j]; ok {
				if _, ok := otherRow[i]; ok {
					otherRow[i] = val
				}
			}
		}
	}
}

// Return a new SymmetricMatrix without the empty rows (and columns) in sym.
// The map returned converts from row indices in the returned matrix to row
// indices in the original matrix.
func (sym *SymmetricMatrix) RemoveEmptyRows() (*SymmetricMatrix, map[int]int) {
	nonEmpty, convert := make([]bool, sym.length), make(map[int]int)
	// Build a list of non-empty rows.
	for i, row := range sym.data {
		for j, val := range row {
			if val != 0.0 {
				if i != j {
					nonEmpty[i] = true
					nonEmpty[j] = true
				} else {
					nonEmpty[i] = true
				}
			}
		}
	}
	// Count the number of non-empty rows and build the map from the old
	// indexing to the indexing without empty rows.
	iNew := 0
	for iOld, val := range nonEmpty {
		if val {
			convert[iNew] = iOld
			iNew++
		}
	}
	// Build the matrix without empty rows and columns.
	newMatrix := NewSymmetricMatrix(iNew)
	for i := 0; i < iNew; i++ {
		for j := i; j < iNew; j++ {
			iOld, jOld := convert[i], convert[j]
			val := sym.Get(iOld, jOld)
			if val != 0.0 {
				newMatrix.Set(i, j, val)
			}
		}
	}
	return newMatrix, convert
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

func (sym *SymmetricMatrix) String() string {
	out := ""
	for i := 0; i < sym.length; i++ {
		outRow := ""
		for j := 0; j < sym.length; j++ {
			outRow += fmt.Sprint(sym.Get(i, j)) + " "
		}
		out += outRow + "\n"
	}
	return out
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
