Percolation model of VO2.

Requires the following packages to build on Mint Debian:
gsl-bin
libgsl0ldbl
libgsl0-dev

__Description of the model__

Ground state: atoms are all dimerized, which prevents long-range hopping. 
Atoms can be thermally excited to the metallic state.

Costs energy Delta to move an atom from dimer to metal.
If both atoms in a dimer are moved to the metal, gain back energy V.

Electrons hop within a percolation network of metallic atoms.
Hopping onto dimer sites is renormalized away, leading to nnn hopping within
 the metal.

First step: Ising-like Monte Carlo in atom excitation variables {nu(i)=0, 1}.
Study the formation of the percolation network without including electrons.

Diagonalizing the electron hopping Hamiltonian isn't trivial.  
Need a strategy to estimate whether a configuration of nu(i)'s is worth allowing
 before diagonalization.

__Notes__

Monte Carlo implementation inspired by [this one in Fortran](http://fraden.brandeis.edu/courses/phys39/simulations/Student%20Ising%20Swarthmore.pdf).

gofmt [pre-commit hook](http://golang.tumblr.com/post/439868556/git-precommit-hook-for-gofmt) is helpful for commits.
