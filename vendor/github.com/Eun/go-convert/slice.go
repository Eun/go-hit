package convert

import (
	"reflect"
)

func (conv *Converter) convertToSlice(src, dst *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Slice, reflect.Array:
		valueType := dst.Base.Type().Elem()
		zeroValue := reflect.Zero(valueType)
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), src.Base.Len(), src.Base.Cap())
		for i := src.Base.Len() - 1; i >= 0; i-- {
			v, err := conv.newNestedConverter().convert(src.Base.Index(i), zeroValue)
			if err != nil {
				return reflect.Value{}, err
			}
			sl.Index(i).Set(v)
		}
		return sl, nil
	case reflect.String:
		return conv.convertStringToSlice(src, dst)
	}

	return reflect.Value{}, nil
}

func (conv *Converter) convertStringToSlice(src, dst *convertValue) (reflect.Value, error) {
	switch dst.Base.Type().Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueType := dst.Base.Type().Elem()
		str := src.Base.String()
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), len(str), len(str))
		for i := src.Base.Len() - 1; i >= 0; i-- {
			sl.Index(i).SetInt(int64(str[i]))
		}
		return sl, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueType := dst.Base.Type().Elem()
		str := src.Base.String()
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), len(str), len(str))
		for i := src.Base.Len() - 1; i >= 0; i-- {
			sl.Index(i).SetUint(uint64(str[i]))
		}
		return sl, nil
	}

	return reflect.Value{}, nil
}
