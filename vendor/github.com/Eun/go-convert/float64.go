package convert

import (
	"reflect"
	"strconv"
)

func (conv *Converter) convertToFloat64(src, _ *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(float64(src.Base.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(float64(src.Base.Uint())), nil
	case reflect.Float32:
		return reflect.ValueOf(float64(src.Base.Float())), nil
	case reflect.Float64:
		return src.Base, nil
	case reflect.Bool:
		if src.Base.Bool() {
			return reflect.ValueOf(float64(1.0)), nil
		}
		return reflect.ValueOf(float64(0.0)), nil
	case reflect.String:
		f, err := strconv.ParseFloat(src.Base.String(), 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f), nil
	}
	return reflect.Value{}, nil
}
