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

		zv := conv.zeroMapValue(dst, key, keyType, valueType, zeroValue)

		newValue, err := conv.newNestedConverter().convert(src.Base.MapIndex(key), zv)
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
		field := reflect.ValueOf(src.Base.Type().Field(i).Name)
		zv := conv.zeroMapValue(dst, field, keyType, valueType, zeroValue)

		newValue, err := conv.newNestedConverter().convert(src.Base.Field(i), zv)
		if err != nil {
			return reflect.Value{}, err
		}

		m.SetMapIndex(field, newValue)
	}

	return m, nil
}

// zeroMapValue returns the zero value for the specific map value
func (conv *Converter) zeroMapValue(dst *convertValue, wantedKey reflect.Value, keyType, valueType reflect.Type, fallbackValue reflect.Value) reflect.Value {
	if valueType.Kind() != reflect.Interface {
		return fallbackValue
	}

	// if the wantedKey type is not the wantedKey type that we expect
	// try to convert it, if it fails, return fallback
	if keyType != wantedKey.Type() {
		var err error
		wantedKey, err = conv.convert(wantedKey, reflect.Zero(keyType))
		if err != nil {
			return fallbackValue
		}
	}

	dstType := dst.Base.MapIndex(wantedKey)
	if !dstType.IsValid() {
		return fallbackValue
	}
	return dstType
}
