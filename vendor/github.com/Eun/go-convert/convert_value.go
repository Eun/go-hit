package convert

import (
	"fmt"
	"reflect"
	"strings"
)

type convertValue struct {
	Base    reflect.Value
	Parents []reflect.Value
}

func newConvertValue(value reflect.Value) (*convertValue, error) {
	var cv convertValue
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		cv.Parents = append(cv.Parents, value)
		value = value.Elem()
	}
	if value.IsValid() {
		cv.Base = value
	} else if size := len(cv.Parents); size > 0 {
		cv.Base = cv.Parents[size-1]
		cv.Parents = cv.Parents[:size-1]
	}

	if !cv.Base.IsValid() {
		return nil, &InvalidTypeError{src: value}
	}

	return &cv, nil
}

func constructValue(v reflect.Value, value *convertValue) (reflect.Value, error) {
	last := v
	for i := len(value.Parents) - 1; i >= 0; i-- {
		parent := value.Parents[i]
		var p reflect.Value
		switch parent.Kind() {
		case reflect.Ptr:
			p = reflect.New(parent.Elem().Type())
			p.Elem().Set(last)
		case reflect.Interface:
			p = reflect.New(parent.Elem().Type())
			p.Elem().Set(last)
			p = p.Elem()
		default:
			return reflect.Value{}, fmt.Errorf("unsure howto handle %s", parent.Kind().String())
		}
		last = p
	}
	return last, nil
}

func (v *convertValue) getFriendlyKindName(typ reflect.Type) string {
	if v.IsNil() {
		return "nil"
	}
	if typ.Kind() == reflect.Interface {
		return "interface{}"
	}
	return strings.Replace(typ.String(), " ", "", -1)
}

func (v *convertValue) getHumanName() string {
	if v.IsNil() {
		return "nil"
	}
	var sb strings.Builder
	for i := 0; i < len(v.Parents); i++ {
		switch v.Parents[i].Kind() {
		case reflect.Ptr:
			sb.WriteString("*")
		case reflect.Interface:
			sb.WriteString("interface{}(")
		}
	}

	switch v.Base.Kind() {
	case reflect.Map:
		fmt.Fprintf(&sb, "map[%s]%s", v.Base.Type().Key().Kind().String(), v.getFriendlyKindName(v.Base.Type().Elem()))
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(&sb, "[]%s", v.getFriendlyKindName(v.Base.Type().Elem()))
	default:
		sb.WriteString(v.getFriendlyKindName(v.Base.Type()))
	}

	for i := 0; i < len(v.Parents); i++ {
		if v.Parents[i].Kind() == reflect.Interface {
			sb.WriteString(")")
		}
	}

	return sb.String()
}

func (v *convertValue) IsNil() bool {
	switch v.Base.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.Base.IsNil()
	}
	return false
}
