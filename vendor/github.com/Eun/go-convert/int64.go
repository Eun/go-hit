package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToInt64(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Int64:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return reflect.ValueOf(int64(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(int64(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(int64(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(int64(1)), nil
		}
		return reflect.ValueOf(int64(0)), nil
	case reflect.String:
		n, err := strconv.ParseInt(src.Base.String(), 0, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(n), nil
	}
	return reflect.Value{}, nil
}
