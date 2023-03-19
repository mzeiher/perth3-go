package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/loader"
)

// write tide data in an mem-mappable format
// all numbers are in little-endian
//
// first 40byte HEADER
// +----------------------------+
// | byte 0-3    uint    yGrid  |
// | byte 4-7    uint    xGrid  |
// | byte 8-15   float64 minLat |
// | byte 16-23  float64 maxLat |
// | byte 24-31  float64 minLon |
// | byte 32-39  float64 maxLon |
// +----------------------------+
// after that in y-x order the data
// offset 40 -> first element (0,0)
// offset 124 -> second element (0,1)
// set data as an array [tide][cos,sin]float32
// +---------------------------------------+
// | byte 0-3    float32 - q1 hcos         |
// | byte 4-7    float32 - o1 hcos         |
// | byte 8-11   float32 - p1 hcos         |
// | byte 12-15  float32 - s1 hcos         |
// | byte 16-19  float32 - k1 hcos         |
// | byte 20-23  float32 - n2 hcos         |
// | byte 24-27  float32 - m2 hcos         |
// | byte 28-31  float32 - s2 hcos         |
// | byte 32-35  float32 - k2 hcos         |
// | byte 36-38  float32 - m4 hcos         |
// | byte 40-43  float32 - q1 hsin         |
// | byte 44-47  float32 - o1 hsin         |
// | byte 48-51  float32 - p1 hsin         |
// | byte 52-55  float32 - s1 hsin         |
// | byte 56-59  float32 - k1 hsin         |
// | byte 60-63  float32 - n2 hsin         |
// | byte 64-67  float32 - m2 hsin         |
// | byte 68-71  float32 - s2 hsin         |
// | byte 72-75  float32 - k2 hsin         |
// | byte 76-79  float32 - m4 hsin         |
// | byte 80-83  float32 - mean sea height |
// +---------------------------------------+
func main() {
	// reader, err := loader.CreateNewAsciiLoader("./.data/tide/tidal_constituents.dat")
	reader, err := loader.CreateNewAsciiLoader("./.data/tide/fort.30")
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("./constituents.dat", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(0, 0)
	// write header
	writeValue(file, constituents.GRID_SIZE_X)
	writeValue(file, constituents.GRID_SIZE_Y)
	writeValue(file, constituents.MIN_LAT)
	writeValue(file, constituents.MAX_LAT)
	writeValue(file, constituents.MIN_LON)
	writeValue(file, constituents.MAX_LON)

	for {
		tideDataAmp, err := reader.GetNextTideGrid()
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Printf("loaded constituent %s %s\n", tideDataAmp.Constituent, tideDataAmp.Type)
		if tideDataAmp.SizeX != int(constituents.GRID_SIZE_X) ||
			tideDataAmp.SizeY != int(constituents.GRID_SIZE_Y) ||
			tideDataAmp.LatitudeMin != constituents.MIN_LAT ||
			tideDataAmp.LatitudeMax != constituents.MAX_LAT ||
			tideDataAmp.LongitudeMin != constituents.MIN_LON ||
			tideDataAmp.LongitudeMax != constituents.MAX_LON {
			panic(fmt.Errorf("dimensions not right"))
		}
		tideDataPhase, err := reader.GetNextTideGrid()
		if err != nil {
			panic(err)
		}
		fmt.Printf("loaded constituent %s %s\n", tideDataPhase.Constituent, tideDataPhase.Type)
		if tideDataAmp.Constituent != tideDataPhase.Constituent {
			panic(fmt.Errorf("wrong constituents for amplitude and phase, can't calc [Hcos,Hsin]"))
		}
		if tideDataPhase.SizeX != int(constituents.GRID_SIZE_X) ||
			tideDataPhase.SizeY != int(constituents.GRID_SIZE_Y) ||
			tideDataPhase.LatitudeMin != constituents.MIN_LAT ||
			tideDataPhase.LatitudeMax != constituents.MAX_LAT ||
			tideDataPhase.LongitudeMin != constituents.MIN_LON ||
			tideDataPhase.LongitudeMax != constituents.MAX_LON {
			panic(fmt.Errorf("dimensions are not compatible"))
		}

		cossin, err := constituents.CalculateHCosHSinFromAmplitudeAndPhase(tideDataAmp, tideDataPhase)
		if err != nil {
			panic(err)
		}
		fmt.Printf("write data for constituent %s\n", tideDataAmp.Constituent)
		for i := 0; i < len(cossin); i++ {
			for j := 0; j < len(cossin[i]); j++ {
				constituents.WriteHCosHSinForConstituent(file, tideDataAmp.Constituent, i, j, cossin[i][j])
			}
		}
	}
	for i := 0; i < int(constituents.GRID_SIZE_Y); i++ {
		for j := 0; j < int(constituents.GRID_SIZE_X); j++ {
			constituents.WriteMeanSeaHeightForPoint(file, i, j, 9999.0)
		}
	}

}

func writeValue(file *os.File, value any) {
	err := binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		panic(err)
	}
}
