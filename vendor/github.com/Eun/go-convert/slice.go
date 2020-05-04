package convert

import (
	"reflect"

	"errors"
)

func (stdRecipes) nilToSlice(_ Converter, _ NilValue, out SliceValue) error {
	out.Elem().Set(reflect.MakeSlice(reflect.SliceOf(out.Type().Elem().Elem()), 0, 0))
	return nil
}

func (stdRecipes) stringToSlice(c Converter, in string, out SliceValue) error {
	valueType := out.Type().Elem().Elem()
	switch valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size := len(in)
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), size, size)
		for i := size - 1; i >= 0; i-- {
			sl.Index(i).SetInt(int64(in[i]))
		}
		out.Elem().Set(sl)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		size := len(in)
		sl := reflect.MakeSlice(reflect.SliceOf(valueType), size, size)
		for i := size - 1; i >= 0; i-- {
			sl.Index(i).SetUint(uint64(in[i]))
		}
		out.Elem().Set(sl)
		return nil
	}
	return errors.New("no recipe")
}

func (stdRecipes) sliceToSlice(c Converter, in, out SliceValue) error {
	o := out.Elem()
	valueType := o.Type().Elem()

	length := in.Len()

	sl := reflect.MakeSlice(reflect.SliceOf(valueType), length, in.Cap())
	for i := length - 1; i >= 0; i-- {
		value := reflect.New(valueType)
		// lookup the original
		if i < o.Len() {
			orig := o.Index(i)
			if orig.IsValid() {
				value.Elem().Set(orig)
			}
		}

		if err := c.ConvertReflectValue(in.Index(i), value); err != nil {
			return err
		}
		sl.Index(i).Set(value.Elem())
	}
	o.Set(sl)
	return nil
}
