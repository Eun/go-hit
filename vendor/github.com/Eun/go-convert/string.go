package convert

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func (conv *Converter) convertToString(src, dst *convertValue) (reflect.Value, error) {
	if src.IsNil() {
		return reflect.Value{}, errors.New("source cannot be nil")
	}
	switch src.Base.Kind() {
	case reflect.String:
		return src.Base, nil
	case reflect.Bool:
		return reflect.ValueOf(strconv.FormatBool(src.Base.Bool())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(strconv.FormatInt(src.Base.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(strconv.FormatUint(src.Base.Uint(), 10)), nil
	case reflect.Float32:
		return reflect.ValueOf(strconv.FormatFloat(src.Base.Float(), 'f', 6, 32)), nil
	case reflect.Float64:
		return reflect.ValueOf(strconv.FormatFloat(src.Base.Float(), 'f', 6, 64)), nil
	case reflect.Array, reflect.Slice:
		return convertSliceToString(src, dst)
	}
	return reflect.Value{}, nil
}

func convertSliceToString(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Type().Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var sb strings.Builder
		for i := 0; i < src.Base.Len(); i++ {
			sb.WriteRune(rune(src.Base.Index(i).Int()))
		}
		return reflect.ValueOf(sb.String()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var sb strings.Builder
		for i := 0; i < src.Base.Len(); i++ {
			sb.WriteRune(rune(src.Base.Index(i).Uint()))
		}
		return reflect.ValueOf(sb.String()), nil
	}
	return reflect.Value{}, nil
}
