package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToFloat32(Converter, NilValue, *float32) error {
	return nil
}
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
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return err
	}
	*out = float32(i)
	return nil
}

func (stdRecipes) timeToFloat32(c Converter, in time.Time, out *float32) error {
	*out = float32(in.Unix())
	return nil
}

func (s stdRecipes) structToFloat32(c Converter, in StructValue, out *float32) error {
	err := s.baseStructToFloat32(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Float32()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToFloat32(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToFloat32(_ Converter, in reflect.Value, out *float32) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toFloat32 interface {
		Float32() float32
	}

	// check for struct.Float32()
	i, ok := in.Interface().(toFloat32)
	if ok {
		*out = i.Float32()
		return nil
	}

	return fmt.Errorf("%s has no Float32() function", in.Type().String())
}
