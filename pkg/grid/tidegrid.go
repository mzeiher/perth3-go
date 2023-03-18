package grid

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
	SizeX        int
	SizeY        int
	LatitudeMin  float64
	LatitudeMax  float64
	LongitudeMin float64
	LongitudeMax float64
	Data         [][]float64
	UndefValue   float64
}
