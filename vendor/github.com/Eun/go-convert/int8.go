package convert

import (
	"strconv"
)

func (stdRecipes) intToInt8(c Converter, in int, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int8ToInt8(c Converter, in int8, out *int8) error {
	*out = in
	return nil
}
func (stdRecipes) int16ToInt8(c Converter, in int16, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int32ToInt8(c Converter, in int32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int64ToInt8(c Converter, in int64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uintToInt8(c Converter, in uint, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint8ToInt8(c Converter, in uint8, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint16ToInt8(c Converter, in uint16, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint32ToInt8(c Converter, in uint32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint64ToInt8(c Converter, in uint64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) boolToInt8(c Converter, in bool, out *int8) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt8(c Converter, in float32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) float64ToInt8(c Converter, in float64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) stringToInt8(c Converter, in string, out *int8) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 8)
	if err != nil {
		return err
	}
	*out = int8(i)
	return nil
}
