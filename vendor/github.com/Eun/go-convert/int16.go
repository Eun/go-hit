package convert

import (
	"errors"
	"reflect"
	"strconv"
)

func (conv *Converter) convertToInt16(src, _ *convertValue) (reflect.Value, error) {
	if src.IsNil() {
		return reflect.Value{}, errors.New("source cannot be nil")
	}
	switch src.Base.Kind() {
	case reflect.Int16:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(int16(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(int16(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(int16(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(int16(1)), nil
		}
		return reflect.ValueOf(int16(0)), nil
	case reflect.String:
		n, err := strconv.ParseInt(src.Base.String(), 0, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int16(n)), nil
	}
	return reflect.Value{}, nil
}
