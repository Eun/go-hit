package convert

import (
	"reflect"
)

func (conv *Converter) convertToMap(src, dst *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Map:
		return conv.convertMapToMap(src, dst)
	case reflect.Struct:
		return conv.convertStructToMap(src, dst)
	}
	return reflect.Value{}, nil
}

func (conv *Converter) convertMapToMap(src, dst *convertValue) (reflect.Value, error) {
	keyType := dst.Base.Type().Key()
	valueType := dst.Base.Type().Elem()
	zeroValue := reflect.Zero(valueType)
	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), src.Base.Len())

	for _, key := range src.Base.MapKeys() {
		newKey, err := conv.newNestedConverter().convert(key, reflect.Zero(keyType))
		if err != nil {
			return reflect.Value{}, err
		}
		newValue, err := conv.newNestedConverter().convert(src.Base.MapIndex(key), zeroValue)
		if err != nil {
			return reflect.Value{}, err
		}

		m.SetMapIndex(newKey, newValue)
	}

	return m, nil
}

func (conv *Converter) convertStructToMap(src, dst *convertValue) (reflect.Value, error) {
	keyType := dst.Base.Type().Key()
	valueType := dst.Base.Type().Elem()

	zeroValue := reflect.Zero(valueType)

	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), src.Base.NumField())

	for i := src.Base.NumField() - 1; i >= 0; i-- {
		newValue, err := conv.newNestedConverter().convert(src.Base.Field(i), zeroValue)
		if err != nil {
			return reflect.Value{}, err
		}

		m.SetMapIndex(reflect.ValueOf(src.Base.Type().Field(i).Name), newValue)
	}

	return m, nil
}
