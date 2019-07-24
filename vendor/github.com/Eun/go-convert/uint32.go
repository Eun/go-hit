package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToUint32(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Uint32:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(uint32(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
		return reflect.ValueOf(uint32(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(uint32(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(uint32(1)), nil
		}
		return reflect.ValueOf(uint32(0)), nil
	case reflect.String:
		n, err := strconv.ParseUint(src.Base.String(), 0, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint32(n)), nil
	}
	return reflect.Value{}, nil
}
