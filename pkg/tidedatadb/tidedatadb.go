/*
This package provides the functions to read and write a constituent database for quick lookup of amplitude and phase
the binary data is structured as a netcdf file with each constituent it's own variable
*/
package tidedatadb

import (
	"errors"
	"io/fs"
	"os"

	"github.com/fhs/go-netcdf/netcdf"
)

type FileMode netcdf.FileMode

const (
	MODE_READONLY FileMode = FileMode(netcdf.NOWRITE)
)

type Dimensions struct {
	MinLat        float32
	MaxLat        float32
	MinLon        float32
	MaxLon        float32
	ResolutionLat float32
	ResolutionLon float32
	GridXSize     uint64
	GridYSize     uint64
}

type TideDataDB struct {
	file *netcdf.Dataset
}

func (t *TideDataDB) Close() error {
	return t.file.Close()
}

// open or creates a new tide data db
// IMPORTANT! currently not threadsafe!
func OpenTideDataDb(filePath string, mode FileMode) (*TideDataDB, error) {
	_, err := os.Stat(filePath)
	var file netcdf.Dataset
	// if file does not exist, create a new file
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		file, err = netcdf.CreateFile(filePath, netcdf.NETCDF4|netcdf.FileMode(mode))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		file, err = netcdf.OpenFile(filePath, netcdf.NETCDF4|netcdf.FileMode(mode))
		if err != nil {
			return nil, err
		}
	}

	return &TideDataDB{
		file: &file,
	}, nil
}
