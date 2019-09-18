package convert

import (
	"reflect"
)

type Option interface {
	private()
}

type skipUnknownFieldsOption struct{}

func (skipUnknownFieldsOption) private() {}

type skipPointersOption struct{}

func (skipPointersOption) private() {}

type customConverterOption struct {
	From reflect.Type
	To   reflect.Type
	Func reflect.Value
}

func (*customConverterOption) private() {}

func customConverter(converterFunc interface{}) Option {
	if converterFunc == nil {
		panic("converterFunc cannot be nil")
	}
	// figure out the parameters
	v := reflect.ValueOf(converterFunc)
	if v.Kind() != reflect.Func {
		panic("converterFunc must be a function")
	}

	t := v.Type()

	if t.NumIn() != 1 {
		panic("only one parameter can be passed (the type to convert from)")
	}

	if t.NumOut() != 2 {
		panic("exactly 2 result values must be provided (type to create to and error)")
	}

	if t.In(0).Kind() == reflect.Ptr {
		panic("parameter cannot be a ptr")
	}

	if t.Out(0).Kind() == reflect.Ptr {
		panic("result cannot be a ptr")
	}

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	if !t.Out(1).Implements(errorInterface) {
		panic("error must be implemented")
	}

	return &customConverterOption{
		From: t.In(0),
		To:   t.Out(0),
		Func: v,
	}
}

var Options = struct {
	SkipUnknownFields func() Option
	SkipPointers      func() Option
	CustomConverter   func(interface{}) Option
}{
	SkipUnknownFields: func() Option {
		return skipUnknownFieldsOption{}
	},
	SkipPointers: func() Option {
		return skipPointersOption{}
	},
	CustomConverter: customConverter,
}

func (conv *Converter) hasOption(option Option) bool {
	for _, opt := range conv.Options {
		if opt == option {
			return true
		}
	}
	return false
}

func (conv *Converter) callCustomConverter(from, to *convertValue) (bool, reflect.Value, error) {
	for _, opt := range conv.Options {
		if v, ok := opt.(*customConverterOption); ok {
			if v.From == from.Base.Type() && v.To == to.Base.Type() {
				result := v.Func.Call([]reflect.Value{from.Base})
				if !result[1].IsNil() {
					return false, reflect.Value{}, result[1].Interface().(error)
				}
				return true, result[0], nil
			}
		}
	}
	return false, reflect.Value{}, nil
}
