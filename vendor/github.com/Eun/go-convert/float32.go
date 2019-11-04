package convert

import (
	"strconv"
)

func (stdRecipes) intToFloat32(c Converter, in int, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) int8ToFloat32(c Converter, in int8, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) int16ToFloat32(c Converter, in int16, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) int32ToFloat32(c Converter, in int32, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) int64ToFloat32(c Converter, in int64, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) uintToFloat32(c Converter, in uint, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) uint8ToFloat32(c Converter, in uint8, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) uint16ToFloat32(c Converter, in uint16, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) uint32ToFloat32(c Converter, in uint32, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) uint64ToFloat32(c Converter, in uint64, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) boolToFloat32(c Converter, in bool, out *float32) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToFloat32(c Converter, in float32, out *float32) error {
	*out = in
	return nil
}
func (stdRecipes) float64ToFloat32(c Converter, in float64, out *float32) error {
	*out = float32(in)
	return nil
}
func (stdRecipes) stringToFloat32(c Converter, in string, out *float32) error {
	i, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return err
	}
	*out = float32(i)
	return nil
}
