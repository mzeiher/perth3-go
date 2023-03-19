package constituents

import (
	"encoding/binary"
	"os"

	"github.com/mzeiher/perth3-go/pkg/utils"
)

const HEADER_SIZE_BYTES = 40
const ENTRY_SIZE_BYTES = 84

// get the data as an array [tide][cos,sin]float32
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
func WriteHCosHSinForConstituent(file *os.File, constituent TideConstituent, gridY int, gridX int, hcoshsin []float64) {
	constituentIndex := 0
	for i := 0; i < len(TideConstituents); i++ {
		if TideConstituents[i] == constituent {
			constituentIndex = i
			break
		}
	}
	constituentOffset := constituentIndex * 4
	yPosOffset := ((gridY * int(GRID_SIZE_X)) + gridX) * ENTRY_SIZE_BYTES

	writeStartOffsetCos := int64(HEADER_SIZE_BYTES + yPosOffset + constituentOffset)
	file.Seek(writeStartOffsetCos, 0)
	err := binary.Write(file, binary.LittleEndian, float32(hcoshsin[0]))
	if err != nil {
		panic(err)
	}
	file.Seek(9*4, 1)
	err = binary.Write(file, binary.LittleEndian, float32(hcoshsin[1]))
	if err != nil {
		panic(err)
	}
}

func WriteMeanSeaHeightForPoint(file *os.File, gridY int, gridX int, value float32) {
	yPosOffset := ((gridY * int(GRID_SIZE_X)) + gridX) * ENTRY_SIZE_BYTES

	writeStartOffsetCos := int64(HEADER_SIZE_BYTES + yPosOffset + 20*4)
	file.Seek(writeStartOffsetCos, 0)
	err := binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		panic(err)
	}
}

// get the data as an array [tide][cos,sin]float32
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
func GetConstituentsForPosition(file *os.File, gridY int, gridX int) [][]float32 {
	yPosOffset := ((gridY * int(GRID_SIZE_X)) + gridX) * ENTRY_SIZE_BYTES
	readOffsetStart := int64(HEADER_SIZE_BYTES + yPosOffset)
	file.Seek(readOffsetStart, 0)

	constituents := make([][]float32, 2)
	for cossin := 0; cossin < 2; cossin++ {
		constituents[cossin] = make([]float32, 10)
		for constituent := 0; constituent < 10; constituent++ {
			binary.Read(file, binary.LittleEndian, &constituents[cossin][constituent])
		}
	}

	return constituents
}

func InterpolateConstituentsAtPositionAndReturnHeight(file *os.File, lat float64, lon float64) [][]float64 {

	xResInDegree := (MAX_LON - MIN_LON) / float64(GRID_SIZE_X-1)
	yResInDegree := (MAX_LAT - MIN_LAT) / float64(GRID_SIZE_Y-1)

	x0PosForLon := mapLon(lon) + 1
	y0PosForLat := mapLat(lat) + 1

	x0PosForLonInt := int(x0PosForLon)
	y0PosForLatInt := int(y0PosForLat)

	x1PosForLonInt := x0PosForLonInt + 1
	if x1PosForLonInt >= int(GRID_SIZE_X) {
		x1PosForLonInt = 0
	}
	y1PosForLatInt := y0PosForLatInt + 1

	weightWest := (xResInDegree - (lon - float64((x0PosForLonInt-1))*xResInDegree - MIN_LON)) / xResInDegree
	weightEast := 1 - weightWest
	weightSouth := (yResInDegree - (lat - float64((y0PosForLatInt-1))*yResInDegree - MIN_LAT)) / yResInDegree
	weightNorth := 1 - weightSouth

	weightNW := weightNorth * weightWest
	weightNE := weightNorth * weightEast
	weightSE := weightSouth * weightEast
	weightSW := weightSouth * weightWest

	combinedWeight := weightNE + weightNE + weightSE + weightSW

	constituentsNW := GetConstituentsForPosition(file, y1PosForLatInt, x1PosForLonInt)
	constituentsNE := GetConstituentsForPosition(file, y1PosForLatInt, x1PosForLonInt)
	constituentsSE := GetConstituentsForPosition(file, y0PosForLatInt, x1PosForLonInt)
	constituentsSW := GetConstituentsForPosition(file, y0PosForLatInt, x0PosForLonInt)

	heights := make([][]float64, 2)

	for cossin := 0; cossin < 2; cossin++ {
		heights[cossin] = make([]float64, 10)
		for constituents := 0; constituents < 10; constituents++ {
			heights[cossin][constituents] = ((weightNW * float64(constituentsNW[cossin][constituents])) +
				(weightNE * float64(constituentsNE[cossin][constituents])) +
				(weightSE * float64(constituentsSE[cossin][constituents])) +
				(weightSW * float64(constituentsSW[cossin][constituents]))) / combinedWeight
		}
	}

	return heights
}

func mapLon(lon float64) float64 {
	return utils.MapValue(lon, MIN_LON, MAX_LON, 0, float64(GRID_SIZE_X)-1)
}
func mapLat(lat float64) float64 {
	return utils.MapValue(lat, MIN_LAT, MAX_LAT, 0, float64(GRID_SIZE_Y)-1)
}
