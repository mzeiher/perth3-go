package loader

import (
	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/mss"
)

type TideDataLoader interface {

	// returns the next set of grid data, where Data[0][0] is the south-west corner and Data[M][N] is the northeast corner
	// if the end is reached and no more data is available a `io.EOF` error is thrown
	GetNextTideGrid() (*constituents.TideConstituentData, error)
}

type MeanSeaSurfaceLoader interface {
	GetMSSData() (*mss.MedianSeaSurfaceData, error)
}
