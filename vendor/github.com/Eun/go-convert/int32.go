package convert

import (
	"strconv"
)

func (stdRecipes) intToInt32(c Converter, in int, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) int8ToInt32(c Converter, in int8, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) int16ToInt32(c Converter, in int16, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) int32ToInt32(c Converter, in int32, out *int32) error {
	*out = in
	return nil
}
func (stdRecipes) int64ToInt32(c Converter, in int64, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) uintToInt32(c Converter, in uint, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) uint8ToInt32(c Converter, in uint8, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) uint16ToInt32(c Converter, in uint16, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) uint32ToInt32(c Converter, in uint32, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) uint64ToInt32(c Converter, in uint64, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) boolToInt32(c Converter, in bool, out *int32) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt32(c Converter, in float32, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) float64ToInt32(c Converter, in float64, out *int32) error {
	*out = int32(in)
	return nil
}
func (stdRecipes) stringToInt32(c Converter, in string, out *int32) error {
	i, err := strconv.ParseInt(in, 0, 32)
	if err != nil {
		return err
	}
	*out = int32(i)
	return nil
}
