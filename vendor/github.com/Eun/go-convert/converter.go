package convert

import (
	"reflect"
)

// NilValue represents a nil value to convert (from/to)
type NilValue struct {
	reflect.Value
}

// MapValue represents a map value to convert (from/to)
type MapValue struct {
	reflect.Value
}

// StructValue represents a struct value to convert (from/to)
type StructValue struct {
	reflect.Value
}

// SliceValue represents a slice value to convert (from/to)
type SliceValue struct {
	reflect.Value
}

// NilType can be used to specify a recipe with the source/destination with a nil value
var NilType = reflect.TypeOf((*NilValue)(nil)).Elem()

// MapType can be used to specify a recipe with the source/destination with a map type
var MapType = reflect.TypeOf((*MapValue)(nil)).Elem()

// StructType can be used to specify a recipe with the source/destination with a struct type
var StructType = reflect.TypeOf((*StructValue)(nil)).Elem()

// SliceType can be used to specify a recipe with the source/destination with a slice type
var SliceType = reflect.TypeOf((*SliceValue)(nil)).Elem()

// Converter is the instance that will be used to convert values
type Converter interface {
	Options() *Options
	Convert(src, dst interface{}, options ...Options) error
	MustConvert(src, dst interface{}, options ...Options)
	ConvertReflectValue(src, dst reflect.Value, options ...Options) error
	MustConvertReflectValue(src, dst reflect.Value, options ...Options)
}

func isGenericType(t reflect.Type) bool {
	switch t {
	case NilType:
		return true
	case MapType:
		return true
	case StructType:
		return true
	case SliceType:
		return true
	}
	return false
}

func getGenericWrapper(v reflect.Type) (wrapParam, bool) {
	switch v {
	case NilType:
		return func(value reflect.Value) reflect.Value {
			return reflect.ValueOf(NilValue{value})
		}, true
	case MapType:
		return func(value reflect.Value) reflect.Value {
			return reflect.ValueOf(MapValue{value})
		}, true
	case StructType:
		return func(value reflect.Value) reflect.Value {
			return reflect.ValueOf(StructValue{value})
		}, true
	case SliceType:
		return func(value reflect.Value) reflect.Value {
			return reflect.ValueOf(SliceValue{value})
		}, true
	}
	return func(value reflect.Value) reflect.Value {
		return value
	}, false
}
