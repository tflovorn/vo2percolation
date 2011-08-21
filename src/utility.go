// Utility functions which don't have a home elsewhere
package percolation

import (
	"os"
	"strings"
)

// Convert string to byte slice
func StringToBytes(str string) ([]byte, os.Error) {
	reader := strings.NewReader(str)
	bytes := make([]byte, len(str))
	for seen := 0; seen < len(str); {
		n, err := reader.Read(bytes)
		if err != nil {
			return nil, err
		}
		seen += n
	}
	return bytes, nil
}
