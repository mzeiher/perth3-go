package perth3

import (
	"math"
	"time"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/lpeqomt"
	"github.com/mzeiher/perth3-go/pkg/tidedatadb"
)

func Solve(constituentDb *tidedatadb.TideDataDB, lat float32, lon float32, timeUtc time.Time) (float64, error) {

	// solver array
	//   [
	//	   [c_hcos, c_hsin, arg, f, u],
	//     ...
	//   ]
	var solver [][]float64 = make([][]float64, 28)

	args := CalculateArguments(timeUtc)
	f, u := CalculateNodalCorrections(timeUtc)
	for i := 0; i < 28; i++ {
		solver[i] = make([]float64, 5)
		solver[i][2] = args[i]
		solver[i][3] = f[i]
		solver[i][4] = u[i]
	}

	constituentsForSolver := []constituents.Constituent{constituents.C_Q1, constituents.C_O1, constituents.C_P1, constituents.C_K1, constituents.C_N2, constituents.C_M2, constituents.C_S2, constituents.C_K2, constituents.C_S1, constituents.C_M4}

	// populate solver with data from the constituent db
	for index, constituent := range constituentsForSolver {
		constituentData, err := constituentDb.GetConstituentData(constituent)
		if err != nil {
			return 0, err
		}
		datum, err := constituentData.GetDataInterpolatedLatLon(lat, lon)
		if err != nil {
			return 0, err
		}
		solver[index][0] = datum.GetHCos()
		solver[index][1] = datum.GetHSin()
	}

	// move S2 to index 26 (like in perth3.f)
	solver[26][0] = solver[8][0]
	solver[26][1] = solver[8][1]
	// move M4 to index 27 (like in perth3.f)
	solver[27][0] = solver[9][0]
	solver[27][1] = solver[9][1]

	// now we infer the minor tide at position
	solver[8][0] = 0.263*solver[0][0] - 0.0252*solver[1][0]    // 2Q1 HCos
	solver[8][1] = 0.263*solver[0][1] - 0.0252*solver[1][1]    // 2Q1 HSin
	solver[9][0] = 0.297*solver[0][0] - 0.0264*solver[1][0]    // sigma1 HCos
	solver[9][1] = 0.297*solver[0][1] - 0.0264*solver[1][1]    // sigma1 HSin
	solver[10][0] = 0.164*solver[0][0] + 0.0048*solver[1][0]   // rho1 HCos
	solver[10][1] = 0.164*solver[0][1] + 0.0048*solver[1][1]   // rho1 HSin
	solver[11][0] = 0.0140*solver[1][0] + 0.0101*solver[3][0]  // M1 HCos
	solver[11][1] = 0.0140*solver[1][1] + 0.0101*solver[3][1]  // M1 HSin
	solver[12][0] = 0.0389*solver[1][0] + 0.0282*solver[3][0]  // M1 HCos
	solver[12][1] = 0.0389*solver[1][1] + 0.0282*solver[3][1]  // M1 HSin
	solver[13][0] = 0.0064*solver[1][0] + 0.0060*solver[3][0]  // chi1 HCos
	solver[13][1] = 0.0064*solver[1][1] + 0.0060*solver[3][1]  // chi1 HSin
	solver[14][0] = 0.0030*solver[1][0] + 0.0171*solver[3][0]  // pi1 HCos
	solver[14][1] = 0.0030*solver[1][1] + 0.0171*solver[3][1]  // pi1 HSin
	solver[15][0] = -0.0015*solver[1][0] + 0.0152*solver[3][0] // phi1 HCos
	solver[15][1] = -0.0015*solver[1][1] + 0.0152*solver[3][1] // phi1 HSin
	solver[16][0] = -0.0065*solver[1][0] + 0.0155*solver[3][0] // theta1 HCos
	solver[16][1] = -0.0065*solver[1][1] + 0.0155*solver[3][1] // theta1 HSin
	solver[17][0] = -0.0389*solver[1][0] + 0.0836*solver[3][0] // J1 HCos
	solver[17][1] = -0.0389*solver[1][1] + 0.0836*solver[3][1] // J1 HSin
	solver[18][0] = -0.0431*solver[1][0] + 0.0613*solver[3][0] // OO1 HCos
	solver[18][1] = -0.0431*solver[1][1] + 0.0613*solver[3][1] // OO1 HSin
	solver[19][0] = 0.264*solver[4][0] - 0.0253*solver[5][0]   // 2N2 HCos
	solver[19][1] = 0.264*solver[4][1] - 0.0253*solver[5][1]   // 2N2 HSin
	solver[20][0] = 0.298*solver[4][0] - 0.0264*solver[5][0]   // mu2 HCos
	solver[20][1] = 0.298*solver[4][1] - 0.0264*solver[5][1]   // mu2 HSin
	solver[21][0] = 0.165*solver[4][0] + 0.00487*solver[5][0]  // nu2 HCos
	solver[21][1] = 0.165*solver[4][1] + 0.00487*solver[5][1]  // nu2 HSin
	solver[22][0] = 0.0040*solver[5][0] + 0.0074*solver[6][0]  // lambda2 HCos
	solver[22][1] = 0.0040*solver[5][1] + 0.0074*solver[6][1]  // lambda2 HSin
	solver[23][0] = 0.0131*solver[5][0] + 0.0326*solver[6][0]  // L2 HCos
	solver[23][1] = 0.0131*solver[5][1] + 0.0326*solver[6][1]  // L2 HSin
	solver[24][0] = 0.0033*solver[5][0] + 0.0082*solver[6][0]  // L2 HCos
	solver[24][1] = 0.0033*solver[5][1] + 0.0082*solver[6][1]  // L2 HSin
	solver[25][0] = 0.0585 * solver[6][0]                      // T2 HCos
	solver[25][1] = 0.0585 * solver[6][1]                      // T2 HSin

	var sum float64 = 0
	// iterate over all heights
	for i := 0; i < 28; i++ {
		heightCos := solver[i][0]
		heightSin := solver[i][1]
		chiu := (solver[i][2] + solver[i][4]) * (math.Pi / 180)
		sum = sum + heightCos*solver[i][3]*math.Cos(chiu) + heightSin*solver[i][3]*math.Sin(chiu)
	}

	lpeqomt := lpeqomt.CalculateLongPeriodEquilibriumOceanMeanTide(timeUtc, lat)

	return sum + lpeqomt, nil
}
