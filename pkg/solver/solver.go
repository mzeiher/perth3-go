package solver

import (
	"errors"
	"time"

	"github.com/mzeiher/perth3-go/pkg/solver/perth3"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

var ErrNoSolverFound = errors.New("no solver found for input")

var availableSolver map[Solver]CreateSolverFunc = make(map[Solver]CreateSolverFunc)

type Solver string

const (
	PERTH_3 Solver = "perth3"
	unknown Solver = "unknown"
)

func (s Solver) String() string {
	switch s {
	case PERTH_3:
		return "perth3"
	}
	return "unknown"
}

func GetSolverFromString(solver string) (Solver, error) {
	switch solver {
	case "perth3":
		return PERTH_3, nil
	}
	return unknown, ErrNoSolverFound
}

func init() {
	availableSolver[PERTH_3] = perth3.Solve
}

type CreateSolverFunc func(constituentDb *tidedatadb.TideDataDB, lat float32, lon float32, timeUtc time.Time) (float64, error)

func GetSolver(solver Solver) (CreateSolverFunc, error) {
	if availableSolver[solver] == nil {
		return nil, ErrNoSolverFound
	}
	return availableSolver[solver], nil
}
