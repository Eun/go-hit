package convert

import (
	"strconv"
)

func (stdRecipes) intToInt(c Converter, in int, out *int) error {
	*out = in
	return nil
}
func (stdRecipes) int8ToInt(c Converter, in int8, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int16ToInt(c Converter, in int16, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int32ToInt(c Converter, in int32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int64ToInt(c Converter, in int64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uintToInt(c Converter, in uint, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint8ToInt(c Converter, in uint8, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint16ToInt(c Converter, in uint16, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint32ToInt(c Converter, in uint32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint64ToInt(c Converter, in uint64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) boolToInt(c Converter, in bool, out *int) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt(c Converter, in float32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) float64ToInt(c Converter, in float64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) stringToInt(c Converter, in string, out *int) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 32)
	if err != nil {
		return err
	}
	*out = int(i)
	return nil
}
