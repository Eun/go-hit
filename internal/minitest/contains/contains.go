// Package contains provides functions to check if a needle is in an haystack.
package contains

import (
	"reflect"
	"strings"

	"github.com/Eun/go-convert"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/internal/misc"
)

// NoRecipeFoundError is an error that will be returned when there is no recipe found to check if needle is in an
// haystack.
type NoRecipeFoundError error

// Contains returns true if needle is in haystack, it returns NoRecipeFoundError as an error if it cannot determinate
// how it can test if needle is in haystack.
func Contains(haystack, needle interface{}) (bool, error) {
	if haystack == nil && needle == nil {
		return true, nil
	}
	if haystack == nil || needle == nil {
		return false, nil
	}

	hay := misc.GetValue(haystack)
	if !hay.IsValid() {
		return false, xerrors.New("hay value is invalid")
	}
	switch hay.Kind() {
	case reflect.String:
		return stringContains(hay.String(), needle)
	case reflect.Slice:
		return sliceContains(hay, needle)
	case reflect.Map:
		return mapContains(hay, needle)
	case reflect.Struct:
		return structContains(hay, needle)
	default:
		return false, NoRecipeFoundError(xerrors.Errorf("unsure how to determinate if %#v is inside %#v", haystack, needle))
	}
}

func stringContains(s string, needle interface{}) (bool, error) {
	var needleStr string
	if err := convert.Convert(needle, &needleStr); err != nil {
		return false, nil
	}
	return strings.Contains(s, needleStr), nil
}

func sliceContains(s reflect.Value, needle interface{}) (bool, error) {
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		needleValue := misc.MakeTypeCopy(v)
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if cmp.Equal(v, needleValue) {
			return true, nil
		}
	}
	return false, nil
}

func mapContains(m reflect.Value, needle interface{}) (bool, error) {
	for _, key := range m.MapKeys() {
		v := key.Interface()
		needleValue := misc.MakeTypeCopy(v)
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true, nil
		}
	}
	return false, nil
}

func structContains(st reflect.Value, needle interface{}) (bool, error) {
	for i := 0; i < st.NumField(); i++ {
		v := st.Type().Field(i).Name
		needleValue := misc.MakeTypeCopy(v)
		if err := convert.Convert(needle, &needleValue); err != nil {
			continue
		}
		if v == needleValue {
			return true, nil
		}
	}
	return false, nil
}
