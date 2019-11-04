package convert

import "reflect"

type NilValue struct{}
type MapValue struct{}
type StructValue struct{}
type SliceValue struct{}

var NilType = reflect.TypeOf((*NilValue)(nil)).Elem()
var MapType = reflect.TypeOf((*MapValue)(nil)).Elem()
var StructType = reflect.TypeOf((*StructValue)(nil)).Elem()
var SliceType = reflect.TypeOf((*SliceValue)(nil)).Elem()

// Converter is the instance that will be used to convert values
type Converter interface {
	Options() *Options
	Convert(src, dst interface{}, options ...Options) error
	MustConvert(src, dst interface{}, options ...Options)
	ConvertReflectValue(src, dst reflect.Value, options ...Options) error
	MustConvertReflectValue(src, dst reflect.Value, options ...Options)
}
