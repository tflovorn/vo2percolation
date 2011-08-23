package percolation

import (
	"fmt"
	"testing"
)

func TestEnvironmentFromJSON(t *testing.T) {
	Delta, V, Beta := 1.0, 0.5, 1.0
	data := fmt.Sprintf("{\"Delta\":%f, \"V\":%f, \"Beta\":%f}", Delta, V, Beta)
	env, err := EnvironmentFromString(data)
	if err != nil {
		t.Fatal(err)
	}
	if env.Delta != Delta || env.V != V {
		t.Fatalf("incorrect value in Environment")
	}
}
