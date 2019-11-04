package convert

import (
	"strconv"
)

func (stdRecipes) intToUint16(c Converter, in int, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int8ToUint16(c Converter, in int8, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int16ToUint16(c Converter, in int16, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int32ToUint16(c Converter, in int32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int64ToUint16(c Converter, in int64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uintToUint16(c Converter, in uint, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint8ToUint16(c Converter, in uint8, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint16ToUint16(c Converter, in uint16, out *uint16) error {
	*out = in
	return nil
}
func (stdRecipes) uint32ToUint16(c Converter, in uint32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint64ToUint16(c Converter, in uint64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) boolToUint16(c Converter, in bool, out *uint16) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint16(c Converter, in float32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) float64ToUint16(c Converter, in float64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) stringToUint16(c Converter, in string, out *uint16) error {
	i, err := strconv.ParseUint(in, 0, 16)
	if err != nil {
		return err
	}
	*out = uint16(i)
	return nil
}
