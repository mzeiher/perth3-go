package constituents

import "math"

type ConstituentDatum struct {
	Constituent Constituent
	Amplitude   float64
	Phase       float64
}

func (c *ConstituentDatum) GetHCos() float64 {
	return float64(c.Amplitude) * math.Cos(float64(c.Phase)*(math.Pi/180))
}
func (c *ConstituentDatum) GetHSin() float64 {
	return float64(c.Amplitude) * math.Sin(float64(c.Phase)*(math.Pi/180))
}
