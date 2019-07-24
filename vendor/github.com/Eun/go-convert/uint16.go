package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToUint16(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Uint16:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(uint16(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(uint16(src.Base.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(uint16(src.Base.Float())), nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(uint16(1)), nil
		}
		return reflect.ValueOf(uint16(0)), nil
	case reflect.String:
		n, err := strconv.ParseUint(src.Base.String(), 0, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint16(n)), nil
	}
	return reflect.Value{}, nil
}
