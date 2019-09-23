package convert

import (
	"errors"
	"fmt"
	"reflect"
)

// Converter is the instance that will be used to convert values
type Converter struct {
	Options     []Option
	source      *convertValue
	destination *convertValue
	parent      *Converter
}

var defaultConverter = Converter{}

// Convert converts the specified value to the specified type and returns it.
// The behavior can be influenced by using the options
// Example:
//     str, err := Convert(8, "")
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("%s\n", str.(string))
func (conv Converter) Convert(src, dstTyp interface{}, options ...Option) (interface{}, error) {
	if dstTyp == nil {
		return nil, errors.New("destination type cannot be nil")
	}

	conv.Options = append(conv.Options, options...)

	var err error

	if src == nil {
		var nilValue *struct{}
		src = nilValue
	}

	conv.source, err = newConvertValue(reflect.ValueOf(src))
	if err != nil {
		return reflect.Value{}, err
	}
	conv.destination, err = newConvertValue(reflect.ValueOf(dstTyp))
	if err != nil {
		return reflect.Value{}, err
	}

	n, err := conv.convert()
	if err != nil {
		return nil, err
	}

	if !n.CanInterface() {
		return nil, fmt.Errorf("cannot create interface for %s", n.Kind().String())
	}

	return n.Interface(), nil
}

// MustConvert calls Convert() but panics if there is an error
func (conv Converter) MustConvert(src, dstTyp interface{}, options ...Option) interface{} {
	v, err := conv.Convert(src, dstTyp, options...)
	if err != nil {
		panic(err)
	}
	return v
}

func (conv Converter) MustConvertReflectValue(src, dstTyp reflect.Value, options ...Option) reflect.Value {
	v, err := conv.ConvertReflectValue(src, dstTyp, options...)
	if err != nil {
		panic(err)
	}
	return v
}
func (conv Converter) ConvertReflectValue(src, dstTyp reflect.Value, options ...Option) (reflect.Value, error) {
	if !src.IsValid() {
		return reflect.Value{}, errors.New("source is invalid")
	}
	if !dstTyp.IsValid() {
		return reflect.Value{}, errors.New("destination is invalid")
	}

	conv.Options = append(conv.Options, options...)

	var err error
	conv.source, err = newConvertValue(src)
	if err != nil {
		return reflect.Value{}, err
	}
	conv.destination, err = newConvertValue(dstTyp)
	if err != nil {
		return reflect.Value{}, err
	}

	return conv.convert()
}

func (conv *Converter) convert() (reflect.Value, error) {
	n, err := conv.convertToType(conv.source, conv.destination)
	if err != nil {
		return reflect.Value{}, &Error{src: conv.source, dst: conv.destination, underlayingError: err}
	}

	if !n.IsValid() {
		return reflect.Value{}, &Error{src: conv.source, dst: conv.destination}
	}

	if conv.hasOption(Options.SkipPointers()) {
		return n, nil
	}

	return constructValue(n, conv.destination)
}

func (conv *Converter) convertToType(src, dst *convertValue) (reflect.Value, error) {
	debug("converting %s to %s\n", src.getHumanName(), dst.getHumanName())
	ok, value, err := conv.callCustomConverter(src, dst)
	if err != nil {
		return reflect.Value{}, err
	}
	if ok {
		return value, nil
	}
	switch dst.Base.Kind() {
	case reflect.Bool:
		return conv.convertToBool(src, dst)
	case reflect.Int:
		return conv.convertToInt(src, dst)
	case reflect.Int8:
		return conv.convertToInt8(src, dst)
	case reflect.Int16:
		return conv.convertToInt16(src, dst)
	case reflect.Int32:
		return conv.convertToInt32(src, dst)
	case reflect.Int64:
		return conv.convertToInt64(src, dst)
	case reflect.Uint:
		return conv.convertToUint(src, dst)
	case reflect.Uint8:
		return conv.convertToUint8(src, dst)
	case reflect.Uint16:
		return conv.convertToUint16(src, dst)
	case reflect.Uint32:
		return conv.convertToUint32(src, dst)
	case reflect.Uint64:
		return conv.convertToUint64(src, dst)
	case reflect.Float32:
		return conv.convertToFloat32(src, dst)
	case reflect.Float64:
		return conv.convertToFloat64(src, dst)
	case reflect.Map:
		return conv.convertToMap(src, dst)
	case reflect.Slice, reflect.Array:
		return conv.convertToSlice(src, dst)
	case reflect.String:
		return conv.convertToString(src, dst)
	case reflect.Struct:
		return conv.convertToStruct(src, dst)
	case reflect.Interface:
		return conv.convertToInterface(src, dst)
	}
	return reflect.Value{}, nil
}

func (conv *Converter) newNestedConverter() *Converter {
	return &Converter{
		parent:  conv,
		Options: conv.Options,
	}
}

// New creates a new converter that can be used multiple times
func New(options ...Option) *Converter {
	return &Converter{
		Options: options,
	}
}

// Convert converts the specified value to the specified type and returns it.
// The behavior can be influenced by using the options
// Example:
//     str, err := Convert(8, "")
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("%s\n", str.(string))
func Convert(src, dstTyp interface{}, options ...Option) (interface{}, error) {
	return defaultConverter.Convert(src, dstTyp, options...)
}

// MustConvert calls Convert() but panics if there is an error
func MustConvert(src, dstTyp interface{}, options ...Option) interface{} {
	return defaultConverter.MustConvert(src, dstTyp, options...)
}

func ConvertReflectValue(src, dstTyp reflect.Value, options ...Option) (reflect.Value, error) {
	return defaultConverter.ConvertReflectValue(src, dstTyp, options...)
}

// MustConvertReflectValue calls MustConvertReflectValue() but panics if there is an error
func MustConvertReflectValue(src, dstTyp reflect.Value, options ...Option) reflect.Value {
	return defaultConverter.MustConvertReflectValue(src, dstTyp, options...)
}

// GetHumanName returns a friendly human readable name for an type
// Example:
//     fmt.Println(GetHumanName(&time.Time{}))
//     prints *time.Time
func GetHumanName(v interface{}) string {
	if v == nil {
		return "nil"
	}
	s, err := newConvertValue(reflect.ValueOf(v))
	if err != nil {
		panic(err.Error())
	}
	return s.getHumanName()
}
