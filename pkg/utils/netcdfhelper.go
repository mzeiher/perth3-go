package utils

import (
	"errors"

	"github.com/fhs/go-netcdf/netcdf"
)

var (
	ErrNetcdfAttributeNotFound = errors.New("attribute not found")
	ErrNetcdfDimensionNotFound = errors.New("dimension not found")
)

func NetcdfGetDimensionFromVariable(name string, variable *netcdf.Var) (*netcdf.Dim, error) {
	dims, err := variable.Dims()
	if err != nil {
		return nil, err
	}
	for _, dim := range dims {
		dimname, err := dim.Name()
		if err != nil {
			return nil, err
		}
		if dimname == name {
			return &dim, nil
		}
	}
	return nil, ErrNetcdfDimensionNotFound

}

func NetcdfGetFloat32FromAttribute(name string, variable *netcdf.Var) ([]float32, error) {
	attr, err := NetcdfGetAttribute(name, variable)
	if err != nil {
		return nil, err
	}

	attrLen, err := attr.Len()
	if err != nil {
		return nil, err
	}
	buffer := make([]float32, attrLen)
	err = attr.ReadFloat32s(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil

}

func NetcdfGetStringFromAttribute(name string, variable *netcdf.Var) (string, error) {
	attr, err := NetcdfGetAttribute(name, variable)
	if err != nil {
		return "", err
	}
	attrLen, err := attr.Len()
	if err != nil {
		return "", err
	}
	buffer := make([]byte, attrLen)
	err = attr.ReadBytes(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func NetcdfGetAttribute(name string, variable *netcdf.Var) (*netcdf.Attr, error) {
	numberAttributes, err := variable.NAttrs()
	if err != nil {
		return nil, err
	}
	if numberAttributes == 0 {
		return nil, ErrNetcdfAttributeNotFound
	}
	for i := 0; i < numberAttributes; i++ {
		attr, err := variable.AttrN(i)
		if err != nil {
			return nil, err
		}
		if attr.Name() == name {
			return &attr, nil
		}
	}
	return nil, ErrNetcdfAttributeNotFound
}
