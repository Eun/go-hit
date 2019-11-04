package convert

import (
	"time"

	"github.com/araddon/dateparse"
)

func (stdRecipes) intToTime(c Converter, in int, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) int8ToTime(c Converter, in int8, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) int16ToTime(c Converter, in int16, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) int32ToTime(c Converter, in int32, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) int64ToTime(_ Converter, in int64, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) uintToTime(c Converter, in uint, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) uint8ToTime(c Converter, in uint8, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) uint16ToTime(c Converter, in uint16, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) uint32ToTime(c Converter, in uint32, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) uint64ToTime(_ Converter, in uint64, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) float32ToTime(_ Converter, in float32, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}

func (stdRecipes) float64ToTime(_ Converter, in float64, out *time.Time) error {
	*out = time.Unix(int64(in), 0)
	return nil
}
func (stdRecipes) stringToTime(_ Converter, in string, out *time.Time) error {
	var err error
	*out, err = dateparse.ParseAny(in)
	return err
}
