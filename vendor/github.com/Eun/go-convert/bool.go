package convert

import (
	"strconv"
)

func (stdRecipes) intToBool(c Converter, in int, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int8ToBool(c Converter, in int8, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int16ToBool(c Converter, in int16, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int32ToBool(c Converter, in int32, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int64ToBool(_ Converter, in int64, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uintToBool(c Converter, in uint, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint8ToBool(c Converter, in uint8, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint16ToBool(c Converter, in uint16, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint32ToBool(c Converter, in uint32, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint64ToBool(_ Converter, in uint64, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) boolToBool(_ Converter, in bool, out *bool) error {
	*out = in
	return nil
}

func (stdRecipes) float32ToBool(_ Converter, in float32, out *bool) error {
	*out = in != 0.0
	return nil
}

func (stdRecipes) float64ToBool(_ Converter, in float64, out *bool) error {
	*out = in != 0.0
	return nil
}
func (stdRecipes) stringToBool(_ Converter, in string, out *bool) error {
	var err error
	*out, err = strconv.ParseBool(in)
	return err
}
