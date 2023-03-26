/*
This package provides the functions to read and write a constituent database for quick lookup of amplitude and phase
the binary data is structured as follows. All data in BIG_ENDIAN format

	byte 0..7  Preamble ([t,i,d,e,d,t,a,VERSION])
	byte 8     Type of Data entry (byte)
	byte 9-11  Padding
	byte 12..n []DataEntry

after the preamble the constituents are written one after another, to search for a constituent read the header, compare the constituent,
if it does not match advance the pointer (HEADER_SIZE_BYTE + header.GRIDX*header.GRIDY) entries relative to the offset of the header start,
repeat until constituent is found or EOF is reached

the constituent data is packaged and aligned as followed

	// data entry
	byte 0...3  Length incl. Data (uint32)
	byte 4...7  Min Latitude of Grid (float32)
	byte 8...11 Max Latitude of Grid (float32)
	byte 12..15 Min Longitude of Grid (float32)
	byte 16..19 Max Longitude of Grid (float32)
	byte 20..23 Resolution Latitude of Grid (float32)
	byte 24..27 Resolution Longitude of Grid (float32)
	byte 28..31 Size of grid in Y direction (latitude) (uint32)
	byte 32..35 Size of grid in X direction (longitude) (uint32)
	byte 36..39 Value of "undefined" (float32)

	byte 40...m Type info
	After the header the data begins and is formatted depending on the type (example for type=Constituent)
	byte m...m+15  [amplitude (float32), phase(float32)] at position y=0, x=0 (minlat, minlon)
	byte m+16...m+15  [amplitude (float32), phase(float32)] at position y=0, x=1 (minlat, minlon+1)
	...
	byte n-16...n  [amplitude (float32), phase(float32)] at position y=GridYSize-1, x=GridXSize-1 (maxlat, maxlon)
*/
package tidedatadb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/fs"
	"os"
)

var (
	ErrInvalidFile     = errors.New("invalid file, magic number mismatch")
	ErrInvalidDataType = errors.New("datatype of database incompatible")
)

var (
	DB_VERSION    byte = 1
	MAGIC_NUMBER       = []byte{'t', 'i', 'd', 'e', 'd', 't', 'a', DB_VERSION}
	PREAMBLE_SIZE      = 12
)

type DataType byte

const (
	TYPE_CONSTITUENT DataType = iota
	TYPE_TIDE_VALUE_CACHE
)

type DataEntry struct {
	Length        uint32
	MinLat        float32
	MaxLat        float32
	MinLon        float32
	MaxLon        float32
	ResolutionLat float32
	ResolutionLon float32
	GridXSize     uint32
	GridYSize     uint32
	UndefValue    float32
}

type TideDataDB struct {
	file *os.File
	Type DataType
}

func (t *TideDataDB) Close() error {
	return t.file.Close()
}

// open or creates a new tide data db
// IMPORTANT! currently not threadsafe!
func OpenTideDataDb(filePath string, dataType DataType, flags int) (*TideDataDB, error) {
	_, err := os.Stat(filePath)
	var file *os.File
	// if file does not exist, create a new file
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		file, err = os.OpenFile(filePath, flags, 0666)
		if err != nil {
			return nil, err
		}

		// write magic number
		err = binary.Write(file, binary.BigEndian, MAGIC_NUMBER)
		if err != nil {
			return nil, err
		}

		// write type of db
		binary.Write(file, binary.BigEndian, dataType)
		if err != nil {
			return nil, err
		}

		// write padding
		binary.Write(file, binary.BigEndian, []byte{0, 0, 0})
		if err != nil {
			return nil, err
		}

		// reset position
		file.Seek(0, 0)
	} else if err != nil {
		return nil, err
	} else {
		file, err = os.OpenFile(filePath, flags, 0666)
		if err != nil {
			return nil, err
		}
		header := make([]byte, 8)
		err = binary.Read(file, binary.BigEndian, &header)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(header, MAGIC_NUMBER) {
			return nil, ErrInvalidFile
		}
	}

	// read type of data
	_, err = file.Seek(8, 0)
	if err != nil {
		return nil, err
	}
	var fileDataType DataType
	err = binary.Read(file, binary.BigEndian, &fileDataType)
	if err != nil {
		return nil, err
	} else if fileDataType != dataType {
		return nil, ErrInvalidDataType
	}

	return &TideDataDB{
		file: file,
		Type: dataType,
	}, nil
}

// func AppendConstituent(file *os.File, header Header, data [][]float32) (int, error) {
// 	// file.
// 	io.EOF
// }
