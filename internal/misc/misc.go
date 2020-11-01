// Package misc provides some helper functions that are used in hit.
package misc

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// MakeURL joins two urls safely.
// example: MakeURL("http://example.com", "%s.html", "index") will return "http://example.com/index.html"
func MakeURL(base, url string, a ...interface{}) string {
	if len(a) > 0 { // only format if a is provided
		url = fmt.Sprintf(url, a...)
	}
	if base == "" {
		return url
	}
	if url == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(url, "/")
}

// GetElem follows all pointers and interfaces until it reaches the base value.
func GetElem(r reflect.Value) reflect.Value {
	for r.IsValid() && r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	return r
}

// GetValue returns a reflect.Value and follows all pointers and interfaces until it reaches the base value.
func GetValue(v interface{}) reflect.Value {
	return GetElem(reflect.ValueOf(v))
}

// MakeTypeCopy creates a copy of src, with the default value(s) for this type.
func MakeTypeCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(original, cpy, false)

	// Return the copy as an interface.
	return cpy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func copyRecursive(original, cpy reflect.Value, zero bool) {
	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, cpy.Elem(), zero)
	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue, zero)
		cpy.Set(copyValue)
	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}
		// Go through each field of the struct and copy it.
		for i := 0; i < original.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.  I'm not sure why.  -mohae
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			copyRecursive(original.Field(i), cpy.Field(i), true)
		}
	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i), true)
		}
	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue, true)
			copyKey := MakeTypeCopy(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}
	default:
		if zero {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		cpy.Set(original)
	}
}
