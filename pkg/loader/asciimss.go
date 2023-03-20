package loader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mzeiher/perth3-go/pkg/mss"
)

func CreateNewAsciiMSSLoader(path string) (MeanSeaSurfaceLoader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	return &asciiMSSFile{
		reader: reader,
	}, nil

}

type asciiMSSHeader struct {
	longitudeMin float64
	longitudeMax float64
	latitudeMin  float64
	latitudeMax  float64

	latRes float64
	lonRes float64
}

type asciiMSSFile struct {
	MeanSeaSurfaceLoader
	reader *bufio.Reader
}

func (a *asciiMSSFile) GetMSSData() (*mss.MedianSeaSurfaceData, error) {
	header, err := a.ParseHeader()
	if err != nil {
		return nil, err
	}

	gridX := int((header.longitudeMax-header.longitudeMin)/header.lonRes) + 1
	gridY := int((header.latitudeMax-header.latitudeMin)/header.latRes) + 1

	data, err := a.ParseData(gridY, gridX)
	if err != nil {
		return nil, err
	}

	return &mss.MedianSeaSurfaceData{
		MinLat: header.latitudeMin,
		MaxLat: header.latitudeMax,
		MinLon: header.longitudeMin,
		MaxLon: header.longitudeMax,

		LonRes: header.lonRes,
		LatRes: header.latRes,

		GridX: gridX,
		GridY: gridY,

		Data: data,
	}, nil
}

func (a *asciiMSSFile) ParseData(gridY int, gridX int) ([][]float32, error) {

	gridData := make([][]float32, gridY)

	numberEntries := gridX * gridY
	currentY := 0
	currentEntry := 0
	gridData[currentY] = make([]float32, gridX)
	for {

		line, err := a.reader.ReadString('\n')
		entries := strings.Fields(line)
		if err == nil && len(entries) == 0 {
			continue
		}
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
			if currentEntry >= gridX*(currentY+1) {
				currentY = currentY + 1
				gridData[currentY] = make([]float32, gridX)
			}

			entryParsed, err := strconv.ParseFloat(entry, 64)
			if err != nil {
				return nil, err
			}
			xPos := ((currentY * gridX) - currentEntry) * -1
			gridData[currentY][xPos] = float32(entryParsed)

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

func (a *asciiMSSFile) ParseHeader() (asciiMSSHeader, error) {
	header := asciiMSSHeader{}

	headerLine, err := a.reader.ReadString('\n')
	if err != nil {
		return header, err
	}

	fields := strings.Fields(headerLine)

	if len(fields) != 6 {
		return header, fmt.Errorf("invalid header, insufficient fields")
	}

	if _, err = fmt.Sscanf(fields[0], "%f", &header.latitudeMin); err != nil {
		return header, err
	}
	if _, err = fmt.Sscanf(fields[1], "%f", &header.latitudeMax); err != nil {
		return header, err
	}

	if _, err = fmt.Sscanf(fields[2], "%f", &header.longitudeMin); err != nil {
		return header, err
	}
	if _, err = fmt.Sscanf(fields[3], "%f", &header.longitudeMax); err != nil {
		return header, err
	}

	if _, err = fmt.Sscanf(fields[4], "%f", &header.latRes); err != nil {
		return header, err
	}
	if _, err = fmt.Sscanf(fields[5], "%f", &header.lonRes); err != nil {
		return header, err
	}

	return header, nil

}
