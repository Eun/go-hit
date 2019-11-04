package convert

import (
	"reflect"

	"errors"
)

func (stdRecipes) stringToSlice(c Converter, in, out reflect.Value) error {
	valueType := out.Type().Elem().Elem()
	switch valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str := in.String()
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), len(str), len(str))
		for i := in.Len() - 1; i >= 0; i-- {
			sl.Index(i).SetInt(int64(str[i]))
		}
		out.Elem().Set(sl)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str := in.String()
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), len(str), len(str))
		for i := in.Len() - 1; i >= 0; i-- {
			sl.Index(i).SetUint(uint64(str[i]))
		}
		out.Elem().Set(sl)
		return nil
	}
	return errors.New("no recipe")
}

func (stdRecipes) sliceToSlice(c Converter, in, out reflect.Value) error {
	out = out.Elem()
	valueType := out.Type().Elem()

	length := in.Len()

	sl := reflect.MakeSlice(reflect.SliceOf(valueType), length, in.Cap())
	for i := length - 1; i >= 0; i-- {
		value := reflect.New(valueType)
		// lookup the original
		if i < out.Len() {
			orig := out.Index(i)
			if orig.IsValid() {
				value.Elem().Set(orig)
			}
		}

		if err := c.ConvertReflectValue(in.Index(i), value); err != nil {
			return err
		}
		sl.Index(i).Set(value.Elem())
	}
	out.Set(sl)
	return nil
}

func (stdRecipes) nilToSlice(c Converter, _, out reflect.Value) error {
	out.Elem().Set(reflect.MakeSlice(reflect.SliceOf(out.Type().Elem().Elem()), 0, 0))
	return nil
}
