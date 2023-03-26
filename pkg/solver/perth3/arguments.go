package perth3

import (
	"time"

	"github.com/mzeiher/perth3-go/pkg/astro"
	"github.com/mzeiher/perth3-go/pkg/datetime"
)

func CalculateArguments(timeUtc time.Time) []float64 {

	// hour := utcTime.Hour()
	mjd := datetime.UTCTimeToMJD(timeUtc)
	julHour := float64(mjd-float64(int(mjd))) * 24
	t1 := float64(15 * julHour)
	t2 := float64(30 * julHour)

	shpn := astro.ComputeAstronomicalMeanLongitudesInDegree(timeUtc)

	arg := make([]float64, 28)

	arg[0] = t1 + shpn.L_h - 3*shpn.L_s + shpn.L_p - 90    // Q1
	arg[1] = t1 + shpn.L_h - 2*shpn.L_s - 90               // O1
	arg[2] = t1 - shpn.L_h - 90                            // P1
	arg[3] = t1 + shpn.L_h + 90                            // K1
	arg[4] = t2 + 2*shpn.L_h - 3*shpn.L_s + shpn.L_p       // N2
	arg[5] = t2 + 2*shpn.L_h - 2*shpn.L_s                  // M2
	arg[6] = t2                                            // S2
	arg[7] = t2 + 2*shpn.L_h                               // K2
	arg[8] = t1 - 4*shpn.L_s + shpn.L_h + 2*shpn.L_p - 90  // 2Q1
	arg[9] = t1 - 4*shpn.L_s + 3*shpn.L_h - 90             // sigma1
	arg[10] = t1 - 3*shpn.L_s + 3*shpn.L_h - shpn.L_p - 90 // rho1
	arg[11] = t1 - shpn.L_s + shpn.L_h - shpn.L_p + 90     // M1
	arg[12] = t1 - shpn.L_s + shpn.L_h + shpn.L_p + 90     // M1
	arg[13] = t1 - shpn.L_s + 3*shpn.L_h - shpn.L_p + 90   // chi1
	arg[14] = t1 - 2*shpn.L_h + shpn.L_P1 - 90             // pi1
	arg[15] = t1 + 3*shpn.L_h + 90                         // phi1
	arg[16] = t1 + shpn.L_s - shpn.L_h + shpn.L_p + 90     // theta1
	arg[17] = t1 + shpn.L_s + shpn.L_h - shpn.L_p + 90     // J1
	arg[18] = t1 + 2*shpn.L_s + shpn.L_h + 90              // OO1
	arg[19] = t2 - 4*shpn.L_s + 2*shpn.L_h + 2*shpn.L_p    // 2N2
	arg[20] = t2 - 4*shpn.L_s + 4*shpn.L_h                 // mu2
	arg[21] = t2 - 3*shpn.L_s + 4*shpn.L_h - shpn.L_p      // nu2
	arg[22] = t2 - shpn.L_s + shpn.L_p + 180               // lambda2
	arg[23] = t2 - shpn.L_s + 2*shpn.L_h - shpn.L_p + 180  // L2
	arg[24] = t2 - shpn.L_s + 2*shpn.L_h + shpn.L_p        // L2
	arg[25] = t2 - shpn.L_h + shpn.L_P1                    // T2
	arg[26] = t1 + 180                                     // S1 (Doodson's phase)
	arg[27] = 2 * arg[5]                                   // M4

	return arg

}
