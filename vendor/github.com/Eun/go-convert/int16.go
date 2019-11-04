package convert

import (
	"strconv"
)

func (stdRecipes) intToInt16(c Converter, in int, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int8ToInt16(c Converter, in int8, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int16ToInt16(c Converter, in int16, out *int16) error {
	*out = in
	return nil
}
func (stdRecipes) int32ToInt16(c Converter, in int32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int64ToInt16(c Converter, in int64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uintToInt16(c Converter, in uint, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint8ToInt16(c Converter, in uint8, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint16ToInt16(c Converter, in uint16, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint32ToInt16(c Converter, in uint32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint64ToInt16(c Converter, in uint64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) boolToInt16(c Converter, in bool, out *int16) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt16(c Converter, in float32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) float64ToInt16(c Converter, in float64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) stringToInt16(c Converter, in string, out *int16) error {
	i, err := strconv.ParseInt(in, 0, 16)
	if err != nil {
		return err
	}
	*out = int16(i)
	return nil
}
