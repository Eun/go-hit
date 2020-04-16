package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToUint(Converter, NilValue, *uint) error {
	return nil
}
func (stdRecipes) intToUint(c Converter, in int, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) int8ToUint(c Converter, in int8, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) int16ToUint(c Converter, in int16, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) int32ToUint(c Converter, in int32, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) int64ToUint(c Converter, in int64, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) uintToUint(c Converter, in uint, out *uint) error {
	*out = in
	return nil
}
func (stdRecipes) uint8ToUint(c Converter, in uint8, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) uint16ToUint(c Converter, in uint16, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) uint32ToUint(c Converter, in uint32, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) uint64ToUint(c Converter, in uint64, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) boolToUint(c Converter, in bool, out *uint) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToUint(c Converter, in float32, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) float64ToUint(c Converter, in float64, out *uint) error {
	*out = uint(in)
	return nil
}
func (stdRecipes) stringToUint(c Converter, in string, out *uint) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseUint(in, 0, 32)
	if err != nil {
		return err
	}
	*out = uint(i)
	return nil
}
func (stdRecipes) timeToUint(c Converter, in time.Time, out *uint) error {
	*out = uint(in.Unix())
	return nil
}

func (s stdRecipes) structToUint(c Converter, in StructValue, out *uint) error {
	err := s.baseStructToUint(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Uint()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToUint(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToUint(_ Converter, in reflect.Value, out *uint) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toUint interface {
		Uint() uint
	}

	// check for struct.Uint()
	i, ok := in.Interface().(toUint)
	if ok {
		*out = i.Uint()
		return nil
	}

	return fmt.Errorf("%s has no Uint() function", in.Type().String())
}
