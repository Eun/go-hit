package convert

import (
	"errors"
	"reflect"
)

func (conv *Converter) convertToInterface(src, _ *convertValue) (reflect.Value, error) {
	if src.IsNil() {
		return reflect.Value{}, errors.New("source cannot be nil")
	}
	if src.Base.CanInterface() {
		return reflect.ValueOf(src.Base.Interface()), nil
	}
	return reflect.Value{}, nil
}
