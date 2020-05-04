package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToFloat64(Converter, NilValue, *float64) error {
	return nil
}
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
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	*out = i
	return nil
}

func (stdRecipes) timeToFloat64(c Converter, in time.Time, out *float64) error {
	*out = float64(in.Unix()) + float64(in.Nanosecond())/1000000000
	return nil
}

func (s stdRecipes) structToFloat64(c Converter, in StructValue, out *float64) error {
	err := s.baseStructToFloat64(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Float64()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToFloat64(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToFloat64(_ Converter, in reflect.Value, out *float64) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toFloat64 interface {
		Float64() float64
	}
	type toFloat64WithErr interface {
		Float64() (float64, error)
	}

	// check for struct.Float64()
	if i, ok := in.Interface().(toFloat64); ok {
		*out = i.Float64()
		return nil
	}
	if i, ok := in.Interface().(toFloat64WithErr); ok {
		var err error
		*out, err = i.Float64()
		return err
	}

	return fmt.Errorf("%s has no Float64() function", in.Type().String())
}
