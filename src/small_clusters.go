package percolation

func AnalyzeClusters(gridLength int, outputFilePath string) {
	g := NewGridWithDims(gridLength, gridLength)
	for {
		// do grid analysis

		// export analysis

		// iterate grid
		err := g.NextGrid()
		if err != nil {
			break
		}
	}
}
