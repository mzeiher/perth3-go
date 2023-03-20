package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/datetime"
	"github.com/mzeiher/perth3-go/pkg/tide"
)

func main() {

	year, month, day, hour, min, sec := 2023, time.March, 20, 10, 0, 0
	// year, month, day, hour, min, sec := 1985, time.January, 1, 0, 0, 0
	t := time.Date(year, month, day, hour, min, sec, 0, time.UTC)

	// brest
	var lat float64 = 48.389999
	var lon float64 = -4.49

	file, err := os.OpenFile("./.data/tide/constituents.dat", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	computedHeights := constituents.InterpolateConstituentsAtPositionAndReturnHeight(file, lat, lon)
	heightsArgument := tide.MapComputedHeightsToHeightArray(computedHeights)
	inferredHeights := tide.InferMinorTideHeights(heightsArgument)
	for i := 0; i <= 500; i++ {
		arguments := tide.DetermineEquilibriumTidalArguments(t)
		nodalCorrectionsF, nodalCorrectionsU := tide.DetermineNodalCorrections(t)

		var sum float64 = 0
		// iterate over all heights
		for i := 0; i < 28; i++ {
			heightCos := inferredHeights[0][i]
			heightSin := inferredHeights[1][i]
			chiu := (arguments[i] + nodalCorrectionsU[i]) * (math.Pi / 180)
			sum = sum + heightCos*nodalCorrectionsF[i]*math.Cos(chiu) + heightSin*nodalCorrectionsF[i]*math.Sin(chiu)
		}

		lpeqomt := tide.CalculateLongPeriodEquilibriumOceanMeanTide(t, lat)
		fmt.Printf("[%s] [%10.8f] tide: %9.4f lpeqomt: %9.4f, sum: %9.4f\n", t.Format(time.RFC3339), datetime.UTCTimeToMJD(t)-46066, sum/100, lpeqomt/100, (sum+lpeqomt)/100)
		t = t.Add(time.Minute * 15)
	}

}
