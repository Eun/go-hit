package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToInt32(Converter, NilValue, *int32) error {
	return nil
}
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
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 32)
	if err != nil {
		return err
	}
	*out = int32(i)
	return nil
}
func (stdRecipes) timeToInt32(c Converter, in time.Time, out *int32) error {
	*out = int32(in.Unix())
	return nil
}

func (s stdRecipes) structToInt32(c Converter, in StructValue, out *int32) error {
	err := s.baseStructToInt32(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Int32()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToInt32(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToInt32(_ Converter, in reflect.Value, out *int32) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toInt32 interface {
		Int32() int32
	}

	// check for struct.Int32()
	i, ok := in.Interface().(toInt32)
	if ok {
		*out = i.Int32()
		return nil
	}

	return fmt.Errorf("%s has no Int32() function", in.Type().String())
}
