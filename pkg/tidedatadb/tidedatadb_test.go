package tidedatadb_test

import (
	"testing"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

func BenchmarkPerth3Solver(b *testing.B) {
	tideDataDb, err := tidedatadb.OpenTideDataDb("../../.data/dtu16.nc", tidedatadb.MODE_READONLY)
	if err != nil {
		b.Fatal(err)
	}
	defer tideDataDb.Close()
	for n := 0; n < b.N; n++ {
		constituent, err := tideDataDb.GetConstituentData(constituents.C_S1)
		if err != nil {
			b.Fatal(err)
		}
		_, err = constituent.GetDataXY(100, 100)
		if err != nil {
			b.Fatal(err)
		}
	}
}
