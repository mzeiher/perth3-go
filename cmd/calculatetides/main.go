package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mzeiher/perth3-go/pkg/solver"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
	"github.com/mzeiher/perth3-go/pkg/tidedatums"
)

const supportedSolvers = "Supported solver:\n" +
	"perth3 - perth3 solver from the dtu\n"

func main() {

	var constituentDbPath string
	flag.StringVar(&constituentDbPath, "constituentdb", "", "Path to constituentdb")

	var tideDataCache string
	flag.StringVar(&tideDataCache, "tideDataCache", "", "Path to tide data cache (optional)")

	var startTimeString string
	flag.StringVar(&startTimeString, "tstart", time.Now().Format(time.RFC3339), "start time in rfc3339 format")

	var endTimeString string
	flag.StringVar(&endTimeString, "tend", "", "end time in rfc3339 format (optional)")

	var stepDurationString string
	flag.StringVar(&stepDurationString, "stepduration", "60s", "step duration in seconds")

	var solverString string
	flag.StringVar(&solverString, "solver", "perth3", "solver to use")

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
	constituentDb, err := tidedatadb.OpenTideDataDb(constituentDbPath, tidedatadb.MODE_READONLY)
	if err != nil {
		printHelpAndExit(err)
	}
	defer constituentDb.Close()

	// parse start/end time and duration
	startTime, err := time.Parse(time.RFC3339, startTimeString)
	if err != nil {
		printHelpAndExit(err)
	}
	if endTimeString == "" {
		endTimeString = startTimeString
	}
	endTime, err := time.Parse(time.RFC3339, endTimeString)
	if err != nil {
		printHelpAndExit(err)
	}
	stepDuration, err := time.ParseDuration(stepDurationString)
	if err != nil {
		printHelpAndExit(err)
	}
	startTimeUTC := startTime.UTC()
	endTimeUTC := endTime.UTC()

	if endTimeUTC.Sub(startTimeUTC) < 0 {
		printHelpAndExit(errors.New("start time must be before end time"))
	}

	solverFunc, err := solver.GetSolver(solverString)
	if err != nil {
		printHelpAndExit(err)
	}

	tideDatums, err := tidedatums.GetDatumsForLatLan(constituentDb, solverString, lat, lon)
	if err != nil {
		panic(err)
	}

	currentTime := startTimeUTC

	fmt.Printf("%-10s %10.4fcm\n", "LAT", tideDatums.LAT)
	fmt.Printf("%-10s %10.4fcm\n", "MSL", tideDatums.MSL)
	fmt.Printf("%-10s %10.4fcm\n", "HAT", tideDatums.HAT)
	fmt.Printf("\n")
	fmt.Printf("%-25s %-11s %-11s\n", "date", "height (cm)", "tide (cm)")
	for {
		tideHeight, err := solverFunc(constituentDb, lat, lon, currentTime)
		if err != nil {
			panic(err)
		}

		tideHeightLAT := tideHeight - (float64(tideDatums.LAT - tideDatums.MSL))

		fmt.Printf("%-25s %11.4f %11.4f\n", currentTime.Local().Format(time.RFC3339), tideHeightLAT, tideHeight)

		currentTime = currentTime.Add(stepDuration)
		if endTimeUTC.Sub(currentTime) < 0 {
			return
		}
	}

}

func printHelpAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Usage of %s %s:\n", os.Args[0], " [OPTIONS] \"lat,lon\"")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s", supportedSolvers)
	if err != nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
