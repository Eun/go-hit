package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToUint(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Uint:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(uint(src.Base.Int())), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(uint(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(uint(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(uint(1)), nil
		}
		return reflect.ValueOf(uint(0)), nil
	case reflect.String:
		n, err := strconv.ParseUint(src.Base.String(), 0, strconv.IntSize)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(n)), nil
	}
	return reflect.Value{}, nil
}
