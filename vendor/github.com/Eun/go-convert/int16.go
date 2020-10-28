package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToInt16(Converter, NilValue, *int16) error {
	return nil
}
func (stdRecipes) intToInt16(c Converter, in int, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int8ToInt16(c Converter, in int8, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int16ToInt16(c Converter, in int16, out *int16) error {
	*out = in
	return nil
}
func (stdRecipes) int32ToInt16(c Converter, in int32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) int64ToInt16(c Converter, in int64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uintToInt16(c Converter, in uint, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint8ToInt16(c Converter, in uint8, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint16ToInt16(c Converter, in uint16, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint32ToInt16(c Converter, in uint32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) uint64ToInt16(c Converter, in uint64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) boolToInt16(c Converter, in bool, out *int16) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt16(c Converter, in float32, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) float64ToInt16(c Converter, in float64, out *int16) error {
	*out = int16(in)
	return nil
}
func (stdRecipes) stringToInt16(c Converter, in string, out *int16) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 16)
	if err != nil {
		return err
	}
	*out = int16(i)
	return nil
}
func (stdRecipes) timeToInt16(c Converter, in time.Time, out *int16) error {
	*out = int16(in.Unix())
	return nil
}

func (s stdRecipes) structToInt16(c Converter, in StructValue, out *int16) error {
	err := s.baseStructToInt16(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Int16()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToInt16(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToInt16(_ Converter, in reflect.Value, out *int16) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toInt16 interface {
		Int16() int16
	}
	type toInt16WithErr interface {
		Int16() (int16, error)
	}

	// check for struct.Int16()
	if i, ok := in.Interface().(toInt16); ok {
		*out = i.Int16()
		return nil
	}
	if i, ok := in.Interface().(toInt16WithErr); ok {
		var err error
		*out, err = i.Int16()
		return err
	}

	return fmt.Errorf("%s has no Int16() function", in.Type().String())
}
