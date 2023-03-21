package datetime

type MoonPhase string

const (
	NewMoon        MoonPhase = "NewMoon"        // 0-1	   New Moon
	WaxingCrescent MoonPhase = "WaxingCrescent" // 2-6	   Waxing Crescent
	FirstQuarter   MoonPhase = "FirstQuarter"   // 7-8	   First Quarter
	WaxingGibbous  MoonPhase = "WaxingGibbous"  // 9-13   Waxing Gibbous
	FullMoon       MoonPhase = "FullMoon"       // 14-15	Full Moon
	WaningGibbous  MoonPhase = "WaningGibbous"  // 16-20	Waning Gibbous
	ThirdQuarter   MoonPhase = "ThirdQuarter"   // 21-22	Third Quarter
	WaningCrescent MoonPhase = "WaningCrescent" // 23-27	Waning Crescent
)

func GetMoonPhaseFromMJD(mjd float64) float64 {
	return ((mjd + 0.5) - 51544.5) / 29.530588853
}
