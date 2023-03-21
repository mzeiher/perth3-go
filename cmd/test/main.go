package main

import (
	"fmt"
	"time"

	"github.com/mzeiher/perth3-go/pkg/loader"
	"github.com/mzeiher/perth3-go/pkg/tide"
)

func main() {

	// t := time.Date(2023, time.March, 20, 16, 0, 0, 0, time.UTC)
	localTime := time.Now()
	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	displayLocale, err := time.LoadLocation("Australia/NSW")
	// displayLocale, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}

	t := localTime.In(localLocation).UTC()

	fmt.Printf("time UTC           : %s\n", t.Format(time.RFC3339))
	fmt.Printf("time Display locale: %s\n", t.In(displayLocale).Format(time.RFC3339))

	// brest
	// var lat float64 = 48.389999
	// var lon float64 = -4.49

	// // cap ferret
	// var lat float64 = 44.617
	// var lon float64 = -1.25

	// // byron bay
	var lat float64 = -28.629794
	var lon float64 = 153.618673

	// // lagos, PT
	// var lat float64 = 37.098181
	// var lon float64 = -8.661505

	// // sagres, PT
	// var lat float64 = 37.010503
	// var lon float64 = -8.962977

	// st helier, channel islands
	// var lat float64 = 49.182736
	// var lon float64 = -2.127271

	// // Foula, Scotland
	// var lat float64 = 60.119397
	// var lon float64 = -2.045731

	// freiburg, germany :)
	// var lat float64 = 48.014915
	// var lon float64 = 7.850038

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
	tideHeights, err := tide.GetTideHeightsAtLocationForTimeSpan("./.data/tide/constituents.dat", t, t.Add(24*time.Hour), 5*time.Minute, lat, lon)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(tideHeights); i++ {
		timeUtc := t.Add((time.Duration(i) * 5) * time.Minute).In(displayLocale).Format(time.RFC3339)
		fmt.Printf("tideHeights %s: %10.4f relative lat: %10.4f +mss: %10.4f\n", timeUtc, tideHeights[i], tideHeights[i]-tideLat, tideHeights[i]-tideLat+mssHeight)
	}

}
