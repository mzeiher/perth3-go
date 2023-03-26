package loader

import (
	"errors"

	"github.com/mzeiher/perth3-go/pkg/loader/constituentdata"
	"github.com/mzeiher/perth3-go/pkg/loader/dtu16ascii"
)

var ErrNoLoaderFound = errors.New("no loader found for selected format")

var loader map[string]constituentdata.CreateLoaderFunction = make(map[string]constituentdata.CreateLoaderFunction)

func init() {
	loader["dtu16ascii"] = dtu16ascii.CreateDTU16Loader
}

func GetLoader(format string, filePath string) (constituentdata.ConstituentDataLoader, error) {
	if loader[format] == nil {
		return nil, ErrNoLoaderFound
	}
	return loader[format](filePath)
}
