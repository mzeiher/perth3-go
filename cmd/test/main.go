package main

import (
	"fmt"
	"time"

	"github.com/mzeiher/perth3-go/pkg/tide"
)

func main() {

	year, month, day, hour, min, sec := 2000, time.January, 1, 0, 0, 0
	t := time.Date(year, month, day, hour, min, sec, 0, time.UTC)

	// brest
	// var lat float64 = 48.389999
	// var lon float64 = -4.49

	// cap ferret
	var lat float64 = 44.617
	var lon float64 = -1.25

	// Get TideHeight
	tideHeight, err := tide.GetTideHeightsAtLocationAtTime("./.data/tide/constituents.dat", t, lat, lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tideHeight: %f\n", tideHeight)

	// Get TideHeight for Timespan
	tideHeights, err := tide.GetTideHeightsAtLocationForTimeSpan("./.data/tide/constituents.dat", t, time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC), 15*time.Minute, lat, lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tideHeights: %f\n", tideHeights)

	tideLat, tideHat, err := tide.GetLatAndHatAtLocation("./.data/tide/constituents.dat", lat, lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LAT: %f, HAT: %f\n", tideLat, tideHat)
}
