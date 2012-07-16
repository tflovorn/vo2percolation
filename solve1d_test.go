package vo2percolation

import (
	"math"
	"testing"
)

func TestSolve1DLinear(t *testing.T) {
	m, b := 1.0, -1.0
	linear := func(x float64) float64 {
		return m*x + b
	}
	expectedRoot := -b / m
	eps := 1e-12
	root, err := Solve1D(linear, -1e9, 1e9, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	neq := func(x, y float64) bool {
		return math.Abs(x-y) > eps
	}
	if neq(root, expectedRoot) {
		t.Fatalf("unexpected value for root")
	}
}
