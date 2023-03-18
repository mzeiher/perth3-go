package datetime

import "time"

func UTCTimeToMJD(utcTime time.Time) float64 {
	return ((float64(utcTime.Unix()) / 86400.0) + 40587.0)
}
