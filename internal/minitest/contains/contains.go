package contains

import (
	"reflect"
	"strings"

	"github.com/Eun/go-convert"
	"github.com/google/go-cmp/cmp"
	"github.com/mohae/deepcopy"
)

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
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		needleValue := deepcopy.Copy(v)
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if cmp.Equal(v, needleValue) {
			return true
		}
	}
	return false
}

func mapContains(m reflect.Value, needle interface{}) bool {
	for _, key := range m.MapKeys() {
		v := key.Interface()
		needleValue := deepcopy.Copy(v)
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
	for i := 0; i < st.NumField(); i++ {
		v := st.Type().Field(i).Name
		needleValue := deepcopy.Copy(v)
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true
		}
	}
	return false
}
