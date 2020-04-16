package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (stdRecipes) nilToString(Converter, NilValue, *string) error {
	return nil
}

func (stdRecipes) intToString(c Converter, in int, out *string) error {
	*out = strconv.FormatInt(int64(in), 10)
	return nil
}

func (stdRecipes) int8ToString(c Converter, in int8, out *string) error {
	*out = strconv.FormatInt(int64(in), 10)
	return nil
}

func (stdRecipes) int16ToString(c Converter, in int16, out *string) error {
	*out = strconv.FormatInt(int64(in), 10)
	return nil
}

func (stdRecipes) int32ToString(c Converter, in int32, out *string) error {
	*out = strconv.FormatInt(int64(in), 10)
	return nil
}

func (stdRecipes) int64ToString(_ Converter, in int64, out *string) error {
	*out = strconv.FormatInt(in, 10)
	return nil
}

func (stdRecipes) uintToString(c Converter, in uint, out *string) error {
	*out = strconv.FormatUint(uint64(in), 10)
	return nil
}

func (stdRecipes) uint8ToString(c Converter, in uint8, out *string) error {
	*out = strconv.FormatUint(uint64(in), 10)
	return nil
}

func (stdRecipes) uint16ToString(c Converter, in uint16, out *string) error {
	*out = strconv.FormatUint(uint64(in), 10)
	return nil
}

func (stdRecipes) uint32ToString(c Converter, in uint32, out *string) error {
	*out = strconv.FormatUint(uint64(in), 10)
	return nil
}

func (stdRecipes) uint64ToString(_ Converter, in uint64, out *string) error {
	*out = strconv.FormatUint(in, 10)
	return nil
}

func (stdRecipes) boolToString(c Converter, in bool, out *string) error {
	*out = strconv.FormatBool(in)
	return nil
}

func (stdRecipes) float32ToString(_ Converter, in float32, out *string) error {
	if out != nil && len(*out) > 0 {
		*out = fmt.Sprintf(*out, in)
		return nil
	}
	*out = strconv.FormatFloat(float64(in), 'f', 6, 32)
	return nil
}

func (stdRecipes) float64ToString(_ Converter, in float64, out *string) error {
	if out != nil && len(*out) > 0 {
		*out = fmt.Sprintf(*out, in)
		return nil
	}
	*out = strconv.FormatFloat(in, 'f', 6, 64)
	return nil
}

func (stdRecipes) stringToString(_ Converter, in string, out *string) error {
	*out = in
	return nil
}

func (stdRecipes) intSliceToString(_ Converter, in []int, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) int8SliceToString(c Converter, in []int8, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) int16SliceToString(c Converter, in []int16, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) int32SliceToString(c Converter, in []int32, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) int64SliceToString(c Converter, in []int64, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) uintSliceToString(c Converter, in []uint, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) uint8SliceToString(c Converter, in []uint8, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) uint16SliceToString(c Converter, in []uint16, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) uint32SliceToString(c Converter, in []uint32, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) uint64SliceToString(c Converter, in []uint64, out *string) error {
	var sb strings.Builder
	for _, i := range in {
		sb.WriteRune(rune(i))
	}
	*out = sb.String()
	return nil
}

func (stdRecipes) timeToString(c Converter, in time.Time, out *string) error {
	*out = in.String()
	return nil
}

func (s stdRecipes) structToString(c Converter, in StructValue, out *string) error {
	err := s.baseStructToString(c, in.Value, out)
	if err == nil {
		return err
	}

	// test for *struct.String()
	v := reflect.New(in.Type())
	v.Elem().Set(in.Value)
	if s.baseStructToString(c, v, out) == nil {
		return nil
	}
	return err
}

func (s stdRecipes) baseStructToString(_ Converter, in reflect.Value, out *string) error {
	if !in.CanInterface() {
		return errors.New("unable to make interface")
	}
	type toString interface {
		String() string
	}

	// check for struct.String()
	i, ok := in.Interface().(toString)
	if ok {
		*out = i.String()
		return nil
	}

	return fmt.Errorf("%s has no String() function", in.Type().String())
}
