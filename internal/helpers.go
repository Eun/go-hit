package internal

import (
	"reflect"
	"strings"

	"github.com/Eun/go-convert"
)

// GetElem follows all pointers and interfaces until it reaches the base value
func GetElem(r reflect.Value) reflect.Value {
	for r.IsValid() && r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	return r
}

// GetValue returns a internal.Value and follows all pointers and interfaces until it reaches the base value
func GetValue(v interface{}) reflect.Value {
	return GetElem(reflect.ValueOf(v))
}

func Contains(heystack interface{}, needle interface{}) bool {
	if heystack == nil && needle == nil {
		return true
	}
	if heystack == nil || needle == nil {
		return false
	}

	hey := reflect.ValueOf(heystack)
	switch hey.Kind() {
	case reflect.String:
		return stringContains(hey.String(), needle)
	case reflect.Slice:
		return sliceContains(hey, needle)
	case reflect.Map:
		return mapContains(hey, needle)
	case reflect.Struct:
		return structContains(hey, needle)
	}

	return false
}

func stringContains(s string, needle interface{}) bool {
	var needleStr string
	if err := convert.Convert(needle, &needleStr); err != nil {
		return false
	}
	return strings.Contains(s, needleStr)
}

func sliceContains(s reflect.Value, needle interface{}) bool {
	for i := s.Len() - 1; i >= 0; i-- {
		v := s.Index(i).Interface()
		needleValue := v
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true
		}
	}
	return false
}

func mapContains(m reflect.Value, needle interface{}) bool {
	for _, key := range m.MapKeys() {
		v := key.Interface()
		needleValue := v
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true
		}
	}
	return false
}

func structContains(st reflect.Value, needle interface{}) bool {
	for i := st.NumField() - 1; i >= 0; i-- {
		v := st.Type().Field(i).Name
		needleValue := v
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true
		}
	}
	return false
}
