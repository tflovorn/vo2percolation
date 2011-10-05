package percolation

import (
	"testing"
	"math"
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
		return math.Fabs(x-y) > eps
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
		return math.Fabs(x-y) > eps
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
