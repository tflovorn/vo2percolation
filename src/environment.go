package percolation

import "os"

// Physical parameters
type Environment struct {
	Delta float64 // energy cost of exciting an atom
	V     float64 // energy gained from exciting a dimer
	Beta  float64 // inverse thermal energy 1 / (k_B * T)
}

// Build an Environment from the JSON file at filePath.
func EnvironmentFromFile(filePath string) (*Environment, os.Error) {
	env := new(Environment)
	err := CopyFromFile(filePath, env)
	if err != nil {
		return nil, err
	}
	return env, nil
}

// Build an Environment from the given JSON string.
func EnvironmentFromString(jsonData string) (*Environment, os.Error) {
	env := new(Environment)
	err := CopyFromString(jsonData, env)
	if err != nil {
		return nil, err
	}
	return env, nil
}
