package convert

import (
	"strconv"
)

func (stdRecipes) intToFloat64(c Converter, in int, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) int8ToFloat64(c Converter, in int8, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) int16ToFloat64(c Converter, in int16, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) int32ToFloat64(c Converter, in int32, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) int64ToFloat64(c Converter, in int64, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) uintToFloat64(c Converter, in uint, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) uint8ToFloat64(c Converter, in uint8, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) uint16ToFloat64(c Converter, in uint16, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) uint32ToFloat64(c Converter, in uint32, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) uint64ToFloat64(c Converter, in uint64, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) boolToFloat64(c Converter, in bool, out *float64) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToFloat64(c Converter, in float32, out *float64) error {
	*out = float64(in)
	return nil
}
func (stdRecipes) float64ToFloat64(c Converter, in float64, out *float64) error {
	*out = in
	return nil
}
func (stdRecipes) stringToFloat64(c Converter, in string, out *float64) error {
	i, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	*out = i
	return nil
}
