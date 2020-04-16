package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToUint16(Converter, NilValue, *uint16) error {
	return nil
}
func (stdRecipes) intToUint16(c Converter, in int, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int8ToUint16(c Converter, in int8, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int16ToUint16(c Converter, in int16, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int32ToUint16(c Converter, in int32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) int64ToUint16(c Converter, in int64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uintToUint16(c Converter, in uint, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint8ToUint16(c Converter, in uint8, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint16ToUint16(c Converter, in uint16, out *uint16) error {
	*out = in
	return nil
}
func (stdRecipes) uint32ToUint16(c Converter, in uint32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) uint64ToUint16(c Converter, in uint64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) boolToUint16(c Converter, in bool, out *uint16) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint16(c Converter, in float32, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) float64ToUint16(c Converter, in float64, out *uint16) error {
	*out = uint16(in)
	return nil
}
func (stdRecipes) stringToUint16(c Converter, in string, out *uint16) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 16)
	if err != nil {
		return err
	}
	*out = uint16(i)
	return nil
}
func (stdRecipes) timeToUint16(c Converter, in time.Time, out *uint16) error {
	*out = uint16(in.Unix())
	return nil
}

func (s stdRecipes) structToUint16(c Converter, in StructValue, out *uint16) error {
	err := s.baseStructToUint16(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Uint16()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToUint16(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToUint16(_ Converter, in reflect.Value, out *uint16) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toUint interface {
		Uint16() uint16
	}

	// check for struct.Uint16()
	i, ok := in.Interface().(toUint)
	if ok {
		*out = i.Uint16()
		return nil
	}

	return fmt.Errorf("%s has no Uint16() function", in.Type().String())
}
