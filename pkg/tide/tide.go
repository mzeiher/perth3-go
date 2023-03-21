package tide

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mzeiher/perth3-go/pkg/constituents"
)

func GetLatAndHatAtLocation(constituentsFile string, lat float64, lon float64) (float64, float64, error) {
	currentTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeEnd := time.Date(2020, time.December, 1, 0, 0, 0, 0, time.UTC)

	tide_lat := 0.0
	tide_hat := 0.0

	inferredHeights, err := GetCalculatedTideHeightsArrayAtLocation(constituentsFile, lat, lon)
	if err != nil {
		return 0, 0, err
	}

	for {
		tideHeight := GetTideHeightFromCalculatedHeightArrayForLocationAndTime(currentTime, inferredHeights, lat, lon)

		tide_hat = math.Max(tide_hat, tideHeight)
		tide_lat = math.Min(tide_lat, tideHeight)

		currentTime = currentTime.Add(60 * time.Minute)
		if timeEnd.Sub(currentTime) <= 0 {
			break
		}
	}

	return tide_lat, tide_hat, nil
}

func GetTideHeightsAtLocationForTimeSpan(constituentsFile string, startTimeUtc time.Time, endTimeUtc time.Time, stepSize time.Duration, lat float64, lon float64) ([]float64, error) {
	inferredHeights, err := GetCalculatedTideHeightsArrayAtLocation(constituentsFile, lat, lon)
	if err != nil {
		return nil, err
	}
	currentTime := startTimeUtc

	numberEntries := int((endTimeUtc.Sub(startTimeUtc) / stepSize))

	data := make([]float64, numberEntries)
	idx := 0
	for {
		if endTimeUtc.Sub(currentTime) <= 0 {
			break
		}
		tideHeight := GetTideHeightFromCalculatedHeightArrayForLocationAndTime(currentTime, inferredHeights, lat, lon)
		data[idx] = tideHeight

		idx++
		currentTime = currentTime.Add(stepSize)
	}

	return data, nil
}

func GetTideHeightAtLocationAtTime(constituentsFile string, timeUtc time.Time, lat float64, lon float64) (float64, error) {

	start := time.Now()
	inferredHeights, err := GetCalculatedTideHeightsArrayAtLocation(constituentsFile, lat, lon)
	if err != nil {
		return 0, err
	}
	fmt.Printf("lookup for tides took %s\n", time.Since(start))

	start = time.Now()
	height := GetTideHeightFromCalculatedHeightArrayForLocationAndTime(timeUtc, inferredHeights, lat, lon)
	fmt.Printf("calculation for tides took %s\n", time.Since(start))

	return height, nil
}

func GetCalculatedTideHeightsArrayAtLocation(constituentFile string, lat float64, lon float64) ([][]float64, error) {
	file, err := os.OpenFile(constituentFile, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	computedHeights := constituents.InterpolateConstituentsAtPositionAndReturnHeight(file, lat, lon)
	heightsArgument := MapComputedHeightsToHeightArray(computedHeights)
	inferredHeights := InferMinorTideHeights(heightsArgument)

	return inferredHeights, nil
}

func GetTideHeightFromCalculatedHeightArrayForLocationAndTime(timeUtc time.Time, inferredHeights [][]float64, lat float64, lon float64) float64 {

	arguments := DetermineEquilibriumTidalArguments(timeUtc)
	nodalCorrectionsF, nodalCorrectionsU := DetermineNodalCorrections(timeUtc)

	var sum float64 = 0
	// iterate over all heights
	for i := 0; i < 28; i++ {
		heightCos := inferredHeights[0][i]
		heightSin := inferredHeights[1][i]
		chiu := (arguments[i] + nodalCorrectionsU[i]) * (math.Pi / 180)
		sum = sum + heightCos*nodalCorrectionsF[i]*math.Cos(chiu) + heightSin*nodalCorrectionsF[i]*math.Sin(chiu)
	}

	lpeqomt := CalculateLongPeriodEquilibriumOceanMeanTide(timeUtc, lat)

	return sum + lpeqomt

}
