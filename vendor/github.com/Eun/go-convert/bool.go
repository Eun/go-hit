package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToBool(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Bool:
		return src.Base, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(src.Base.Int() != 0), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(src.Base.Uint() != 0), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(src.Base.Float() != 0.0), nil
	case reflect.String:
		b, err := strconv.ParseBool(src.Base.String())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	}
	return reflect.Value{}, nil
}
