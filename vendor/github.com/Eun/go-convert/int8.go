package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToInt8(Converter, NilValue, *int8) error {
	return nil
}
func (stdRecipes) intToInt8(c Converter, in int, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int8ToInt8(c Converter, in int8, out *int8) error {
	*out = in
	return nil
}
func (stdRecipes) int16ToInt8(c Converter, in int16, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int32ToInt8(c Converter, in int32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) int64ToInt8(c Converter, in int64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uintToInt8(c Converter, in uint, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint8ToInt8(c Converter, in uint8, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint16ToInt8(c Converter, in uint16, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint32ToInt8(c Converter, in uint32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) uint64ToInt8(c Converter, in uint64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) boolToInt8(c Converter, in bool, out *int8) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt8(c Converter, in float32, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) float64ToInt8(c Converter, in float64, out *int8) error {
	*out = int8(in)
	return nil
}
func (stdRecipes) stringToInt8(c Converter, in string, out *int8) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 8)
	if err != nil {
		return err
	}
	*out = int8(i)
	return nil
}
func (stdRecipes) timeToInt8(c Converter, in time.Time, out *int8) error {
	*out = int8(in.Unix())
	return nil
}

func (s stdRecipes) structToInt8(c Converter, in StructValue, out *int8) error {
	err := s.baseStructToInt8(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Int8()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToInt8(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToInt8(_ Converter, in reflect.Value, out *int8) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toInt8 interface {
		Int8() int8
	}
	type toInt8WithErr interface {
		Int8() (int8, error)
	}

	// check for struct.Int8()
	if i, ok := in.Interface().(toInt8); ok {
		*out = i.Int8()
		return nil
	}
	if i, ok := in.Interface().(toInt8WithErr); ok {
		var err error
		*out, err = i.Int8()
		return err
	}

	return fmt.Errorf("%s has no Int8() function", in.Type().String())
}
