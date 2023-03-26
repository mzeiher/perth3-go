package lpeqomt

import (
	"math"
	"time"

	"github.com/mzeiher/perth3-go/pkg/astro"
)

func CalculateLongPeriodEquilibriumOceanMeanTide(timeUtc time.Time, lat float32) float64 {

	const PSOL = 283 * (math.Pi / 180)
	shpnp := astro.ComputeAstronomicalMeanLongitudesInDegree(timeUtc)
	// mjdInSeconds := datetime.UTCTimeToMJD(timeUtc) * 86400
	// et := (mjdInSeconds - 4043174400) / 86400
	// shpnp := astro.SHPNP{
	// 	L_s: (math.Pi / 180) * math.Mod(290.21+et*13.1763965, 360),
	// 	L_h: (math.Pi / 180) * math.Mod(280.12+et*0.9856473, 360),
	// 	L_p: (math.Pi / 180) * math.Mod(274.35+et*0.1114041, 360),
	// 	L_N: (math.Pi / 180) * math.Mod(343.51+et*0.0529539, 360),
	// }

	zlp := 2.79*math.Cos(shpnp.L_N) - 0.49*math.Cos(shpnp.L_h-PSOL) - 3.1*math.Cos(2*shpnp.L_h)
	ph := shpnp.L_s

	zlp = zlp -
		0.67*math.Cos(ph-2*shpnp.L_h+shpnp.L_p) -
		(3.52-0.46*math.Cos(shpnp.L_N))*math.Cos(ph-shpnp.L_p)

	ph = ph + shpnp.L_s
	zlp = zlp - 6.66*math.Cos(ph) -
		2.76*math.Cos(ph+shpnp.L_N) -
		0.26*math.Cos(ph+2*shpnp.L_N) -
		0.58*math.Cos(ph-2*shpnp.L_h) -
		0.29*math.Cos(ph-2*shpnp.L_p)

	ph = ph + shpnp.L_s
	zlp = zlp - 1.27*math.Cos(ph-shpnp.L_p) -
		0.53*math.Cos(ph-shpnp.L_p+shpnp.L_N) -
		0.24*math.Cos(ph-2*shpnp.L_h+shpnp.L_p)

	s := math.Sin(0 * (math.Pi / 180))

	lpeqomt := 0.437 * zlp * (1.5*s*s - 0.5)

	return lpeqomt
}
