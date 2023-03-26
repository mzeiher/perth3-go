package tidedatadb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/utils"
)

var (
	ErrConstituentNotFound    = errors.New("constituent not found")
	ErrConstituentAlreadyInDb = errors.New("constituent already in DB")
)

type ConstituentAmplitudeUnit byte

const (
	UNIT_CM ConstituentAmplitudeUnit = iota
	UNIT_METER
	UNIT_FEET
)

type ConstituentPhaseUnit byte

const (
	UNIT_DEGREE ConstituentPhaseUnit = iota
	UNIT_RADIAN
)

type ConstituentInfo struct {
	Constituent   constituents.Constituent
	AmplitudeUnit ConstituentAmplitudeUnit
	PhaseUnit     ConstituentPhaseUnit
}

func (t *TideDataDB) GetConstituentData(constituent constituents.Constituent) (*ConstituentData, error) {
	if t.Type != TYPE_CONSTITUENT {
		return nil, ErrInvalidDataType
	}
	t.file.Seek(int64(PREAMBLE_SIZE), 0)
	offset := PREAMBLE_SIZE
	for {
		t.file.Seek(int64(offset), 0)
		var currentEntry DataEntry
		err := binary.Read(t.file, binary.BigEndian, &currentEntry)
		if errors.Is(err, io.EOF) {
			return nil, ErrConstituentNotFound
		} else if err != nil {
			return nil, err
		} else {
			var currentConstituentInfo ConstituentInfo
			err := binary.Read(t.file, binary.BigEndian, &currentConstituentInfo)
			if err != nil {
				return nil, err
			}
			if currentConstituentInfo.Constituent == constituent {
				return &ConstituentData{
					db:              t,
					Offset:          uint32(offset),
					ConstituentInfo: currentConstituentInfo,
					Header:          currentEntry,
				}, nil
			} else {
				// advance offset
				offset = offset + int(currentEntry.Length)
			}
		}
	}

}

func (t *TideDataDB) CreateNewConstituentData(entryToCreate DataEntry, constituentInfoToCreate ConstituentInfo) (*ConstituentData, error) {
	if t.Type != TYPE_CONSTITUENT {
		return nil, ErrInvalidDataType
	}
	offset := PREAMBLE_SIZE
	for {
		t.file.Seek(int64(offset), 0)
		var currentEntry DataEntry
		err := binary.Read(t.file, binary.BigEndian, &currentEntry)
		if errors.Is(err, io.EOF) {
			// eof reached, create new entry to add data
			t.file.Seek(int64(offset), 0)
			// header + constituent info + data (2x4 byte per entry)
			length := 40 + 6 + ((entryToCreate.GridXSize * 8) * entryToCreate.GridYSize)
			entryToCreate.Length = length

			binary.Write(t.file, binary.BigEndian, entryToCreate)
			binary.Write(t.file, binary.BigEndian, constituentInfoToCreate)

			return &ConstituentData{
				db:              t,
				Offset:          uint32(offset),
				ConstituentInfo: constituentInfoToCreate,
				Header:          entryToCreate,
			}, nil
		} else if err != nil {
			return nil, err
		} else {
			// check currentConstituent
			var currentConstituent ConstituentInfo
			err := binary.Read(t.file, binary.BigEndian, &currentConstituent)
			if err != nil {
				return nil, err
			} else if currentConstituent.Constituent == constituentInfoToCreate.Constituent {
				// constituent already in file
				return nil, ErrConstituentAlreadyInDb
			} else {
				// advance offset
				offset = offset + int(currentEntry.Length)
			}

		}
	}

}

type ConstituentData struct {
	db              *TideDataDB
	Offset          uint32
	ConstituentInfo ConstituentInfo
	Header          DataEntry
}

func (c *ConstituentData) WriteDataXY(amplitudePhase []float32, x uint32, y uint32) error {
	if y >= c.Header.GridYSize {
		return errors.New("invalid index")
	}
	if x >= c.Header.GridXSize {
		return errors.New("invalid index")
	}
	if len(amplitudePhase) != 2 {
		return errors.New("invalid amplitude phase length")
	}
	offset := c.Offset + 40 + 6 + (y*c.Header.GridXSize*8 + x*8)
	c.db.file.Seek(int64(offset), 0)
	binary.Write(c.db.file, binary.BigEndian, amplitudePhase)

	return nil
}

func (c *ConstituentData) GetDataXY(x uint32, y uint32) ([]float32, error) {
	if y >= c.Header.GridYSize {
		return nil, errors.New("invalid index")
	}
	if x >= c.Header.GridXSize {
		return nil, errors.New("invalid index")
	}

	offset := c.Offset + 40 + 6 + (y*c.Header.GridXSize*8 + x*8)
	c.db.file.Seek(int64(offset), 0)

	data := make([]float32, 2)
	binary.Read(c.db.file, binary.BigEndian, data)

	return data, nil
}

func (c *ConstituentData) GetDataInterpolatedLatLon(lat float32, lon float32) (*constituents.ConstituentDatum, error) {
	rawData, err := utils.InterpolateValues(lat, lon, c.Header.MinLat, c.Header.MaxLat, c.Header.MinLon, c.Header.MaxLon, c.Header.GridXSize, c.Header.GridYSize, c, true)
	if err != nil {
		return nil, err
	}
	if c.ConstituentInfo.AmplitudeUnit == UNIT_METER {
		rawData[0] = rawData[0] * 100
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
