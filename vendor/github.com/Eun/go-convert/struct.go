package convert

import (
	"reflect"

	"fmt"
	"strings"
)

func (stdRecipes) nilToStruct(Converter, NilValue, StructValue) error {
	return nil
}

func (stdRecipes) mapToStruct(c Converter, in MapValue, out StructValue) error {
	fieldNameValue := reflect.New(reflect.TypeOf(""))
	for _, key := range in.MapKeys() {
		// convert key
		if err := c.ConvertReflectValue(key, fieldNameValue); err != nil {
			return err
		}
		fieldName := fieldNameValue.Elem().String()

		// find the destination field with the converted value
		field := out.Elem().FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(fieldName, s)
		})
		if !field.IsValid() || !field.CanSet() {
			if c.Options().SkipUnknownFields {
				continue
			}
			return fmt.Errorf("unable to find %s in %s", fieldName, out.Elem().Type().String())
		}

		// convert value
		value := reflect.New(field.Type())
		if err := c.ConvertReflectValue(in.MapIndex(key), value); err != nil {
			return err
		}

		field.Set(value.Elem())
	}
	return nil
}

func (stdRecipes) structToStruct(c Converter, in, out StructValue) error {
	for i := in.NumField() - 1; i >= 0; i-- {
		fieldName := in.Type().Field(i).Name
		// find the destination field

		// find the destination field with the converted value
		field := out.Elem().FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(fieldName, s)
		})
		if !field.IsValid() || !field.CanSet() {
			if c.Options().SkipUnknownFields {
				continue
			}
			return fmt.Errorf("unable to find %s in %s", fieldName, out.Elem().Type().String())
		}

		// convert value
		value := reflect.New(field.Type())
		err := c.ConvertReflectValue(in.Field(i), value)
		if err != nil {
			return err
		}

		field.Set(value.Elem())
	}

	return nil
}
