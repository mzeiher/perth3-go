package datetime

import (
	"math"
	"time"
)

func GetEphemerisTimeSimple(utcTime time.Time) float64 {
	mjd := UTCTimeToMJD(utcTime)
	// Convert MJD to Julian Day (JD)
	jd := mjd + 2400000.5

	// Compute the correction factor between Universal Time (UT) and Dynamical Time (TD)
	t0 := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	daysSinceJ2000 := utcTime.Sub(t0).Hours() / 24.0
	deltaT := 62.92 + 0.32217*daysSinceJ2000 + 0.005589*math.Pow(daysSinceJ2000, 2.0)

	// Convert JD to Terrestrial Time (TT)
	tt := jd + deltaT/86400.0

	// Convert TT to Ephemeris Time (ET)
	return (tt - 2451545.0) / 36525.0
}

func GetEphemerisTimeCorrected(utcTime time.Time) float64 {
	mjd := UTCTimeToMJD(utcTime)

	// Convert MJD to Julian Day (JD)
	jd := mjd + 2400000.5
	deltaT := getDeltaTCorrection(utcTime)

	// Convert JD to Terrestrial Time (TT)
	tt := jd + deltaT/86400.0

	// Convert TT to Ephemeris Time (ET)
	return (tt - 2451545.0) / 36525.0
}

func getDeltaTCorrection(utcTime time.Time) float64 {
	year := utcTime.Year()

	var secdif float64 = 0.0

	//if we are outside the table on the low end
	// use the stephenson and morrison expression 948 to 1660,
	// and the borkowski formula for earlier years
	if year < 1660 {
		if year >= 948 {
			b := 0.01 * float64(year-2000)
			secdif = b*(b*23.58+100.3) + 101.6
		} else {
			b := 0.01*float64(year-2000) + 3.75
			secdif = 35.0*b*b + 40.0
		}
	} else if year >= 1660 && year < 2023 {
		// we are in the table do linear interpolate
		currentYearVal := correctionDeltaTTable[year-1660]
		nextYearVal := correctionDeltaTTable[year-1660+1]

		currentYearMDJ := UTCTimeToMJD(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC))
		nextYearMDJ := UTCTimeToMJD(time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC))
		yearToInterpolate := UTCTimeToMJD(utcTime)
		secdif = float64(currentYearVal + (yearToInterpolate-currentYearMDJ)*(nextYearVal-currentYearVal)/(nextYearMDJ-currentYearMDJ))
	} else {
		// otherwise do linear extrapolate
		b := float64(year - 1660)
		lastEntry := correctionDeltaTTable[len(correctionDeltaTTable)-1]
		secondToLastEntry := correctionDeltaTTable[len(correctionDeltaTTable)-2]
		secdif = lastEntry + b*(lastEntry-secondToLastEntry)
	}

	//.the astronomical almanac table is corrected by adding the expression
	//      -0.000091 (ndot + 26)(year-1955)^2  seconds
	//to entries prior to 1955 (page K8), where ndot is the secular tidal
	//term in the mean motion of the moon. entries after 1955 are referred
	//to atomic time standards and are not affected by errors in lunar
	//or planetary theory.  a value of ndot = -25.8 arcsec per century squared
	//is the value used in jpl's de403 ephemeris, the earlier de200 ephemeris
	//used the value -23.8946. note for years below the table (less than 1620)
	//the time difference is not adjusted for small improvements in the
	//current estimate of ndot because the formulas were derived from
	//studies of ancient eclipses and other historical information, whose
	//interpretation depends only partly on ndot.
	//here we make the ndot correction.

	if year < 1955 {
		b := float64(year - 1955)
		secdif = secdif - 0.000091*(-25.8+26.0)*b*b
	}

	return secdif
}

// https://maia.usno.navy.mil/ser7/deltat.data
// 1660 - 2023
var correctionDeltaTTable = []float64{
	38, 37, 36, 37, 38, 36, 35, 34, 33, 32,
	31, 30, 29, 29, 28, 27, 26, 25, 25, 26,
	26, 25, 24, 24, 24, 24, 24, 23, 23, 22,
	22, 22, 21, 21, 21, 21, 20, 20, 20, 20,
	21, 21, 20, 20, 19, 19, 19, 20, 20, 20,
	20, 20, 21, 21, 21, 21, 21, 21, 21, 21,
	21.1, 21.0, 20.9, 20.7, 20.4, 20.0, 19.4, 18.7, 17.8, 17.0,
	16.6, 16.1, 15.7, 15.3, 14.7, 14.3, 14.1, 14.1, 13.7, 13.5,
	13.5, 13.4, 13.4, 13.3, 13.2, 13.2, 13.1, 13.0, 13.3, 13.5,
	13.7, 13.9, 14.0, 14.1, 14.1, 14.3, 14.4, 14.6, 14.7, 14.7,
	14.8, 14.9, 15.0, 15.2, 15.4, 15.6, 15.6, 15.9, 15.9, 15.7,
	15.7, 15.7, 15.9, 16.1, 15.9, 15.7, 15.3, 15.5, 15.6, 15.6,
	15.6, 15.5, 15.4, 15.2, 14.9, 14.6, 14.3, 14.1, 14.2, 13.7,
	13.3, 13.0, 13.2, 13.1, 13.3, 13.5, 13.2, 13.1, 13.0, 12.6,
	12.6, 12.0, 11.8, 11.4, 11.1, 11.1, 11.1, 11.1, 11.2, 11.5,
	11.2, 11.7, 11.9, 11.8, 11.8, 11.8, 11.6, 11.5, 11.4, 11.3,
	11.13, 10.94, 10.29, 9.94, 9.88, 9.72, 9.66, 9.51, 9.21, 8.60,
	7.95, 7.59, 7.36, 7.10, 6.89, 6.73, 6.39, 6.25, 6.25, 6.22,
	6.22, 6.30, 6.35, 6.32, 6.33, 6.37, 6.40, 6.46, 6.48, 6.53,
	6.55, 6.69, 6.84, 7.03, 7.15, 7.26, 7.23, 7.21, 6.99, 7.19,
	7.35, 7.41, 7.36, 6.95, 6.45, 5.92, 5.15, 4.11, 2.94, 1.97,
	1.04, 0.11, -0.82, -1.70, -2.48, -3.19, -3.84, -4.43, -4.79, -5.09,
	-5.36, -5.37, -5.34, -5.40, -5.58, -5.74, -5.69, -5.67, -5.73, -5.78,
	-5.86, -6.01, -6.28, -6.53, -6.50, -6.41, -6.11, -5.63, -4.68, -3.72,
	-2.70, -1.48, -0.08, 1.26, 2.59, 3.92, 5.20, 6.29, 7.68, 9.13,
	10.38, 11.64, 13.23, 14.69, 16.00, 17.19, 18.19, 19.13, 20.14, 20.86,
	21.41, 22.06, 22.51, 23.01, 23.46, 23.63, 23.95, 24.39, 24.34, 24.10,
	24.02, 23.98, 23.89, 23.93, 23.88, 23.91, 23.76, 23.91, 23.96, 24.04,
	24.35, 24.82, 25.30, 25.77, 26.27, 26.76, 27.27, 27.77, 28.25, 28.70,
	29.15, 29.57, 29.97, 30.36, 30.72, 31.07, 31.349, 31.677, 32.166, 32.671,
	33.150, 33.584, 33.992, 34.466, 35.030, 35.738, 36.546, 37.429, 38.291, 39.204,
	40.182, 41.170, 42.227, 43.373, 44.486, 45.477, 46.458, 47.521, 48.535, 49.589,
	50.540, 51.382, 52.168, 52.957, 53.789, 54.3427, 54.8713, 55.3222, 55.8197, 56.3000,
	56.8553, 57.5653, 58.3092, 59.1218, 59.9845, 60.7853, 61.6287, 62.2950, 62.9659, 63.4673,
	63.8285, 64.0908, 64.2998, 64.4734, 64.5736, 64.6876, 64.8452, 65.1464, 65.4573, 65.7768,
	66.0699, 66.3246, 66.6030, 66.9069, 67.2810, 67.6439, 68.1024, 68.5927, 68.9676, 69.2202,
	69.3612, 69.3593, 69.2945, 69.2038,
}
