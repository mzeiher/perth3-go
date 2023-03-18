package grid

import (
	"fmt"
	"math"
)

type TideConstituent string

const (
	Q1 TideConstituent = "Q1"
	O1 TideConstituent = "O1"
	P1 TideConstituent = "P1"
	S1 TideConstituent = "S1"
	K1 TideConstituent = "K1"
	N2 TideConstituent = "N2"
	M2 TideConstituent = "M2"
	S2 TideConstituent = "S2"
	K2 TideConstituent = "K2"
	M4 TideConstituent = "M4"
)

var TideConstituents = []TideConstituent{
	Q1, O1, P1, S1, K1, N2, M2, S2, K2, M4,
}

type TideValueType string

const (
	PHASE     TideValueType = "PHASE"
	AMPLITUDE TideValueType = "AMPLITUDE"
)

type TideGridData struct {
	Constituent  TideConstituent
	Type         TideValueType
	LatitudeMin  float64
	LatitudeMax  float64
	LongitudeMin float64
	LongitudeMax float64
	Data         [][]float64
	UndefValue   float64
}

// calculate the hcos(G) and hsin(G) from amplitude and phase
// hcos = amplitude * cos(phase * (pi/180))
// hsin = amplitude * sin(phase * (pi/180))
// outputs array[longitude][latitude][0 = hcos, 1 = hsin]
func CalculateHCosHSinFromAmplitudeAndPhase(amplitude *TideGridData, phase *TideGridData) ([][][]float64, error) {
	if amplitude.Constituent != phase.Constituent {
		return nil, fmt.Errorf("constituent does not match got %s expected %s", phase.Constituent, amplitude.Constituent)
	}
	if amplitude.Type != AMPLITUDE {
		return nil, fmt.Errorf("type does not match, got %s expected %s", amplitude.Type, AMPLITUDE)
	}
	if phase.Type != PHASE {
		return nil, fmt.Errorf("type does not match, got %s expected %s", amplitude.Type, PHASE)
	}
	if len(amplitude.Data) == 0 || len(amplitude.Data) != len(phase.Data) {
		return nil, fmt.Errorf("data-y length does not match amplitude: %d, phase %d", len(amplitude.Data), len(phase.Data))
	}

	var data [][][]float64 = make([][][]float64, len(amplitude.Data))

	for i := 0; i < len(amplitude.Data); i++ {
		if len(amplitude.Data[i]) != len(phase.Data[i]) {
			return nil, fmt.Errorf("data-x length does not match amplitude: %d, phase %d", len(amplitude.Data[i]), len(phase.Data[i]))
		}
		data[i] = make([][]float64, len(amplitude.Data[i]))
		for j := 0; j < len(amplitude.Data[i]); j++ {
			data[i][j] = make([]float64, 2)
			if amplitude.Data[i][j] == amplitude.UndefValue || phase.Data[i][j] == phase.UndefValue {
				data[i][j][0] = math.NaN()
				data[i][j][1] = math.NaN()
			} else {
				data[i][j][0] = amplitude.Data[i][j] * math.Cos(phase.Data[i][j]*(math.Pi/180))
				data[i][j][1] = amplitude.Data[i][j] * math.Sin(phase.Data[i][j]*(math.Pi/180))
			}
		}
	}

	return data, nil
}
