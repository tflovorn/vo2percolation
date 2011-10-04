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

// Return the length of the matrix represented by sym
func (sym *SymmetricMatrix) Length() int {
	return sym.length
}

// Return the value at row i, column j in sym.
func (sym *SymmetricMatrix) Get(i, j int) float64 {
	if i > sym.length || j > sym.length {
		panic("matrix access out of bounds")
	}
	if i > j {
		i, j = j, i
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
	if i > j {
		i, j = j, i
	}
	row, ok := sym.data[i]
	if !ok {
		// row doesn't exist yet, need to create it
		row = make(map[int]float64)
		row[j] = val
		sym.data[i] = row
	} else {
		sym.data[i][j] = val
	}
}

// Return true if and only if sym and comp represent the same matrix.
func (sym *SymmetricMatrix) Equals(comp *SymmetricMatrix) bool {
	if sym.length != comp.length {
		return false
	}
	// Need to iterate over both matrices since we skip zeroed elements.
	for i, row := range sym.data {
		for j, val := range row {
			if comp.Get(i, j) != val {
				return false
			}
		}
	}
	for i, row := range comp.data {
		for j, val := range row {
			if sym.Get(i, j) != val {
				return false
			}
		}
	}
	return true
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

// Return a new SymmetricMatrix derived from sym with convert, which is a map
// from the indices in sym to the indices in the returned matrix.
// The length of the returned matrix is given by length.
// All rows of the returned matrix which are not given in convert are zeroed.
func (sym *SymmetricMatrix) ReconstructEmptyRows(convert map[int]int, length int) *SymmetricMatrix {
	retMatrix := NewSymmetricMatrix(length)
	for i, row := range sym.data {
		for j, val := range row {
			iRet, jRet := convert[i], convert[j]
			retMatrix.Set(iRet, jRet, val)
		}
	}
	return retMatrix
}

func InsertEmptyRows(orig [][]float64, convert map[int]int, length int) [][]float64 {
	// initialize the return slice to zeros
	retSlice := make([][]float64, length)
	for i := 0; i < length; i++ {
		retSlice[i] = make([]float64, length)
	}
	// iterate over possible new labels
	for i := 0; i < length; i++ {
		oldI, ok := convert[i]
		if ok {
			// row/col is nonzero, copy values from orig
			for j := 0; j < length; j++ {
				oldJ, ok := convert[j]
				if ok {
					// this element is nonzero
					val := orig[oldI][oldJ]
					retSlice[i][j] = val
				}
			}
		}
	}
	return retSlice
}

// Return an ordered slice of the eigenvalues of sym, and a slice of the
// eigenvectors in the same order.
func (sym *SymmetricMatrix) Eigensystem() ([]float64, [][]float64) {
	originalSize := sym.length
	reduced, convert := sym.RemoveEmptyRows()
	size := C.size_t(reduced.length)
	eigenvalues := C.gsl_vector_alloc(size)
	eigenvectors := C.gsl_matrix_alloc(size, size)
	matrix := reduced.toMatrix()
	work := C.gsl_eigen_symmv_alloc(size)
	err := C.gsl_eigen_symmv(matrix, eigenvalues, eigenvectors, work)
	if err != 0 {
		// handle it
	}
	goEigenvalues := vectorToSlice(eigenvalues)
	goEigenvectors := matrixColumnsToSlices(eigenvectors)
	C.gsl_vector_free(eigenvalues)
	C.gsl_matrix_free(eigenvectors)
	C.gsl_matrix_free(matrix)
	C.gsl_eigen_symmv_free(work)
	retEigenvectors := InsertEmptyRows(goEigenvectors, convert, originalSize)
	return goEigenvalues, retEigenvectors
}

// Return the GSL matrix representation of sym.
func (sym *SymmetricMatrix) toMatrix() *C.gsl_matrix {
	// start with a zeroed matrix
	size := C.size_t(sym.length)
	matrix := C.gsl_matrix_calloc(size, size)
	// iterate over sym, setting the corresponding elements in matrix
	for i, row := range sym.data {
		for j, val := range row {
			it, jt := C.size_t(i), C.size_t(j)
			dval := C.double(val)
			C.gsl_matrix_set(matrix, it, jt, dval)
			if i != j {
				C.gsl_matrix_set(matrix, jt, it, dval)
			}
		}
	}
	return matrix
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

func matrixColumnsToSlices(m *C.gsl_matrix) [][]float64 {
	vectors := [][]float64{}
	var i, j C.size_t
	for i = 0; i < m.size1; i++ {
		xs := []float64{}
		for j = 0; j < m.size2; j++ {
			xs = append(xs, float64(C.gsl_matrix_get(m, j, i)))
		}
		vectors = append(vectors, xs)
	}
	return vectors
}
