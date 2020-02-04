package convert

import (
	"strconv"
)

func (stdRecipes) intToUint32(c Converter, in int, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int8ToUint32(c Converter, in int8, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int16ToUint32(c Converter, in int16, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int32ToUint32(c Converter, in int32, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int64ToUint32(c Converter, in int64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uintToUint32(c Converter, in uint, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint8ToUint32(c Converter, in uint8, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint16ToUint32(c Converter, in uint16, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint32ToUint32(c Converter, in uint32, out *uint32) error {
	*out = in
	return nil
}
func (stdRecipes) uint64ToUint32(c Converter, in uint64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) boolToUint32(c Converter, in bool, out *uint32) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint32(c Converter, in float32, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) float64ToUint32(c Converter, in float64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) stringToUint32(c Converter, in string, out *uint32) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 32)
	if err != nil {
		return err
	}
	*out = uint32(i)
	return nil
}
