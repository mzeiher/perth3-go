package utils

type Interpolatable interface {
	GetDataXY(x uint32, y uint32) ([]float32, error)
}

func InterpolateValues(lat float32, lon float32, minLat float32, maxLat float32, minLon float32, maxLon float32, gridSizeX uint32, gridSizeY uint32, dataGrid Interpolatable, wrap bool) ([]float32, error) {

	xResInDegree := (maxLon - minLon) / float32(gridSizeX-1)
	yResInDegree := (maxLat - minLat) / float32(gridSizeY-1)

	x0PosForLon := MapValue(lon, minLon, maxLon, 0, float32(gridSizeX)-1)
	y0PosForLat := MapValue(lat, minLat, maxLat, 0, float32(gridSizeY)-1)

	x0PosForLonInt := uint32(x0PosForLon)
	y0PosForLatInt := uint32(y0PosForLat)

	x1PosForLonInt := x0PosForLonInt + 1
	y1PosForLatInt := y0PosForLatInt + 1

	// calculate weights
	weightWest := (xResInDegree - (lon - float32((x0PosForLonInt))*xResInDegree - minLon)) / xResInDegree
	weightEast := 1 - weightWest
	weightSouth := (yResInDegree - (lat - float32((y0PosForLatInt))*yResInDegree - minLat)) / yResInDegree
	weightNorth := 1 - weightSouth

	weightNW := weightNorth * weightWest
	weightNE := weightNorth * weightEast
	weightSE := weightSouth * weightEast
	weightSW := weightSouth * weightWest

	combinedWeight := weightNE + weightNW + weightSE + weightSW

	// special case, for example if 90,0 is queried, we need to wrap at the northpole, if 0,180 we need to wrap along the null meridian
	y1x0WithPossibleWrap := x0PosForLonInt
	y1x1WithPossibleWrap := x1PosForLonInt
	if wrap {
		if x1PosForLonInt >= gridSizeX {
			x1PosForLonInt = 0
		}
		if y1PosForLatInt >= gridSizeY {
			y1PosForLatInt = y0PosForLatInt
			y1x0WithPossibleWrap = (gridSizeX - 1) - x0PosForLonInt
			y1x1WithPossibleWrap = (gridSizeX - 1) - x1PosForLonInt
		}
	}

	constituentsNW, err := dataGrid.GetDataXY(y1x0WithPossibleWrap, y1PosForLatInt)
	if err != nil {
		return nil, err
	}
	constituentsNE, err := dataGrid.GetDataXY(y1x1WithPossibleWrap, y1PosForLatInt)
	if err != nil {
		return nil, err
	}
	constituentsSE, err := dataGrid.GetDataXY(x1PosForLonInt, y0PosForLatInt)
	if err != nil {
		return nil, err
	}
	constituentsSW, err := dataGrid.GetDataXY(x0PosForLonInt, y0PosForLatInt)
	if err != nil {
		return nil, err
	}

	values := make([]float32, len(constituentsNE))

	for i := 0; i < len(constituentsNE); i++ {

		values[i] = ((weightNW * constituentsNW[i]) +
			(weightNE * constituentsNE[i]) +
			(weightSE * constituentsSE[i]) +
			(weightSW * constituentsSW[i])) / combinedWeight

	}

	return values, nil
}
