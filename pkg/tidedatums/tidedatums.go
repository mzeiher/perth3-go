package tidedatums

import (
	"math"
	"time"

	"github.com/mzeiher/perth3-go/pkg/solver"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

type TideDatums struct {
	HAT  float32
	MHWS float32
	MHW  float32
	MHWN float32
	MSL  float32
	MLWN float32
	MLW  float32
	MLWS float32
	LAT  float32
}

func GetDatumsForLatLan(constituentDb *tidedatadb.TideDataDB, solverName string, lat float32, lon float32) (*TideDatums, error) {

	currentTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	step := 15 * time.Minute

	tideSolver, err := solver.GetSolver(solverName)
	if err != nil {
		return nil, err
	}

	lat_datum := 0.0
	hat_datum := 0.0

	tideHeight := 0.0
	nbrTideData := 0

	for {
		height, err := tideSolver(constituentDb, lat, lon, currentTime)
		if err != nil {
			return nil, err
		}

		lat_datum = math.Min(height, lat_datum)
		hat_datum = math.Max(height, hat_datum)
		tideHeight = height + tideHeight

		currentTime = currentTime.Add(step)
		if end.Sub(currentTime) < 0 {
			break
		}
		nbrTideData = nbrTideData + 1
	}

	return &TideDatums{
		HAT:  float32(hat_datum),
		LAT:  float32(lat_datum),
		MHWS: 0,
		MHW:  0,
		MHWN: 0,
		MSL:  (float32(tideHeight) / float32(nbrTideData)),
		MLWN: 0,
		MLW:  0,
		MLWS: 0,
	}, nil
}
