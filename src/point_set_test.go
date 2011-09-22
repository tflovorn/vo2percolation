package percolation

import "testing"

func TestPointSetContains(t *testing.T) {
	L := 5
	g := NewGridWithDims(L, L)
	ps := g.PointSet()
	for i := 0; i < L; i++ {
		ps.Add(i, i)
		if !ps.Contains(i, i) {
			t.Fatalf("ps does not contain element added to it")
		}
	}
	for i := 0; i < L; i++ {
		ps.Remove(i, i)
		if ps.Contains(i, i) {
			t.Fatalf("ps contains element removed from it")
		}
	}

}

func TestPointSetElements(t *testing.T) {
	L := 5
	g := NewGridWithDims(L, L)
	ps := g.PointSet()
	for i := 0; i < L; i++ {
		ps.Add(i, i)
	}
	elems := ps.Elements()
	for _, point := range elems {
		x, y := point[0], point[1]
		if !ps.Contains(x, y) {
			t.Fatalf("ps does not contain point in Elements()")
		}
	}
}
