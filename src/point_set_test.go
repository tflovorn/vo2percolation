package percolation

import "testing"

func TestPointSetContains(t *testing.T) {
	L := 5
	g := NewGridWithDims(L, L)
	ps := g.PointSet()
	for i := 0; i < L; i++ {
		p := Point{i, i}
		ps.Add(p)
		if !ps.Contains(p) {
			t.Fatalf("ps does not contain element added to it")
		}
	}
	for i := 0; i < L; i++ {
		p := Point{i, i}
		ps.Remove(p)
		if ps.Contains(p) {
			t.Fatalf("ps contains element removed from it")
		}
	}

}

func TestPointSetElements(t *testing.T) {
	L := 5
	g := NewGridWithDims(L, L)
	ps := g.PointSet()
	for i := 0; i < L; i++ {
		ps.Add(Point{i, i})
	}
	elems := ps.Elements()
	for _, point := range elems {
		if !ps.Contains(point) {
			t.Fatalf("ps does not contain point in Elements()")
		}
	}
}
