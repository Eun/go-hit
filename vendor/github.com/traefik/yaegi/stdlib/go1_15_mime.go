// Code generated by 'yaegi extract mime'. DO NOT EDIT.

// +build go1.15,!go1.16

package stdlib

import (
	"mime"
	"reflect"
)

func init() {
	Symbols["mime"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"AddExtensionType":         reflect.ValueOf(mime.AddExtensionType),
		"BEncoding":                reflect.ValueOf(mime.BEncoding),
		"ErrInvalidMediaParameter": reflect.ValueOf(&mime.ErrInvalidMediaParameter).Elem(),
		"ExtensionsByType":         reflect.ValueOf(mime.ExtensionsByType),
		"FormatMediaType":          reflect.ValueOf(mime.FormatMediaType),
		"ParseMediaType":           reflect.ValueOf(mime.ParseMediaType),
		"QEncoding":                reflect.ValueOf(mime.QEncoding),
		"TypeByExtension":          reflect.ValueOf(mime.TypeByExtension),

		// type definitions
		"WordDecoder": reflect.ValueOf((*mime.WordDecoder)(nil)),
		"WordEncoder": reflect.ValueOf((*mime.WordEncoder)(nil)),
	}
}
