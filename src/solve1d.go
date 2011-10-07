// Interface to GSL one-dimensional function solver.
package percolation

/*
#cgo LDFLAGS: -lgsl
#include <gsl/gsl_errno.h>
#include <gsl/gsl_math.h>
#include <gsl/gsl_roots.h>

#define SOLVE1D_MAX_ITER 1024

extern double goEvaluateSolve1D(double, void*);
extern void goPackDataSolve1D(void*, double, int);

static double solve1D(double left, double right, double epsAbs, double epsRel, void* userdata) {
	int status, iter = 0;
	int converged = 1;
	const gsl_root_fsolver_type *T;
	gsl_root_fsolver *s;
	double r = (left + right) / 2.0;
	double x_lo = left, x_hi = right;
	gsl_function F;

	F.function = &goEvaluateSolve1D;
	F.params = userdata;
	T = gsl_root_fsolver_brent;
	s = gsl_root_fsolver_alloc(T);
	gsl_root_fsolver_set(s, &F, x_lo, x_hi);

	do {
		iter++;
		status = gsl_root_fsolver_iterate(s);
		r = gsl_root_fsolver_root(s);
		x_lo = gsl_root_fsolver_x_lower(s);
		x_hi = gsl_root_fsolver_x_upper(s);
		status = gsl_root_test_interval(x_lo, x_hi, epsAbs, epsRel);
	} while (status == GSL_CONTINUE && iter < SOLVE1D_MAX_ITER);

	gsl_root_fsolver_free(s);
	if (iter >= SOLVE1D_MAX_ITER || status != GSL_SUCCESS) {
		converged = 0;
	}
	goPackDataSolve1D(userdata, r, converged);
	return r;
}
*/
import "C"

import (
	"os"
	"fmt"
	"unsafe"
)

type dataSolve1D struct {
	f         func(float64) float64
	root      float64
	converged int
}

//export goEvaluateSolve1D
func goEvaluateSolve1D(x C.double, dataPtr unsafe.Pointer) C.double {
	data := (*dataSolve1D)(dataPtr)
	val := data.f(float64(x))
	return C.double(val)
}

//export goPackDataSolve1D
func goPackDataSolve1D(dataPtr unsafe.Pointer, root C.double, converged C.int) {
	data := (*dataSolve1D)(dataPtr)
	data.root = float64(root)
	data.converged = int(converged)
}

// Find the root of error bracketed by left and right to absolute precision
// epsAbs and relative precision epsRel.
func Solve1D(error func(float64) float64, left, right, epsAbs, epsRel float64) (float64, os.Error) {
	errLeft, errRight := error(left), error(right)
	if (errLeft > 0 && errRight > 0) || (errLeft < 0 && errRight < 0) {
		return 0.0, fmt.Errorf("left and right do not bracket a root")
	}
	data := &dataSolve1D{error, 0.0, 0}
	dataPtr := unsafe.Pointer(data)
	C.solve1D(C.double(left), C.double(right), C.double(epsAbs), C.double(epsRel), dataPtr)
	if data.converged == 0 {
		return 0.0, fmt.Errorf("failed to find root to desired accuracy")
	}
	return data.root, nil
}
