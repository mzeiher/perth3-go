package mss

import "github.com/mzeiher/perth3-go/pkg/utils"

type MedianSeaSurfaceData struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64

	LonRes float64
	LatRes float64

	GridX int
	GridY int

	Data [][]float32
}

func (mss *MedianSeaSurfaceData) InterpolateDataForPoint(lat float64, lon float64) float64 {

	// shift lon
	if lon <= 0 {
		lon = lon + 180
	}

	x0PosForLon := utils.MapValue(lon, mss.MinLon, mss.MaxLon, 0, float64(mss.GridX))
	y0PosForLat := utils.MapValue(lat, mss.MinLat, mss.MaxLat, 0, float64(mss.GridY))

	x0PosForLonInt := int(x0PosForLon)
	y0PosForLatInt := int(y0PosForLat)

	x1PosForLonInt := x0PosForLonInt + 1
	if x1PosForLonInt >= int(mss.GridX) {
		x1PosForLonInt = 0
	}
	y1PosForLatInt := y0PosForLatInt + 1

	weightWest := (mss.LonRes - (lon - float64((x0PosForLonInt))*mss.LonRes - mss.MinLon)) / mss.LonRes
	weightEast := 1 - weightWest
	weightSouth := (mss.LatRes - (lat - float64((y0PosForLatInt))*mss.LatRes - mss.MinLat)) / mss.LatRes
	weightNorth := 1 - weightSouth

	weightNW := weightNorth * weightWest
	weightNE := weightNorth * weightEast
	weightSE := weightSouth * weightEast
	weightSW := weightSouth * weightWest

	combinedWeight := weightNE + weightNW + weightSE + weightSW

	mssNW := mss.Data[y1PosForLatInt][x0PosForLonInt]
	mssNE := mss.Data[y1PosForLatInt][x1PosForLonInt]
	mssSE := mss.Data[y0PosForLatInt][x1PosForLonInt]
	mssSW := mss.Data[y0PosForLatInt][x0PosForLonInt]

	return ((weightNW * float64(mssNW)) +
		(weightNE * float64(mssNE)) +
		(weightSE * float64(mssSE)) +
		(weightSW * float64(mssSW))) / combinedWeight
}
