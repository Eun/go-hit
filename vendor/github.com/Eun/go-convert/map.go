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
	zeroKey := reflect.Zero(keyType)
	zeroValue := reflect.Zero(valueType)

	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), src.Base.Len())

	for _, key := range src.Base.MapKeys() {
		haveKey, err := conv.newNestedConverter().convert(key, zeroKey)
		if err != nil {
			return reflect.Value{}, err
		}

		// lookup key in dst
		zv := dst.Base.MapIndex(haveKey)
		if !zv.IsValid() {
			zv = zeroValue
		}

		haveValue := src.Base.MapIndex(key)

		newValue, err := conv.newNestedConverter().convert(haveValue, zv)
		if err != nil {
			return reflect.Value{}, err
		}

		m.SetMapIndex(haveKey, newValue)
	}

	return m, nil
}

func (conv *Converter) convertStructToMap(src, dst *convertValue) (reflect.Value, error) {
	keyType := dst.Base.Type().Key()
	valueType := dst.Base.Type().Elem()
	zeroValue := reflect.Zero(valueType)

	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), src.Base.NumField())

	for i := src.Base.NumField() - 1; i >= 0; i-- {
		field := reflect.ValueOf(src.Base.Type().Field(i).Name)
		haveValue := src.Base.Field(i)

		// lookup key in dst
		zv := dst.Base.MapIndex(field)
		if !zv.IsValid() {
			zv = zeroValue
		}

		newValue, err := conv.newNestedConverter().convert(haveValue, zv)
		if err != nil {
			return reflect.Value{}, err
		}

		m.SetMapIndex(field, newValue)
	}

	return m, nil
}
