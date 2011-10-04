package percolation

import (
	"testing"
	/* for performance test
	"math"
	"time"
	*/
)

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

/*
// Retrieving an arbitrary point from a PointSet should be O(1) with respect
// to the size of the PointSet.
func TestPointFromPointSetPerformance(t *testing.T) {
	startN := 1024
	scaleN := 4
	N, firstRetrieveTime := startN, int64(0)
	minIters := 3
	minStopN := startN * int(math.Ceil(math.Pow(float64(scaleN), float64(minIters-1))))
	for totalElapsed := int64(0); totalElapsed < 5e8 || N < minStopN; N *= scaleN {
		initTime := time.Nanoseconds()
		// Initialize the PointSet with elements.
		L := int(math.Ceil(math.Sqrt(float64(N))))
		g := NewGridWithDims(L, L)
		ps := g.PointSet()
		convert := g.ConvertFrom1D()
		for i := 0; i < N; i++ {
			ps.Add(convert(i))
		}
		// How long does it take to retrieve an arbitrary point?
		// Retrieve a constant number of them.
		pointCount := 1000
		retrieveStartTime := time.Nanoseconds()
		for i := 0; i < pointCount; i++ {
			ps.Point()
		}
		endTime := time.Nanoseconds()
		retrieveTime := endTime - retrieveStartTime
		totalElapsed = endTime - initTime
		if N == startN {
			firstRetrieveTime = retrieveTime
		} else {
			someFactor := int64(2)
			if retrieveTime > someFactor*firstRetrieveTime {
				t.Fatalf("time to run Point() method of PointSet is scaling faster than O(1) wrt PointSet size")
			}
		}
	}
}
*/
