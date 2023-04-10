package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mzeiher/perth3-go/pkg/loader"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

const supportedFormats = "Supported Formats:\n" +
	"dtu16ascii - ascii representation of DTU16 files (.fort30)\n" +
	"             all tide constituents should be concatenated before loading\n" +
	"             cat q1.d o1.d p1.d s1.d k1.d n2.d m2.d s2.d k2.d m4.d > fort.30\n"

// this command line utility creates a lookup database for the sin and cos components of
// the provided constituents.
//
// currently only the ASCII format used by the DTU-10/16 model is supported.
func main() {

	var format string
	flag.StringVar(&format, "format", "dtu16ascii", "format of input file")

	var help bool
	flag.BoolVar(&help, "help", false, "print help")

	flag.Parse()

	if help {
		printHelpAndExit(nil)
	}

	if format == "" {
		printHelpAndExit(errors.New("format option missing"))
	}
	if format != "dtu16ascii" {
		printHelpAndExit(errors.New("invalid format"))
		return
	}

	inFile := flag.Arg(0)
	outFile := flag.Arg(1)

	if inFile == "" || outFile == "" {
		printHelpAndExit(errors.New("must provide an INPUT and OUTPUT file"))
	}

	constituentReader, err := loader.GetLoader(format, inFile)
	if err != nil {
		printHelpAndExit(err)
	}
	defer constituentReader.Close()
	tideDbWriter, err := tidedatadb.OpenTideDataDb(outFile)
	if err != nil {
		printHelpAndExit(err)
	}
	defer tideDbWriter.Close()

	for {
		tideDataAmp, err := constituentReader.GetNextConstituentData()
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			printHelpAndExit(err)
		}
		fmt.Printf("loaded constituent %s %s\n", tideDataAmp.Constituent, tideDataAmp.Type)

		tideDataPhase, err := constituentReader.GetNextConstituentData()
		if err != nil {
			printHelpAndExit(err)
		}
		fmt.Printf("loaded constituent %s %s\n", tideDataPhase.Constituent, tideDataPhase.Type)
		if tideDataAmp.Constituent != tideDataPhase.Constituent {
			printHelpAndExit(fmt.Errorf("wrong constituents for amplitude and phase"))
		}
		if tideDataPhase.SizeX != tideDataAmp.SizeX ||
			tideDataPhase.SizeY != tideDataAmp.SizeY ||
			tideDataPhase.LatitudeMin != tideDataAmp.LatitudeMin ||
			tideDataPhase.LatitudeMax != tideDataAmp.LatitudeMax ||
			tideDataPhase.LongitudeMin != tideDataAmp.LongitudeMin ||
			tideDataPhase.LongitudeMax != tideDataAmp.LongitudeMax {
			printHelpAndExit(fmt.Errorf("dimensions of amplitude and phase data are not compatible"))
		}

		constituentEntry, err := tideDbWriter.CreateNewConstituentData(tidedatadb.Dimensions{
			MinLat:        tideDataAmp.LatitudeMin,
			MaxLat:        tideDataAmp.LatitudeMax,
			MinLon:        tideDataAmp.LongitudeMin,
			MaxLon:        tideDataAmp.LongitudeMax,
			ResolutionLat: (tideDataAmp.LatitudeMax - tideDataAmp.LatitudeMin) / float32(tideDataAmp.SizeY-1),
			ResolutionLon: (tideDataAmp.LongitudeMax - tideDataAmp.LongitudeMin) / float32(tideDataAmp.SizeX-1),
			GridXSize:     uint64(tideDataAmp.SizeX),
			GridYSize:     uint64(tideDataAmp.SizeY),
		}, tidedatadb.ConstituentInfo{
			Constituent:   tideDataAmp.Constituent,
			AmplitudeUnit: tidedatadb.UNIT_CM,
			PhaseUnit:     tidedatadb.UNIT_DEGREE,
		})
		if err != nil {
			printHelpAndExit(err)
		}

		fmt.Printf("Write constituent %s\n", constituentEntry.ConstituentInfo.Constituent)

		for y := 0; y < tideDataAmp.SizeY; y++ {
			for x := 0; x < tideDataAmp.SizeX; x++ {
				err = constituentEntry.WriteDataXY([]float32{tideDataAmp.Data[y][x], tideDataPhase.Data[y][x]}, uint64(x), uint64(y))
				if err != nil {
					printHelpAndExit(err)
				}
			}
		}
	}

}

func printHelpAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Usage of %s %s:\n", os.Args[0], " [OPTIONS] INPUT OUTPUT")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s", supportedFormats)
	if err != nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
