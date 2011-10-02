package percolation

import (
	"os"
	"fmt"
)

const EnvironmentValidateError = "Environment input parameters are invalid"

// Physical parameters
type Environment struct {
	Beta  float64 // inverse thermal energy 1 / (k_B * T)
	Delta float64 // energy cost of exciting an atom
	V     float64 // energy gained from exciting a dimer
	// hopping energies:
	t_alpha      float64 // dimer direction, a_1g orbital
	t_beta_dimer float64 // dimer direction, e_pi orbital
	t_beta_diag  float64 // diagonal direction, e_pi orbital
}

// Build an Environment from the JSON file at filePath.
func EnvironmentFromFile(filePath string) (*Environment, os.Error) {
	return buildEnvironment(func(env *Environment) os.Error {
		return CopyFromFile(filePath, env)
	})
}

// Build an Environment from the given JSON string.
func EnvironmentFromString(jsonData string) (*Environment, os.Error) {
	return buildEnvironment(func(env *Environment) os.Error {
		return CopyFromString(jsonData, env)
	})
}

// Build an environment using the given copy function.
func buildEnvironment(copier func(*Environment) os.Error) (*Environment, os.Error) {
	env := new(Environment)
	err := copier(env)
	if err != nil {
		return nil, err
	}
	if !env.validate() {
		return nil, fmt.Errorf(EnvironmentValidateError)
	}
	return env, nil
}

// Do the fields of env have acceptable values?
func (env *Environment) validate() bool {
	return env.Delta > 0 && env.V > 0 && env.Beta > 0
}
