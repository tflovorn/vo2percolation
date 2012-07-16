package vo2percolation

import (
	"math"
	"testing"
)

func TestBuildSymmetricMatrix(t *testing.T) {
	sym := NewSymmetricMatrix(7)
	val := 5.0
	sym.Set(4, 6, val)
	if sym.Get(4, 6) != val || sym.Get(6, 4) != val {
		t.Fatalf("unexpected value returned from sym.Get after sym.Set")
	}
}

func TestRemoveEmptyRowsSingle(t *testing.T) {
	sym := NewSymmetricMatrix(5)
	val := 5.0
	sym.Set(2, 2, val)
	newSym, _ := sym.RemoveEmptyRows()
	if newSym.length != 1 {
		t.Fatalf("unexpected length after removing empty rows")
	}
	if newSym.Get(0, 0) != val {
		t.Fatalf("unexpected value returned from sym.Get after RemoveEmptyRows()")
	}
}

func TestRemoveEmptyRowsDouble(t *testing.T) {
	sym := NewSymmetricMatrix(5)
	val := 5.0
	sym.Set(2, 2, val)
	sym.Set(1, 3, val)
	newSym, convert := sym.RemoveEmptyRows()
	if newSym.length != 3 {
		t.Fatalf("unexpected length after removing empty rows")
	}
	if newSym.Get(0, 2) != val || newSym.Get(1, 1) != val || newSym.Get(2, 0) != val {
		t.Fatalf("unexpected value returned from sym.Get after RemoveEmptyRows()")
	}
	rebuild := newSym.ReconstructEmptyRows(convert, 5)
	if !rebuild.Equals(sym) {
		t.Fatalf("reconstructed matrix differs from original")
	}
}

func TestInsertEmptyRows(t *testing.T) {
	reduced := [][]float64{[]float64{2.0, 1.0}, []float64{1.0, 2.0}}
	convert := map[int]int{0: 0, 1: 2}
	n := InsertEmptyRows(reduced, convert, 3)
	expected := [][]float64{[]float64{2.0, 0.0, 1.0}, []float64{0.0, 0.0, 0.0}, []float64{1.0, 0.0, 2.0}}
	for i := 0; i < len(n); i++ {
		for j := 0; j < len(n); j++ {
			if n[i][j] != expected[i][j] {
				t.Fatalf("incorrect insertion of empty rows")
			}
		}
	}
}

func TestEigensystem2x2(t *testing.T) {
	sym := NewSymmetricMatrix(2)
	sym.Set(0, 0, 2.0)
	sym.Set(1, 0, 1.0)
	sym.Set(1, 1, 2.0)
	vals, vs := sym.Eigensystem()
	eps := 1e-12
	neq := func(x float64, y float64) bool {
		return math.Abs(x-y) > eps
	}
	// this could fail if the order of returned eigenvalues changes
	if neq(vals[0], 3.0) || neq(vals[1], 1.0) {
		t.Fatalf("incorrect eigenvalue returned")
	}
	x := 1.0 / math.Sqrt(2.0)
	if neq(vs[0][0], x) || neq(vs[0][1], x) || neq(vs[1][0], -x) || neq(vs[1][1], x) {
		t.Fatalf("incorrect eigenvector returned")
	}
}

func TestEigensystem3x3WithZeros(t *testing.T) {
	sym := NewSymmetricMatrix(3)
	sym.Set(0, 0, 2.0)
	sym.Set(2, 0, 1.0)
	sym.Set(2, 2, 2.0)
	vals, vs := sym.Eigensystem()
	eps := 1e-12
	neq := func(x float64, y float64) bool {
		return math.Abs(x-y) > eps
	}
	// this could fail if the order of returned eigenvalues changes
	if neq(vals[0], 3.0) || neq(vals[1], 1.0) {
		t.Fatalf("incorrect eigenvalue returned")
	}
	x := 1.0 / math.Sqrt(2.0)
	if neq(vs[0][0], x) || neq(vs[0][2], x) || neq(vs[2][0], -x) || neq(vs[2][2], x) {
		t.Fatalf("incorrect eigenvector returned")
	}

}

func TestEigensystem7x7Really3x3(t *testing.T) {
	sym := NewSymmetricMatrix(7)
	sym.Set(1, 1, 1.0)
	sym.Set(1, 3, 1.0)
	sym.Set(1, 5, 1.0)
	sym.Set(3, 3, 2.0)
	sym.Set(3, 5, 1.0)
	sym.Set(5, 5, 2.0)
	vals, vs := sym.Eigensystem()
	eps := 1e-12
	neq := func(x float64, y float64) bool {
		return math.Abs(x-y) > eps
	}
	s3 := math.Sqrt(3)
	// Inspect eigenvalues:
	// this could fail if the order of returned eigenvalues changes.
	if neq(vals[0], 2-s3) || neq(vals[1], 2+s3) || neq(vals[2], 1.0) {
		t.Fatalf("incorrect eigenvalue returned")
	}
	// Inspect zero eigenvectors.
	for i := 0; i <= 6; i += 2 {
		for j := 0; j < 7; j++ {
			if neq(vs[i][j], 0.0) {
				t.Fatalf("unexpected nonzero eigenvector")
			}
		}
	}
	// Inspect nonzero eigenvectors.
	norms := []float64{math.Sqrt(2 * (3 + s3)), math.Sqrt(2 * (3 - s3)), math.Sqrt(2)}
	expected := [][]float64{[]float64{1.0 + s3, -1, -1}, []float64{-s3 + 1.0, -1, -1}, []float64{0, -1, 1}}
	for i := 0; i < 3; i++ {
		for j := 0; j < 7; j++ {
			var val float64
			if j%2 == 0 { // even elements are zeroed
				val = 0.0
			} else {
				val = expected[i][(j-1)/2]
			}
			comp := norms[i] * vs[2*i+1][j]
			if neq(val, comp) {
				t.Fatalf("unexpected eigenvector element")
			}
		}
	}
}
