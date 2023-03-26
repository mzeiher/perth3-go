package dtu16ascii

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mzeiher/perth3-go/pkg/constituents"
	"github.com/mzeiher/perth3-go/pkg/loader/constituentdata"
)

type asciiTideFile struct {
	constituentdata.ConstituentDataLoader
	reader *bufio.Reader
	file   *os.File
}

type asciiTideHeader struct {
	longitudeMin float32
	longitudeMax float32
	latitudeMin  float32
	latitudeMax  float32

	constituentType constituentdata.ConstituentValueType
	constituent     constituents.Constituent

	undefValue     float32
	entriesPerLine int

	gridX int
	gridY int
}

func CreateDTU16Loader(filePath string) (constituentdata.ConstituentDataLoader, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	return &asciiTideFile{
		reader: reader,
		file:   file,
	}, nil
}

func (a *asciiTideFile) Close() error {
	return a.file.Close()
}

func (a *asciiTideFile) GetNextConstituentData() (*constituentdata.TideConstituentData, error) {

	header, err := a.ParseHeader()
	if err != nil {
		return nil, err
	}

	gridData := &constituentdata.TideConstituentData{
		Constituent:  header.constituent,
		Type:         header.constituentType,
		LatitudeMin:  header.latitudeMin,
		LatitudeMax:  header.latitudeMax,
		LongitudeMin: header.longitudeMin,
		LongitudeMax: header.longitudeMax,

		SizeX: header.gridX,
		SizeY: header.gridY,

		UndefValue: header.undefValue,
		Data:       make([][]float32, header.gridY),
	}

	numberEntries := header.gridX * header.gridY
	currentY := 0
	currentEntry := 0
	gridData.Data[currentY] = make([]float32, header.gridX)
	for {

		line, err := a.reader.ReadString('\n')
		entries := strings.Fields(line)
		if err != nil {
			return nil, err
		}

		if len(entries) == 0 && currentEntry < numberEntries {
			return nil, fmt.Errorf("too few entries in table")
		}
		for index, entry := range entries {
			if currentEntry > numberEntries {
				return nil, fmt.Errorf("too many entries in table")
			}
			if currentEntry >= header.gridX*(currentY+1) {
				currentY = currentY + 1
				gridData.Data[currentY] = make([]float32, header.gridX)
			}

			entryParsed, err := strconv.ParseFloat(entry, 32)
			if err != nil {
				return nil, err
			}
			xPos := ((currentY * header.gridX) - currentEntry) * -1
			gridData.Data[currentY][xPos] = float32(entryParsed)

			if currentEntry == numberEntries-1 && index == len(entries)-1 {
				// we reached the end and red all entries
				return gridData, nil
			} else if currentEntry == numberEntries-1 && index < len(entries)-1 {
				return nil, fmt.Errorf("to many entries in file")
			} else {
				currentEntry = currentEntry + 1
			}
		}
	}

}

func (a *asciiTideFile) ParseHeader() (asciiTideHeader, error) {
	asciiHeader := asciiTideHeader{}

	// try to read the title
	title, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}

	constituent, err := constituents.FromString(strings.Fields(title)[0])
	if err != nil {
		return asciiHeader, fmt.Errorf("unknown constituent in title %s", title)
	}
	asciiHeader.constituent = constituent

	if strings.Contains(strings.ToLower(title), "amplitude") {
		asciiHeader.constituentType = constituentdata.AMPLITUDE
	} else if strings.Contains(strings.ToLower(title), "phase") {
		asciiHeader.constituentType = constituentdata.PHASE
	}

	// try to get type in second line
	description, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}
	if strings.Contains(strings.ToLower(description), "amplitude") {
		asciiHeader.constituentType = constituentdata.AMPLITUDE
	} else if strings.Contains(strings.ToLower(description), "phase") {
		asciiHeader.constituentType = constituentdata.PHASE
	}
	if asciiHeader.constituentType == "" {
		return asciiHeader, fmt.Errorf("constituent type not found")
	}

	// read grid size (x,y)
	gridSize, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}

	_, err = fmt.Sscanf(gridSize, "%d %d", &asciiHeader.gridY, &asciiHeader.gridX)
	if err != nil {
		return asciiHeader, err
	}

	// read latitude min/max
	latMinMax, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}

	_, err = fmt.Sscanf(latMinMax, "%f %f", &asciiHeader.latitudeMin, &asciiHeader.latitudeMax)
	if err != nil {
		return asciiHeader, err
	}

	// read latitude min/max
	lonMinMax, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}

	_, err = fmt.Sscanf(lonMinMax, "%f %f", &asciiHeader.longitudeMin, &asciiHeader.longitudeMax)
	if err != nil {
		return asciiHeader, err
	}

	// read UNDEF value
	undefValue, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}
	_, err = fmt.Sscanf(undefValue, "%f", &asciiHeader.undefValue)
	if err != nil {
		return asciiHeader, err
	}

	// read and forget fortran format
	entriesPerLine, err := a.reader.ReadString('\n')
	if err != nil {
		return asciiHeader, err
	}
	_, err = fmt.Sscanf(entriesPerLine, "(%d", &asciiHeader.entriesPerLine)
	if err != nil {
		return asciiHeader, err
	}

	return asciiHeader, nil
}
