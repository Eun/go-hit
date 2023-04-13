package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToUint8(Converter, NilValue, *uint8) error {
	return nil
}
func (stdRecipes) intToUint8(c Converter, in int, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int8ToUint8(c Converter, in int8, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int16ToUint8(c Converter, in int16, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int32ToUint8(c Converter, in int32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) int64ToUint8(c Converter, in int64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uintToUint8(c Converter, in uint, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint8ToUint8(c Converter, in uint8, out *uint8) error {
	*out = in
	return nil
}
func (stdRecipes) uint16ToUint8(c Converter, in uint16, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint32ToUint8(c Converter, in uint32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) uint64ToUint8(c Converter, in uint64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) boolToUint8(c Converter, in bool, out *uint8) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint8(c Converter, in float32, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) float64ToUint8(c Converter, in float64, out *uint8) error {
	*out = uint8(in)
	return nil
}
func (stdRecipes) stringToUint8(c Converter, in string, out *uint8) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 8)
	if err != nil {
		return err
	}
	*out = uint8(i)
	return nil
}

func (stdRecipes) timeToUint8(c Converter, in time.Time, out *uint8) error {
	*out = uint8(in.Unix())
	return nil
}

func (s stdRecipes) structToUint8(c Converter, in StructValue, out *uint8) error {
	err := s.baseStructToUint8(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Uint8()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToUint8(c, v, out) == nil {
		return nil
	}
	return err
}

type toUint8 interface {
	Uint8() uint8
}
type toUint8WithErr interface {
	Uint8() (uint8, error)
}

func (s stdRecipes) baseStructToUint8(_ Converter, in reflect.Value, out *uint8) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}

	// check for struct.Uint8()
	if i, ok := in.Interface().(toUint8); ok {
		*out = i.Uint8()
		return nil
	}
	if i, ok := in.Interface().(toUint8WithErr); ok {
		var err error
		*out, err = i.Uint8()
		return err
	}

	if ok, i, err := genericIntConvert(in); ok {
		if err != nil {
			return err
		}
		*out = uint8(i)
		return nil
	}

	return fmt.Errorf("%s has no Uint8() function", in.Type().String())
}
