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
			zv := conv.zeroSliceValue(dst, i, valueType, zeroValue)
			v, err := conv.newNestedConverter().convert(src.Base.Index(i), zv)
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

// zeroSliceValue returns the zero value for the specific map value
func (conv *Converter) zeroSliceValue(dst *convertValue, index int, valueType reflect.Type, fallbackValue reflect.Value) reflect.Value {
	if index >= dst.Base.Len() {
		return fallbackValue
	}

	dstType := dst.Base.Index(index)
	if !dstType.IsValid() {
		return fallbackValue
	}
	return dstType
}
