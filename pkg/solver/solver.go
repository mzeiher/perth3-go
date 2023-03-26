package solver

import (
	"errors"
	"time"

	"github.com/mzeiher/perth3-go/pkg/solver/perth3"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

var ErrNoSolverFound = errors.New("no solver found for input")

var availableSolver map[string]Solver = make(map[string]Solver)

func init() {
	availableSolver["perth3"] = perth3.Solve
}

type Solver func(constituentDb *tidedatadb.TideDataDB, lat float32, lon float32, timeUtc time.Time) (float64, error)

func GetSolver(solver string) (Solver, error) {
	if availableSolver[solver] == nil {
		return nil, ErrNoSolverFound
	}
	return availableSolver[solver], nil
}
