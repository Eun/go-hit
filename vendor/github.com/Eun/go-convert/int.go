package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (stdRecipes) nilToInt(Converter, NilValue, *int) error {
	return nil
}
func (stdRecipes) intToInt(c Converter, in int, out *int) error {
	*out = in
	return nil
}
func (stdRecipes) int8ToInt(c Converter, in int8, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int16ToInt(c Converter, in int16, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int32ToInt(c Converter, in int32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) int64ToInt(c Converter, in int64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uintToInt(c Converter, in uint, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint8ToInt(c Converter, in uint8, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint16ToInt(c Converter, in uint16, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint32ToInt(c Converter, in uint32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) uint64ToInt(c Converter, in uint64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) boolToInt(c Converter, in bool, out *int) error {
	switch in {
	case true:
		*out = 1
	default:
		*out = 0
	}
	return nil
}

func (stdRecipes) float32ToInt(c Converter, in float32, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) float64ToInt(c Converter, in float64, out *int) error {
	*out = int(in)
	return nil
}
func (stdRecipes) stringToInt(c Converter, in string, out *int) error {
	if in == "" {
		*out = 0
		return nil
	}
	i, err := strconv.ParseInt(in, 0, 32)
	if err != nil {
		return err
	}
	*out = int(i)
	return nil
}
func (stdRecipes) timeToInt(c Converter, in time.Time, out *int) error {
	*out = int(in.Unix())
	return nil
}

func (s stdRecipes) structToInt(c Converter, in StructValue, out *int) error {
	err := s.baseStructToInt(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.Int()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToInt(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToInt(_ Converter, in reflect.Value, out *int) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toInt interface {
		Int() int
	}
	type toIntWithErr interface {
		Int() (int, error)
	}

	// check for struct.Int()
	if i, ok := in.Interface().(toInt); ok {
		*out = i.Int()
		return nil
	}
	if i, ok := in.Interface().(toIntWithErr); ok {
		var err error
		*out, err = i.Int()
		return err
	}

	return fmt.Errorf("%s has no Int() function", in.Type().String())
}
