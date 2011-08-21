package percolation

import (
	"fmt"
	"testing"
)

func TestEnvironmentFromJSON(t *testing.T) {
	Delta, V := 1.0, 0.5
	data := fmt.Sprintf("{\"Delta\":%f, \"V\":%f}", Delta, V)
	env, err := EnvironmentFromString(data)
	if err != nil {
		t.Fatal(err)
	}
	if env.Delta != Delta || env.V != V {
		t.Fatalf("incorrect value in Environment")
	}
}
