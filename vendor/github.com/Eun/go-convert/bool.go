package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func (stdRecipes) nilToBool(Converter, NilValue, *bool) error {
	return nil
}

func (stdRecipes) intToBool(c Converter, in int, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int8ToBool(c Converter, in int8, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int16ToBool(c Converter, in int16, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int32ToBool(c Converter, in int32, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) int64ToBool(_ Converter, in int64, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uintToBool(c Converter, in uint, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint8ToBool(c Converter, in uint8, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint16ToBool(c Converter, in uint16, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint32ToBool(c Converter, in uint32, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) uint64ToBool(_ Converter, in uint64, out *bool) error {
	*out = in != 0
	return nil
}

func (stdRecipes) boolToBool(_ Converter, in bool, out *bool) error {
	*out = in
	return nil
}

func (stdRecipes) float32ToBool(_ Converter, in float32, out *bool) error {
	*out = in != 0.0
	return nil
}

func (stdRecipes) float64ToBool(_ Converter, in float64, out *bool) error {
	*out = in != 0.0
	return nil
}

func (stdRecipes) stringToBool(_ Converter, in string, out *bool) error {
	if in == "" {
		*out = false
		return nil
	}
	var err error
	*out, err = strconv.ParseBool(in)
	return err
}

func (s stdRecipes) structToBool(c Converter, in StructValue, out *bool) error {
	err := s.baseStructToBool(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Bool()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToBool(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToBool(_ Converter, in reflect.Value, out *bool) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toBool interface {
		Bool() bool
	}
	type toBoolWithErr interface {
		Bool() (bool, error)
	}

	// check for struct.String()
	if i, ok := in.Interface().(toBool); ok {
		*out = i.Bool()
		return nil
	}
	if i, ok := in.Interface().(toBoolWithErr); ok {
		var err error
		*out, err = i.Bool()
		return err
	}

	return fmt.Errorf("%s has no Bool() function", in.Type().String())
}
