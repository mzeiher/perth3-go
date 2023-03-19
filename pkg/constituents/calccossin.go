package constituents

import (
	"fmt"
	"math"
)

// calculate the hcos(G) and hsin(G) from amplitude and phase
// hcos = amplitude * cos(phase * (pi/180))
// hsin = amplitude * sin(phase * (pi/180))
// outputs array[longitude][latitude][0 = hcos, 1 = hsin]
func CalculateHCosHSinFromAmplitudeAndPhase(amplitude *TideConstituentData, phase *TideConstituentData) ([][][]float64, error) {
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
