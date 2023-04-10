package tidedatadb

import (
	"errors"
	"math"

	"github.com/fhs/go-netcdf/netcdf"
	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/utils"
)

var (
	ErrConstituentNotFound    = errors.New("constituent not found")
	ErrConstituentAlreadyInDb = errors.New("constituent already in DB")
	ErrUnitNotFound           = errors.New("unit not found")
)

type ConstituentAmplitudeUnit byte

const (
	UNIT_CM ConstituentAmplitudeUnit = iota
	UNIT_METER
	UNIT_FEET
)

func ConstituentAmplitudeUnitFromString(name string) (ConstituentAmplitudeUnit, error) {
	switch name {
	case "CM":
		return UNIT_CM, nil
	case "M":
		return UNIT_METER, nil
	case "FT":
		return UNIT_FEET, nil
	}
	return 0, ErrUnitNotFound
}

func (c ConstituentAmplitudeUnit) String() string {
	switch c {
	case UNIT_CM:
		return "CM"
	case UNIT_METER:
		return "M"
	case UNIT_FEET:
		return "FT"
	}
	return ""
}

type ConstituentPhaseUnit byte

const (
	UNIT_DEGREE ConstituentPhaseUnit = iota
	UNIT_RADIAN
)

func ConstituentPhaseUnitFromString(name string) (ConstituentPhaseUnit, error) {
	switch name {
	case "DEGREE":
		return UNIT_DEGREE, nil
	case "RADIAN":
		return UNIT_RADIAN, nil
	}
	return 0, ErrUnitNotFound
}

func (c ConstituentPhaseUnit) String() string {
	switch c {
	case UNIT_DEGREE:
		return "DEGREE"
	case UNIT_RADIAN:
		return "RADIAN"
	}
	return ""
}

const ATTR_UNIT_AMPLITUDE = "UNIT_AMP"
const ATTR_UNIT_PHASE = "UNIT_PHASE"

type ConstituentInfo struct {
	Constituent   constituents.Constituent
	AmplitudeUnit ConstituentAmplitudeUnit
	PhaseUnit     ConstituentPhaseUnit
}

func (t *TideDataDB) GetConstituentData(constituent constituents.Constituent) (*ConstituentData, error) {
	variable, err := t.file.Var(constituent.String())
	if err != nil {
		return nil, err
	}

	attrAmpUnit, err := utils.NetcdfGetStringFromAttribute(ATTR_UNIT_AMPLITUDE, &variable)
	if err != nil {
		return nil, err
	}
	attrPhaseUnit, err := utils.NetcdfGetStringFromAttribute(ATTR_UNIT_PHASE, &variable)
	if err != nil {
		return nil, err
	}

	ampUnit, err := ConstituentAmplitudeUnitFromString(attrAmpUnit)
	if err != nil {
		return nil, err
	}
	phaseUnit, err := ConstituentPhaseUnitFromString(attrPhaseUnit)
	if err != nil {
		return nil, err
	}

	dimensionsLat, err := t.file.Dim("lat")
	if err != nil {
		return nil, err
	}
	dimensionLatLen, err := dimensionsLat.Len()
	if err != nil {
		return nil, err
	}

	dimensionsLon, err := t.file.Dim("lon")
	if err != nil {
		return nil, err
	}
	dimensionLonLen, err := dimensionsLon.Len()
	if err != nil {
		return nil, err
	}

	latVar, err := t.file.Var("lat")
	if err != nil {
		return nil, err
	}

	lonVar, err := t.file.Var("lon")
	if err != nil {
		return nil, err
	}

	minLat, err := latVar.ReadFloat64At([]uint64{0})
	if err != nil {
		return nil, err
	}
	maxLat, err := latVar.ReadFloat64At([]uint64{dimensionLatLen - 1})
	if err != nil {
		return nil, err
	}
	minLon, err := lonVar.ReadFloat64At([]uint64{0})
	if err != nil {
		return nil, err
	}
	maxLon, err := lonVar.ReadFloat64At([]uint64{dimensionLonLen - 1})
	if err != nil {
		return nil, err
	}

	return &ConstituentData{
		variable: &variable,
		Dimensions: Dimensions{
			MinLat:        float32(minLat),
			MaxLat:        float32(maxLat),
			MinLon:        float32(minLon),
			MaxLon:        float32(maxLon),
			GridXSize:     dimensionLonLen,
			GridYSize:     dimensionLatLen,
			ResolutionLat: (float32(maxLat) - float32(minLat)) / float32(dimensionLatLen-1),
			ResolutionLon: (float32(maxLon) - float32(minLon)) / float32(dimensionLonLen-1),
		},
		ConstituentInfo: ConstituentInfo{
			Constituent:   constituent,
			AmplitudeUnit: ampUnit,
			PhaseUnit:     phaseUnit,
		},
	}, nil
}

func (t *TideDataDB) CreateNewConstituentData(dimensionsToCreate Dimensions, constituentInfoToCreate ConstituentInfo) (*ConstituentData, error) {
	// create dimensions if not exist
	var dimLat netcdf.Dim
	var dimLon netcdf.Dim
	var dimData netcdf.Dim
	var err error

	if dimLat, err = t.file.Dim("lat"); err != nil {
		dimLat, err = t.file.AddDim("lat", dimensionsToCreate.GridYSize)
		if err != nil {
			return nil, err
		}
	}
	if dimLon, err = t.file.Dim("lon"); err != nil {
		dimLon, err = t.file.AddDim("lon", dimensionsToCreate.GridXSize)
		if err != nil {
			return nil, err
		}
	}
	if dimData, err = t.file.Dim("data"); err != nil {
		dimData, err = t.file.AddDim("data", 2)
		if err != nil {
			return nil, err
		}
	}
	// create dimension data if not exist
	if _, err := t.file.Var("lat"); err != nil {
		dimLatVar, err := t.file.AddVar("lat", netcdf.DOUBLE, []netcdf.Dim{dimLat})
		if err != nil {
			return nil, err
		}
		for i := 0; i < int(dimensionsToCreate.GridYSize); i++ {
			dimLatVar.WriteFloat64At([]uint64{uint64(i)}, float64(dimensionsToCreate.MinLat)+(float64(i)*float64(dimensionsToCreate.ResolutionLat)))
		}
	}
	if _, err := t.file.Var("lon"); err != nil {
		dimLonVar, err := t.file.AddVar("lon", netcdf.DOUBLE, []netcdf.Dim{dimLon})
		if err != nil {
			return nil, err
		}
		for i := 0; i < int(dimensionsToCreate.GridXSize); i++ {
			dimLonVar.WriteFloat64At([]uint64{uint64(i)}, float64(dimensionsToCreate.MinLon)+(float64(i)*float64(dimensionsToCreate.ResolutionLon)))
		}
	}

	constituentVariable, err := t.file.AddVar(constituentInfoToCreate.Constituent.String(), netcdf.FLOAT, []netcdf.Dim{dimLat, dimLon, dimData})
	if err != nil {
		return nil, err
	}

	attrUnitAmplitude := constituentVariable.Attr(ATTR_UNIT_AMPLITUDE)
	err = attrUnitAmplitude.WriteBytes([]byte(constituentInfoToCreate.AmplitudeUnit.String()))
	if err != nil {
		return nil, err
	}

	attrUnitPhase := constituentVariable.Attr(ATTR_UNIT_PHASE)
	if err != nil {
		return nil, err
	}
	err = attrUnitPhase.WriteBytes([]byte(constituentInfoToCreate.PhaseUnit.String()))
	if err != nil {
		return nil, err
	}

	return &ConstituentData{
		variable:        &constituentVariable,
		Dimensions:      dimensionsToCreate,
		ConstituentInfo: constituentInfoToCreate,
	}, nil
}

type ConstituentData struct {
	variable        *netcdf.Var
	Dimensions      Dimensions
	ConstituentInfo ConstituentInfo
}

func (c *ConstituentData) WriteDataXY(amplitudePhase []float32, x uint64, y uint64) error {
	err := c.variable.WriteFloat32At([]uint64{y, x, 0}, amplitudePhase[0])
	if err != nil {
		return err
	}
	err = c.variable.WriteFloat32At([]uint64{y, x, 1}, amplitudePhase[1])
	if err != nil {
		return err
	}
	return nil
}

func (c *ConstituentData) GetDataXY(x uint64, y uint64) ([]float32, error) {
	amp, err := c.variable.ReadFloat32At([]uint64{y, x, 0})
	if err != nil {
		return nil, err
	}
	phase, err := c.variable.ReadFloat32At([]uint64{y, x, 1})
	if err != nil {
		return nil, err
	}
	return []float32{amp, phase}, nil
}

func (c *ConstituentData) GetDataInterpolatedLatLon(lat float32, lon float32) (*constituents.ConstituentDatum, error) {
	rawData, err := utils.InterpolateValues(lat, lon, c.Dimensions.MinLat, c.Dimensions.MaxLat, c.Dimensions.MinLon, c.Dimensions.MaxLon, c.Dimensions.GridXSize, c.Dimensions.GridYSize, c, true)
	if err != nil {
		return nil, err
	}
	if c.ConstituentInfo.AmplitudeUnit == UNIT_METER {
		rawData[0] = rawData[0] / 100
	} else if c.ConstituentInfo.AmplitudeUnit == UNIT_FEET {
		rawData[0] = rawData[0] * 30.48
	}
	if c.ConstituentInfo.PhaseUnit == UNIT_RADIAN {
		rawData[1] = rawData[1] * (180 / math.Pi)
	}
	return &constituents.ConstituentDatum{
		Constituent: c.ConstituentInfo.Constituent,
		Amplitude:   float64(rawData[0]),
		Phase:       float64(rawData[1]),
	}, nil
}
