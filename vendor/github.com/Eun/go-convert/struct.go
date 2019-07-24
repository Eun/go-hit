package convert

import (
	"fmt"
	"reflect"
	"strings"
)

func (conv *Converter) convertToStruct(src, dst *convertValue) (reflect.Value, error) {
	switch src.Base.Kind() {
	case reflect.Map:
		return conv.convertMapToStruct(src, dst)
	case reflect.Struct:
		return conv.convertStructToStruct(src, dst)
	}
	return reflect.Value{}, nil
}

func (conv *Converter) convertMapToStruct(src, dst *convertValue) (reflect.Value, error) {
	st := reflect.New(dst.Base.Type())
	zeroString := reflect.ValueOf("")
	for _, key := range src.Base.MapKeys() {
		// convert key
		fieldNameValue, err := conv.newNestedConverter().convert(key, zeroString)
		if err != nil {
			return reflect.Value{}, err
		}
		fieldName := fieldNameValue.String()

		// find the destination field with the converted value
		field := st.Elem().FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(fieldName, s)
		})
		if !field.IsValid() || !field.CanSet() {
			if conv.hasOption(Options.SkipUnknownFields()) {
				continue
			}
			return reflect.Value{}, fmt.Errorf("unable to find %s in %s", fieldName, dst.getHumanName())
		}

		newValue, err := conv.newNestedConverter().convert(src.Base.MapIndex(key), reflect.Zero(field.Type()))
		if err != nil {
			return reflect.Value{}, err
		}

		field.Set(newValue)
	}
	return st.Elem(), nil
}

func (conv *Converter) convertStructToStruct(src, dst *convertValue) (reflect.Value, error) {
	st := reflect.New(dst.Base.Type())

	for i := src.Base.NumField() - 1; i >= 0; i-- {
		fieldName := src.Base.Type().Field(i).Name
		// find the destination field

		// find the destination field with the converted value
		field := st.Elem().FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(fieldName, s)
		})
		if !field.IsValid() || !field.CanSet() {
			if conv.hasOption(Options.SkipUnknownFields()) {
				continue
			}
			return reflect.Value{}, fmt.Errorf("unable to find %s in %s", fieldName, dst.getHumanName())
		}

		newValue, err := conv.newNestedConverter().convert(src.Base.Field(i), field)
		if err != nil {
			return reflect.Value{}, err
		}

		field.Set(newValue)
	}

	return st.Elem(), nil
}
