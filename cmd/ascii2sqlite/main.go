package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/mzeiher/perth3-go/pkg/grid"
	"github.com/mzeiher/perth3-go/pkg/loader"
)

func main() {
	reader, err := loader.CreateNewAsciiLoader("./.data/tide/tidal_constituents.dat")
	// reader, err := loader.CreateNewAsciiLoader("./.data/tide/fort.30")
	if err != nil {
		panic(err)
	}

	// frequencies := make([][][][]float64, 0)

	for {
		tideDataAmp, err := reader.GetNextTideGrid()
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}
		tideDataPhase, err := reader.GetNextTideGrid()
		if err != nil {
			panic(err)
		}

		cossin, err := grid.CalculateHCosHSinFromAmplitudeAndPhase(tideDataAmp, tideDataPhase)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d\n", len(cossin))
	}
}
