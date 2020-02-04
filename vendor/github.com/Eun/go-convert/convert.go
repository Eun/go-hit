package convert

import (
	"errors"
	"fmt"
	"reflect"
)

// Converter is the instance that will be used to convert values
type defaultConverter struct {
	options Options
}

var defaultConverterInstance = defaultConverter{
	options: Options{
		SkipUnknownFields: false,
		Recipes:           getStdRecipes(),
	},
}

// Convert converts the specified value to the specified type and returns it.
// The behavior can be influenced by using the options
// Example:
//     str, err := Convert(8, "")
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("%s\n", str.(string))
func (conv defaultConverter) Convert(src, dst interface{}, options ...Options) error {
	if dst == nil {
		return errors.New("destination type cannot be nil")
	}

	if src == nil {
		src = NilValue{}
	}

	return conv.ConvertReflectValue(reflect.ValueOf(src), reflect.ValueOf(dst), options...)
}

// MustConvert calls Convert() but panics if there is an error
func (conv defaultConverter) MustConvert(src, dstTyp interface{}, options ...Options) {
	if err := conv.Convert(src, dstTyp, options...); err != nil {
		panic(err)
	}
}

func (conv defaultConverter) MustConvertReflectValue(src, dstTyp reflect.Value, options ...Options) {
	if err := conv.ConvertReflectValue(src, dstTyp, options...); err != nil {
		panic(err)
	}
}
func (conv defaultConverter) ConvertReflectValue(src, dst reflect.Value, options ...Options) error {
	if !src.IsValid() {
		return errors.New("source is invalid")
	}
	if !dst.IsValid() {
		return errors.New("destination is invalid")
	}

	// unpackInterface src
	src = unpackInterface(src)

	if dst.Kind() != reflect.Ptr {
		return errors.New("destination must be a pointer")
	}

	if !src.IsValid() {
		// src was an interface to nothing
		return nil
	}

	// interface(string) -> string
	out := unpackInterface(dst.Elem())

	if !out.IsValid() {
		// dst was an interface to nothing
		dst.Elem().Set(src)
		return nil
	}

	// create a new temporary variable (this variable is the real type that we want to convert to)
	tmp := reflect.New(out.Type())
	tmp.Elem().Set(out)
	out = tmp

	debug("converting %s `%v' to %s\n",
		src.Type().String(),
		printValue(src),
		out.Elem().Type().String())

	if len(options) > 0 {
		conv.options.SkipUnknownFields = options[0].SkipUnknownFields
		conv.options.Recipes = append(options[0].Recipes, conv.options.Recipes...)
	}

	genericFrom := conv.getGenericType(src.Type())
	genericTo := conv.getGenericType(out.Elem().Type())

	if genericFrom != nil {
		debug2(">> generic from: %s\n", genericFrom.String())
	}

	if genericTo != nil {
		debug2(">> generic to:   %s\n", genericTo.String())
	}

	for _, recipe := range conv.options.Recipes {
		if recipe.From != src.Type() && recipe.From != genericFrom {
			debug2(">> skipping %v because src %s != %s\n", recipe.Func, recipe.From.String(), src.Type().String())
			continue
		}
		if recipe.To != out.Type() && recipe.To != genericTo {
			debug2(">> skipping %v because dst %s != %s\n", recipe.Func, recipe.To.String(), out.Type().String())
			continue
		}
		// debug2("entering %s | %s==%s==%s %s==%s==%s", recipe.Func, recipe.From.String(), src.Type().String(), genericFrom.String(), recipe.To.String(), dst.Type().String(), genericTo.String())
		if err := recipe.Func(&conv, src, out); err != nil {
			return fmt.Errorf("unable to convert %s to %s: %s", src.Type().String(), out.Elem().Type().String(), err)
		}
		debug(">> successful (1) %s `%v' to %s `%v'\n", src.Type().String(), printValue(src), out.Elem().Type().String(), printValue(out.Elem()))
		// everything was good, set the dst to the copy
		dst.Elem().Set(out.Elem())
		return nil
	}

	if src.Kind() == reflect.Ptr {
		if src.Elem().IsValid() {
			debug(">> following src because no recipe for %s to %s was found\n", src.Type().String(), dst.Type().String())
			return conv.ConvertReflectValue(src.Elem(), dst, options...)
		}
		if src.IsNil() {
			debug(">> src is nil\n")
			// make a new instance if src is nil
			return conv.ConvertReflectValue(reflect.Zero(src.Type().Elem()), dst, options...)
		}
	}

	if out.Elem().Kind() == reflect.Ptr {
		if out.Elem().Elem().IsValid() {
			debug(">> following dst because no recipe for %s to %s was found\n", src.Type().String(), out.Elem().Type().String())
			return conv.ConvertReflectValue(src, out.Elem(), options...)
		}
		debug(">> Creating new %s\n", out.Type().Elem().Elem().String())
		out = reflect.New(out.Type().Elem().Elem())
		err := conv.ConvertReflectValue(src, out, options...)
		if err != nil {
			return err
		}
		dst.Elem().Set(out)
		debug(">> successful (2) %s `%v' to %s `%v'\n", src.Type().String(), printValue(src), out.Elem().Type().String(), printValue(out.Elem()))
		return nil
	}

	return fmt.Errorf("unable to convert %s to %s: no recipe", src.Type().String(), out.Elem().Type().String())
}

func (conv *defaultConverter) getGenericType(p reflect.Type) reflect.Type {
	if p == NilType {
		return NilType
	}
	switch p.Kind() {
	case reflect.Struct:
		return StructType
	case reflect.Map:
		return MapType
	case reflect.Slice, reflect.Array:
		return SliceType
	}
	return nil
}

func (conv *defaultConverter) Options() *Options {
	return &conv.options
}

// New creates a new converter that can be used multiple times
func New(options ...Options) Converter {
	conv := defaultConverterInstance
	if len(options) > 0 {
		conv.options.SkipUnknownFields = options[0].SkipUnknownFields
		conv.options.Recipes = append(options[0].Recipes, conv.options.Recipes...)
	}
	return &conv
}

// Convert converts the specified value to the specified type and returns it.
// The behavior can be influenced by using the options
// Example:
//     var str string
//     if err := Convert(8, &str); err != nil {
//         panic(err)
//     }
//     fmt.Printf("%s\n", str)
func Convert(src, dst interface{}, options ...Options) error {
	return defaultConverterInstance.Convert(src, dst, options...)
}

// MustConvert calls Convert() but panics if there is an error
func MustConvert(src, dst interface{}, options ...Options) {
	defaultConverterInstance.MustConvert(src, dst, options...)
}

func ConvertReflectValue(src, dstTyp reflect.Value, options ...Options) error {
	return defaultConverterInstance.ConvertReflectValue(src, dstTyp, options...)
}

// MustConvertReflectValue calls MustConvertReflectValue() but panics if there is an error
func MustConvertReflectValue(src, dstTyp reflect.Value, options ...Options) {
	defaultConverterInstance.MustConvertReflectValue(src, dstTyp, options...)
}

func unpackInterface(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	return value
}

func printValue(value reflect.Value) string {
	v := value
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Ptr, reflect.Interface:
			v = v.Elem()
		default:
			if v.CanInterface() {
				return fmt.Sprintf("%v", v.Interface())
			}
			return "unknown"
		}
	}

	return ""
}
