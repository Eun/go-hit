package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToInt(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Int:
		return src.Base, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(int(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(int(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(int(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(1), nil
		}
		return reflect.ValueOf(0), nil
	case reflect.String:
		n, err := strconv.ParseInt(src.Base.String(), 0, strconv.IntSize)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(n)), nil
	}
	return reflect.Value{}, nil
}
