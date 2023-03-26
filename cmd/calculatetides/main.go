package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mzeiher/perth3-go/v2/pkg/astro"
	"github.com/mzeiher/perth3-go/v2/pkg/constituents"
	"github.com/mzeiher/perth3-go/v2/pkg/tidedatadb"
)

func main() {

	var constituentDbPath string
	flag.StringVar(&constituentDbPath, "constituentdb", "", "Path to constituentdb")

	var tideDataCache string
	flag.StringVar(&tideDataCache, "tideDataCache", "", "Path to tide data cache (optional)")

	var startTimeString string
	flag.StringVar(&startTimeString, "tstart", time.Now().Format(time.RFC3339), "start time in rfc3339 format (default: now")

	var endTimeString string
	flag.StringVar(&endTimeString, "tend", time.Now().Format(time.RFC3339), "end time in rfc3339 format (optional)")

	var stepDurationSecondsInt int
	flag.IntVar(&stepDurationSecondsInt, "stepDurationSeconds", 60, "step duration in seconds (default: 60)")

	var help bool
	flag.BoolVar(&help, "help", false, "print help")

	flag.Parse()

	if help {
		printHelpAndExit(nil)
	}

	var lat, lon float32
	_, err := fmt.Sscanf(flag.Arg(0), "%f,%f", &lat, &lon)
	if err != nil {
		printHelpAndExit(err)
	}

	// load constituent db for lookup
	constituentDb, err := tidedatadb.OpenTideDataDb(constituentDbPath, tidedatadb.TYPE_CONSTITUENT, os.O_RDONLY)
	if err != nil {
		printHelpAndExit(err)
	}
	defer constituentDb.Close()

	// parse start/end time and duration
	startTime, err := time.Parse(time.RFC3339, startTimeString)
	if err != nil {
		printHelpAndExit(err)
	}
	endTime, err := time.Parse(time.RFC3339, endTimeString)
	if err != nil {
		printHelpAndExit(err)
	}
	stepDuration, err := time.ParseDuration(fmt.Sprintf("%ds", stepDurationSecondsInt))
	if err != nil {
		printHelpAndExit(err)
	}
	startTimeUTC := startTime.UTC()
	endTimeUTC := endTime.UTC()

	fmt.Printf("%s, %s, %s\n", startTimeUTC.String(), endTimeUTC.String(), stepDuration.String())

	fmt.Printf("%+v\n", astro.ComputeAstronomicalMeanLongitudesInDegree(startTimeUTC))

	constituentData, err := constituentDb.GetConstituentData(constituents.C_Q1)
	if err != nil {
		printHelpAndExit(err)
	}

	data, err := constituentData.GetDataInterpolatedLatLon(lat, lon)
	if err != nil {
		printHelpAndExit(err)
	}

	fmt.Printf("%v", data)

}

func printHelpAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Usage of %s %s:\n", os.Args[0], " [OPTIONS] \"lat,lon\"")
	flag.PrintDefaults()
	if err != nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
