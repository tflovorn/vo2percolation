package percolation

import (
	"os"
	"json"
)

type GridAnalysis struct {
	totalSites int
	fermi float64
	grid [][]bool
}

const separatorRepeat = 3

// Iterate over all possible grid configurations with given grid length.
// For each configuration with only one cluster, collect data on it and export
// that data.
func BruteForceSurvey(gridLength int, ener Energetics, outputFilePath string) os.Error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	g := NewGridWithDims(gridLength, gridLength)
	for {
		// get analysis
		gAnalysis, err := AnalyzeCluster(g, ener)
		if err != nil {
			return err
		}
		// export analysis
		if gAnalysis != nil {
			marshalled, err := json.Marshal(gAnalysis)
			if err != nil {
				return err
			}
			marshalled = appendSeparator(marshalled)
			if _, err := outputFile.Write(marshalled); err != nil {
				return err
			}
		}
		// iterate grid
		if done := g.NextGrid(); done {
			// if we have covered all grids, stop
			break
		}
	}
	if err := outputFile.Close(); err != nil {
		return err
	}
	return nil
}

func AnalyzeCluster(g *Grid, ener Energetics) (*GridAnalysis, os.Error) {
	// only interested in grids with a single cluster
	if len(g.AllClusters()) != 1 {
		return nil, nil
	}
	// -- do cluster shape analysis --
	totalSites := g.ActiveSiteCount()

	// -- do energetic analysis --
	// fermi energy, one electron per site
	fermi, err := ener.FermiEnergy(g, totalSites)
	if err != nil {
		return nil, err
	}
	// -- pack data --
	ga := GridAnalysis{totalSites, fermi, g.data}
	return &ga, nil
}

func appendSeparator(data []byte) []byte {
	for i := 0; i < separatorRepeat; i++ {
		data = append(data, byte('\n'))
	}
	return data
}
