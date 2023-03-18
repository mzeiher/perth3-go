package main

import (
	"fmt"
	"time"

	"github.com/mzeiher/perth3-go/pkg/tide"
)

func main() {
	// Define the Modified Julian Date (MJD) for a specific date and time
	year, month, day, hour, min, sec := 1985, time.January, 1, 0, 0, 10
	t := time.Date(year, month, day, hour, min, sec, 0, time.UTC)

	args := tide.DetermineEquilibriumTidalArguments(t)
	f, u := tide.DetermineNodalCorrections(t)
	fmt.Printf("%v\n", args)
	fmt.Printf("%v\n%v\n", f, u)
}
