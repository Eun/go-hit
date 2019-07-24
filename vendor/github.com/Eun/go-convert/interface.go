package convert

import (
	"reflect"
)

func (conv *Converter) convertToInterface(src, _ *convertValue) (reflect.Value, error) {
	if src.Base.CanInterface() {
		return reflect.ValueOf(src.Base.Interface()), nil
	}
	return reflect.Value{}, nil
}
