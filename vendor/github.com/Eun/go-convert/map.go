package convert

import (
	"reflect"
)

func (stdRecipes) nilToMap(Converter, NilValue, MapValue) error {
	return nil
}

func (stdRecipes) mapToMap(c Converter, in MapValue, out MapValue) error {
	keyType := out.Elem().Type().Key()
	valueType := out.Elem().Type().Elem()

	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), in.Len())

	for _, key := range in.MapKeys() {
		keyValue := reflect.New(keyType)
		if err := c.ConvertReflectValue(key, keyValue); err != nil {
			return err
		}

		var valueValue reflect.Value
		// lookup key in dst
		existingValue := out.Elem().MapIndex(keyValue.Elem())
		if existingValue.IsValid() {
			valueValue = reflect.New(existingValue.Type())
			valueValue.Elem().Set(existingValue)
		} else {
			valueValue = reflect.New(valueType)
		}

		if err := c.ConvertReflectValue(in.MapIndex(key), valueValue); err != nil {
			return err
		}

		m.SetMapIndex(keyValue.Elem(), valueValue.Elem())
	}
	out.Elem().Set(m)

	return nil
}

func (stdRecipes) structToMap(c Converter, in StructValue, out MapValue) error {
	keyType := out.Elem().Type().Key()
	valueType := out.Elem().Type().Elem()

	m := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), in.NumField())

	for i := in.NumField() - 1; i >= 0; i-- {
		field := reflect.ValueOf(in.Type().Field(i).Name)

		var valueValue reflect.Value
		// lookup key in dst
		existingValue := out.Elem().MapIndex(field)
		if existingValue.IsValid() {
			valueValue = reflect.New(existingValue.Type())
			valueValue.Elem().Set(existingValue)
		} else {
			valueValue = reflect.New(valueType)
		}

		if err := c.ConvertReflectValue(in.Field(i), valueValue); err != nil {
			return err
		}

		m.SetMapIndex(field, valueValue.Elem())
	}
	out.Elem().Set(m)

	return nil
}
