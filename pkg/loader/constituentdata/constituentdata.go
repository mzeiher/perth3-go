package constituentdata

import (
	"io"

	"github.com/mzeiher/perth3-go/pkg/constituents"
)

type ConstituentValueType string

const (
	PHASE     ConstituentValueType = "PHASE"
	AMPLITUDE ConstituentValueType = "AMPLITUDE"
)

type TideConstituentData struct {
	Constituent  constituents.Constituent
	Type         ConstituentValueType
	SizeX        int
	SizeY        int
	LatitudeMin  float32
	LatitudeMax  float32
	LongitudeMin float32
	LongitudeMax float32
	Data         [][]float32
	UndefValue   float32
}

type ConstituentDataLoader interface {
	io.Closer
	// returns the next set of grid data, where Data[0][0] is the south-west corner and Data[M][N] is the northeast corner
	// if the end is reached and no more data is available a `io.EOF` error is thrown
	GetNextConstituentData() (*TideConstituentData, error)
}

type CreateLoaderFunction func(filePath string) (ConstituentDataLoader, error)
