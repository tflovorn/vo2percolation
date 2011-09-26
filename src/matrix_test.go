package percolation

import (
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
	newSym, _ := sym.RemoveEmptyRows()
	if newSym.length != 3 {
		t.Fatalf("unexpected length after removing empty rows")
	}
	if newSym.Get(0, 2) != val || newSym.Get(1, 1) != val || newSym.Get(2, 0) != val {
		t.Fatalf("unexpected value returned from sym.Get after RemoveEmptyRows()")
	}
}
