package tide

import (
	"math"
	"time"

	"github.com/mzeiher/perth3-go/pkg/astro"
)

// return f[0..27] and u[0..27] for time
func DetermineNodalCorrections(utcTime time.Time) ([]float64, []float64) {

	shpn := astro.Compute5BasicAstronomicalMeanLongitudesInDegree(utcTime)

	sinn := math.Sin(shpn.L_N * (math.Pi / 180))
	cosn := math.Cos(shpn.L_N * (math.Pi / 180))
	sin2n := math.Sin(2 * shpn.L_N * (math.Pi / 180))
	cos2n := math.Cos(2 * shpn.L_N * (math.Pi / 180))

	f := make([]float64, 28)
	u := make([]float64, 28)

	f[0] = 1.009 + 0.187*cosn - 0.015*cos2n
	f[1] = f[0]

	f[2] = 1.0

	f[3] = 1.006 + 0.115*cosn - 0.009*cos2n
	f[4] = 1.000 - 0.037*cosn
	f[5] = f[4]

	f[6] = 1.0

	f[7] = 1.024 + 0.286*cosn + 0.008*cos2n
	f[8] = math.Sqrt(math.Pow(1.0+0.189*cosn-0.0058*cos2n, 2) + math.Pow(0.189*sinn-0.0058*sin2n, 2))

	f[9] = f[8]
	f[10] = f[8]

	f[11] = math.Sqrt(math.Pow(1.0+0.185*cosn, 2) + math.Pow(0.185*sinn, 2))
	f[12] = math.Sqrt(math.Pow(1.0+0.201*cosn, 2) + math.Pow(0.201*sinn, 2))
	f[13] = math.Sqrt(math.Pow(1.0+0.221*cosn, 2) + math.Pow(0.221*sinn, 2))

	f[14] = 1.0
	f[15] = 1.0
	f[16] = 1.0

	f[17] = math.Sqrt(math.Pow(1.0+0.198*cosn, 2) + math.Pow(0.198*sinn, 2))
	f[18] = math.Sqrt(math.Pow(1.0+0.640*cosn+0.134*cos2n, 2) + math.Pow(0.640*sinn+0.134*sin2n, 2))
	f[19] = math.Sqrt(math.Pow(1.0-0.0373*cosn, 2) + math.Pow(0.0373*sinn, 2))

	f[20] = f[19]
	f[21] = f[19]

	f[22] = 1.0

	f[23] = f[19]

	f[24] = math.Sqrt(math.Pow(1.0+0.441*cosn, 2) + math.Pow(0.441*sinn, 2))

	f[25] = 1.0
	f[26] = 1.0

	f[27] = f[5] * f[5]

	u[0] = 10.8*sinn - 1.3*sin2n
	u[1] = u[0]

	u[2] = 0.0

	u[3] = -8.9*sinn + 0.7*sin2n
	u[4] = -2.1 * sinn
	u[5] = u[4]

	u[6] = 0.0

	u[7] = -17.7*sinn + 0.7*sin2n
	u[8] = math.Atan2(0.189*sinn-0.0058*sin2n, 1.0+0.189*cosn-0.0058*sin2n) / (math.Pi / 180)
	u[9] = u[8]
	u[10] = u[8]
	u[11] = math.Atan2(0.185*sinn, 1.0+0.185*cosn) / (math.Pi / 180)
	u[12] = math.Atan2(-0.201*sinn, 1.0+0.201*cosn) / (math.Pi / 180)
	u[13] = math.Atan2(-0.221*sinn, 1.0+0.221*cosn) / (math.Pi / 180)

	u[14] = 0.0
	u[15] = 0.0
	u[16] = 0.0

	u[17] = math.Atan2(-0.198*sinn, 1.0+0.198*cosn) / (math.Pi / 180)
	u[18] = math.Atan2(-0.640*sinn-0.134*sin2n, 1.0+0.640*cosn+0.134*cos2n) / (math.Pi / 180)
	u[19] = math.Atan2(-0.0373*sinn, 1.0-0.0373*cosn) / (math.Pi / 180)
	u[20] = u[19]
	u[21] = u[19]

	u[22] = 0.0

	u[23] = u[19]
	u[24] = math.Atan2(-0.441*sinn, 1.0+0.441*cosn) / (math.Pi / 180)

	u[25] = 0.0
	u[26] = 0.0

	u[27] = 2 * u[5]

	return f, u
}
