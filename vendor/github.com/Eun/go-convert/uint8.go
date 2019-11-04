package convert

import (
	"strconv"
)

func (stdRecipes) intToUint8(c Converter, in int, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int8ToUint8(c Converter, in int8, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int16ToUint8(c Converter, in int16, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int32ToUint8(c Converter, in int32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int64ToUint8(c Converter, in int64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uintToUint8(c Converter, in uint, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint8ToUint8(c Converter, in uint8, out *uint8) error {
	*out = in
	return nil
}
func (stdRecipes) uint16ToUint8(c Converter, in uint16, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint32ToUint8(c Converter, in uint32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint64ToUint8(c Converter, in uint64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) boolToUint8(c Converter, in bool, out *uint8) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint8(c Converter, in float32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) float64ToUint8(c Converter, in float64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) stringToUint8(c Converter, in string, out *uint8) error {
	i, err := strconv.ParseUint(in, 0, 8)
	if err != nil {
		return err
	}
	*out = uint8(i)
	return nil
}
