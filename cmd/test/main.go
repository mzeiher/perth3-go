package main

import (
	"fmt"
	"time"

	"github.com/mzeiher/perth3-go/pkg/loader"
	"github.com/mzeiher/perth3-go/pkg/tide"
)

func main() {

	// t := time.Date(2023, time.March, 20, 16, 0, 0, 0, time.UTC)
	t := time.Now().UTC()

	// brest
	// var lat float64 = 48.389999
	// var lon float64 = -4.49

	// cap ferret
	var lat float64 = 44.617
	var lon float64 = -1.25

	mssLoader, err := loader.CreateNewAsciiMSSLoader("./.data/meanseasurface/DTU15MSS_2min.mss")
	if err != nil {
		panic(err)
	}
	data, err := mssLoader.GetMSSData()
	if err != nil {
		panic(err)
	}
	mssHeight := data.InterpolateDataForPoint(lat, lon)
	fmt.Printf("mssHeight: %f\n", mssHeight)

	tideLat, tideHat, err := tide.GetLatAndHatAtLocation("./.data/tide/constituents.dat", lat, lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LAT: %10.4f, HAT: %10.4f\n", tideLat, tideHat)

	// Get TideHeight
	tideHeight, err := tide.GetTideHeightAtLocationAtTime("./.data/tide/constituents.dat", t, lat, lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tideHeight: %10.4f\n", tideHeight)
	fmt.Printf("tideHeight relative LAT: %10.4f +mss %10.4f\n", (tideHeight - tideLat), (tideHeight-tideLat)+mssHeight)

	// Get TideHeight for Timespan
	tideHeights, err := tide.GetTideHeightsAtLocationForTimeSpan("./.data/tide/constituents.dat", t, t.Add(24*time.Hour), 15*time.Minute, lat, lon)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(tideHeights); i++ {
		timeUtc := t.Add((time.Duration(i) * 15) * time.Minute).Local().Format(time.RFC3339)
		fmt.Printf("tideHeights %s: %10.4f relative lat: %10.4f +mss: %10.4f\n", timeUtc, tideHeights[i], tideHeights[i]-tideLat, tideHeights[i]-tideLat+mssHeight)
	}

}
