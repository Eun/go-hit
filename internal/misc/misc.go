package misc

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// GetLastArgument returns the last argument
func GetLastArgument(params []interface{}) (interface{}, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}

// GetLastStringArgument returns the last argument
func GetLastStringArgument(params []string) (string, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return "", false
}

// GetLastIntArgument returns the last argument
func GetLastIntArgument(params []int) (int, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return 0, false
}

func MakeURL(base, url string, a ...interface{}) string {
	url = fmt.Sprintf(url, a...)
	if base == "" {
		return url
	}
	if url == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(url, "/")
}

type devNullWriter struct{}

func (devNullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// DevNullWriter returns a writer to the abyss
func DevNullWriter() io.Writer {
	return devNullWriter{}
}

// StringSliceHasPrefixStringSlice returns true if haystack starts with needle
func StringSliceHasPrefixSlice(haystack, needle []string) bool {
	haySize := len(haystack)
	needleSize := len(needle)
	if needleSize > haySize {
		return false
	}
	haySize = needleSize

	for i := 0; i < haySize; i++ {
		if haystack[i] != needle[i] {
			return false
		}
	}
	return true
}

// GetElem follows all pointers and interfaces until it reaches the base value
func GetElem(r reflect.Value) reflect.Value {
	for r.IsValid() && r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	return r
}

// GetValue returns a reflect.Value and follows all pointers and interfaces until it reaches the base value
func GetValue(v interface{}) reflect.Value {
	return GetElem(reflect.ValueOf(v))
}

func GetGenericFunc(v interface{}) reflect.Value {
	r := GetValue(v)
	if !r.IsValid() {
		return r
	}
	if r.Kind() != reflect.Func {
		return reflect.Value{} // return an invalid func
	}
	return r
}

func CallGenericFunc(f reflect.Value) {
	t := f.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		in[i] = reflect.Zero(t.In(i))
	}
	f.Call(in)
}
