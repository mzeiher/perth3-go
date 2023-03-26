package astro

import (
	"math"
	"time"

	"github.com/mzeiher/perth3-go/pkg/datetime"
)

type MeanLongitudes struct {
	//mean longitude of moon
	L_s float64
	// mean longitude of sun
	L_h float64
	//mean longitude of lunar perigee
	L_p float64
	// mean longitude of ascending lunar node
	L_N float64
	// mean longitude of solar perigee
	L_P1 float64
}

// Computes the 5 basic astronomical mean longitudes  s, h, p, N, p'.
//
// Note N is not N', i.e. N is decreasing with time.
//
// TIME is UTC in decimal Modified Julian Day (MJD).
// All longitudes returned in degrees.
//
// R. D. Ray, NASA/GSFC   August 2003
//
// Most of the formulae for mean longitudes are extracted from
// Jean Meeus, Astronomical Algorithms, 2nd ed., 1998.
// Page numbers below refer to this book.
//
// Note: This routine uses TIME in UT and does not distinguish between
//
//	the subtle differences of UTC, UT1, etc.  This is more than adequate
//	for the calculation of these arguments, especially in tidal studies.
func ComputeAstronomicalMeanLongitudesInDegree(utcTime time.Time) MeanLongitudes {
	const CIRCLE float64 = 360

	et := datetime.GetEphemerisTimeCorrected(utcTime)

	meanLongitudes := MeanLongitudes{}

	// mean longitude of moon (p.338)
	meanLongitudes.L_s = (((-1.53388e-8*et+1.855835e-6)*et-1.5786e-3)*et+481267.88123421)*et + 218.3164477

	// mean elongation of moon (p.338)
	D := (((-8.8445-9*et+1.83195e-6)*et-1.8819e-3)*et+445267.1114034)*et + 297.8501921

	// mean longitude of sun
	meanLongitudes.L_h = meanLongitudes.L_s - D

	// mean longitude of lunar perigee (p.343)
	meanLongitudes.L_p = ((-1.249172e-5*et-1.032e-2)*et+4069.0137287)*et + 83.3532465

	// mean longitude of ascending lunar node (p.144)
	meanLongitudes.L_N = ((2.22222e-6*et+2.0708e-3)*et-1934.136261)*et + 125.04452

	// mean longitude of solar perigee (Simon et al., 1994)
	meanLongitudes.L_P1 = 282.94 + 1.7192*et

	meanLongitudes.L_s = math.Mod(meanLongitudes.L_s, CIRCLE)
	meanLongitudes.L_h = math.Mod(meanLongitudes.L_h, CIRCLE)
	meanLongitudes.L_p = math.Mod(meanLongitudes.L_p, CIRCLE)
	meanLongitudes.L_N = math.Mod(meanLongitudes.L_N, CIRCLE)
	meanLongitudes.L_P1 = math.Mod(meanLongitudes.L_P1, CIRCLE)

	if meanLongitudes.L_s < 0 {
		meanLongitudes.L_s = meanLongitudes.L_s + CIRCLE
	}
	if meanLongitudes.L_h < 0 {
		meanLongitudes.L_h = meanLongitudes.L_h + CIRCLE
	}
	if meanLongitudes.L_p < 0 {
		meanLongitudes.L_p = meanLongitudes.L_p + CIRCLE
	}
	if meanLongitudes.L_N < 0 {
		meanLongitudes.L_N = meanLongitudes.L_N + CIRCLE
	}
	if meanLongitudes.L_P1 < 0 {
		meanLongitudes.L_P1 = meanLongitudes.L_P1 + CIRCLE
	}

	return meanLongitudes
}
