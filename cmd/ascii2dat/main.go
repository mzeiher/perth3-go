package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mzeiher/perth3-go/pkg/grid"
	"github.com/mzeiher/perth3-go/pkg/loader"
)

const HEADER_SIZE = 40
const ENTRY_SIZE = 84

const GRID_SIZE_Y uint32 = 2881
const GRID_SIZE_X uint32 = 5761

const MIN_LAT float64 = -90.0
const MAX_LAT float64 = 90.0

const MIN_LON float64 = -180.0
const MAX_LON float64 = 180.0

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
// offset
// +---------------------------------------+
// | byte 0-3    float32 - q1 cos          |
// | byte 4-7    float32 - q1 sin          |
// | byte 8-11   float32 - o1 cos          |
// | byte 12-15  float32 - o1 sin          |
// | byte 16-19  float32 - p1 cos          |
// | byte 20-23  float32 - p1 sin          |
// | byte 24-27  float32 - s1 cos          |
// | byte 28-31  float32 - s1 sin          |
// | byte 32-35  float32 - k1 cos          |
// | byte 36-39  float32 - k1 sin          |
// | byte 40-32  float32 - n2 cos          |
// | byte 44-47  float32 - n2 sin          |
// | byte 48-51  float32 - m2 cos          |
// | byte 52-55  float32 - m2 sin          |
// | byte 56-59  float32 - s2 cos          |
// | byte 60-63  float32 - s2 sin          |
// | byte 64-67  float32 - k2 cos          |
// | byte 68-71  float32 - k2 sin          |
// | byte 72-75  float32 - m4 cos          |
// | byte 76-79  float32 - m4 sin          |
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

	// write header
	writeValue(file, GRID_SIZE_X)
	writeValue(file, GRID_SIZE_Y)
	writeValue(file, MIN_LAT)
	writeValue(file, MAX_LAT)
	writeValue(file, MIN_LON)
	writeValue(file, MAX_LON)

	for {
		tideDataAmp, err := reader.GetNextTideGrid()
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}
		if tideDataAmp.SizeX != int(GRID_SIZE_X) || tideDataAmp.SizeY != int(GRID_SIZE_Y) || tideDataAmp.LatitudeMin != MIN_LAT || tideDataAmp.LatitudeMax != MAX_LAT || tideDataAmp.LongitudeMin != MIN_LON || tideDataAmp.LongitudeMax != MAX_LON {
			panic(fmt.Errorf("dimensions not right"))
		}
		tideDataPhase, err := reader.GetNextTideGrid()
		if err != nil {
			panic(err)
		}
		if tideDataAmp.Constituent != tideDataPhase.Constituent {
			panic(fmt.Errorf("wrong constituents for amplitude and phase"))
		}
		if tideDataPhase.SizeX != int(GRID_SIZE_X) || tideDataPhase.SizeY != int(GRID_SIZE_Y) || tideDataPhase.LatitudeMin != MIN_LAT || tideDataPhase.LatitudeMax != MAX_LAT || tideDataPhase.LongitudeMin != MIN_LON || tideDataPhase.LongitudeMax != MAX_LON {
			panic(fmt.Errorf("dimensions not right"))
		}

		cossin, err := grid.CalculateHCosHSinFromAmplitudeAndPhase(tideDataAmp, tideDataPhase)
		if err != nil {
			panic(err)
		}
		fmt.Printf("write data for constituent %s\n", tideDataAmp.Constituent)
		for i := 0; i < len(cossin); i++ {
			for j := 0; j < len(cossin[i]); j++ {
				writeCosSinForConstituent(file, tideDataAmp.Constituent, i, j, cossin[i][j])
			}
		}
	}
}

func writeValue(file *os.File, value any) {
	err := binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		panic(err)
	}
}

func writeCosSinForConstituent(file *os.File, constituent grid.TideConstituent, gridY int, gridX int, cossin []float64) {
	constituentIndex := 0
	for i := 0; i < len(grid.TideConstituents); i++ {
		if grid.TideConstituents[i] == constituent {
			constituentIndex = i
			break
		}
	}
	constituentOffset := constituentIndex * 8
	yPosOffset := ((gridY * int(GRID_SIZE_X)) + gridX) * ENTRY_SIZE
	writeStartOffset := int64(HEADER_SIZE + yPosOffset + constituentOffset)
	file.Seek(writeStartOffset, 0)
	err := binary.Write(file, binary.LittleEndian, float32(cossin[0]))
	if err != nil {
		panic(err)
	}
	err = binary.Write(file, binary.LittleEndian, float32(cossin[1]))
	if err != nil {
		panic(err)
	}
}
