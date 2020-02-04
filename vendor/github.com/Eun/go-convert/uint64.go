package convert

import (
	"strconv"
)

func (stdRecipes) intToUint64(c Converter, in int, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) int8ToUint64(c Converter, in int8, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) int16ToUint64(c Converter, in int16, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) int32ToUint64(c Converter, in int32, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) int64ToUint64(c Converter, in int64, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) uintToUint64(c Converter, in uint, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) uint8ToUint64(c Converter, in uint8, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) uint16ToUint64(c Converter, in uint16, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) uint32ToUint64(c Converter, in uint32, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) uint64ToUint64(c Converter, in uint64, out *uint64) error {
	*out = in
	return nil
}
func (stdRecipes) boolToUint64(c Converter, in bool, out *uint64) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint64(c Converter, in float32, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) float64ToUint64(c Converter, in float64, out *uint64) error {
	*out = uint64(in)
	return nil
}
func (stdRecipes) stringToUint64(c Converter, in string, out *uint64) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 64)
	if err != nil {
		return err
	}
	*out = uint64(i)
	return nil
}
