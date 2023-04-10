package perth3_test

import (
	"testing"
	"time"

	"github.com/mzeiher/perth3-go/pkg/solver/perth3"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

func BenchmarkPerth3Solver(b *testing.B) {
	tideDataDb, err := tidedatadb.OpenTideDataDb("../../../.data/dtu16.nc")
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := perth3.Solve(tideDataDb, 37.010503, -8.962977, time.Date(2023, 1, 1, 00, 00, 00, 00, time.UTC))
		if err != nil {
			b.Fatal(err)
		}
	}
}
