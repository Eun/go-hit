package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToUint32(Converter, NilValue, *uint32) error {
	return nil
}
func (stdRecipes) intToUint32(c Converter, in int, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int8ToUint32(c Converter, in int8, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int16ToUint32(c Converter, in int16, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int32ToUint32(c Converter, in int32, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) int64ToUint32(c Converter, in int64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uintToUint32(c Converter, in uint, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint8ToUint32(c Converter, in uint8, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint16ToUint32(c Converter, in uint16, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) uint32ToUint32(c Converter, in uint32, out *uint32) error {
	*out = in
	return nil
}
func (stdRecipes) uint64ToUint32(c Converter, in uint64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) boolToUint32(c Converter, in bool, out *uint32) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint32(c Converter, in float32, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) float64ToUint32(c Converter, in float64, out *uint32) error {
	*out = uint32(in)
	return nil
}
func (stdRecipes) stringToUint32(c Converter, in string, out *uint32) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 32)
	if err != nil {
		return err
	}
	*out = uint32(i)
	return nil
}
func (stdRecipes) timeToUint32(c Converter, in time.Time, out *uint32) error {
	*out = uint32(in.Unix())
	return nil
}

func (s stdRecipes) structToUint32(c Converter, in StructValue, out *uint32) error {
	err := s.baseStructToUint32(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Uint32()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToUint32(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToUint32(_ Converter, in reflect.Value, out *uint32) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toUint interface {
		Uint32() uint32
	}

	// check for struct.Uint32()
	i, ok := in.Interface().(toUint)
	if ok {
		*out = i.Uint32()
		return nil
	}

	return fmt.Errorf("%s has no Uint32() function", in.Type().String())
}
