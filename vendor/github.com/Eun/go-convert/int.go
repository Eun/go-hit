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

type toInt interface {
	Int() int
}
type toIntWithErr interface {
	Int() (int, error)
}

func (s stdRecipes) baseStructToInt(_ Converter, in reflect.Value, out *int) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
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

	if ok, i, err := genericIntConvert(in); ok {
		if err != nil {
			return err
		}
		*out = int(i)
		return nil
	}

	return fmt.Errorf("%s has no Int() function", in.Type().String())
}

func genericIntConvert(in reflect.Value) (bool, int64, error) {
	// check for struct.Int()
	if i, ok := in.Interface().(toInt); ok {
		return true, int64(i.Int()), nil
	}
	if i, ok := in.Interface().(toIntWithErr); ok {
		v, err := i.Int()
		return true, int64(v), err
	}

	// check for struct.Int8()
	if i, ok := in.Interface().(toInt8); ok {
		return true, int64(i.Int8()), nil
	}
	if i, ok := in.Interface().(toInt8WithErr); ok {
		v, err := i.Int8()
		return true, int64(v), err
	}

	// check for struct.Int16()
	if i, ok := in.Interface().(toInt16); ok {
		return true, int64(i.Int16()), nil
	}
	if i, ok := in.Interface().(toInt16WithErr); ok {
		v, err := i.Int16()
		return true, int64(v), err
	}

	// check for struct.Int32()
	if i, ok := in.Interface().(toInt32); ok {
		return true, int64(i.Int32()), nil
	}
	if i, ok := in.Interface().(toInt32WithErr); ok {
		v, err := i.Int32()
		return true, int64(v), err
	}

	// check for struct.Int64()
	if i, ok := in.Interface().(toInt64); ok {
		return true, i.Int64(), nil
	}
	if i, ok := in.Interface().(toInt64WithErr); ok {
		v, err := i.Int64()
		return true, v, err
	}

	// check for struct.Uint()
	if i, ok := in.Interface().(toUint); ok {
		return true, int64(i.Uint()), nil
	}
	if i, ok := in.Interface().(toUintWithErr); ok {
		v, err := i.Uint()
		return true, int64(v), err
	}

	// check for struct.Uint8()
	if i, ok := in.Interface().(toUint8); ok {
		return true, int64(i.Uint8()), nil
	}
	if i, ok := in.Interface().(toUint8WithErr); ok {
		v, err := i.Uint8()
		return true, int64(v), err
	}

	// check for struct.Uint16()
	if i, ok := in.Interface().(toUint16); ok {
		return true, int64(i.Uint16()), nil
	}
	if i, ok := in.Interface().(toUint16WithErr); ok {
		v, err := i.Uint16()
		return true, int64(v), err
	}

	// check for struct.Uint32()
	if i, ok := in.Interface().(toUint32); ok {
		return true, int64(i.Uint32()), nil
	}
	if i, ok := in.Interface().(toUint32WithErr); ok {
		v, err := i.Uint32()
		return true, int64(v), err
	}
	// check for struct.Uint64()
	if i, ok := in.Interface().(toUint64); ok {
		return true, int64(i.Uint64()), nil
	}
	if i, ok := in.Interface().(toUint64WithErr); ok {
		v, err := i.Uint64()
		return true, int64(v), err
	}
	return false, 0, nil
}
