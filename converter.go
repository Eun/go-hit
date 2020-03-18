package hit

import (
	"net/http"
	"reflect"

	"strings"

	"github.com/Eun/go-convert"
	"golang.org/x/xerrors"
)

var converter convert.Converter

// create a new converter with custom converts
func init() {
	converter = convert.New(convert.Options{
		Recipes: []convert.Recipe{
			// convert http.Header to different map types (map[string]string, map[string]interface{}, ...)
			{
				From: reflect.TypeOf(http.Header{}),
				To:   convert.MapType,
				Func: func(c convert.Converter, in reflect.Value, out reflect.Value) error {
					if !in.CanInterface() {
						return xerrors.New("CanInterface() == false")
					}
					hdr, ok := in.Interface().(http.Header)
					if !ok {
						return xerrors.New("source is not a http.Header")
					}
					if hdr == nil {
						return xerrors.New("source is nil")
					}
					m := make(map[string]string)
					for key, value := range hdr {
						m[key] = strings.Join(value, "")
					}
					return c.ConvertReflectValue(reflect.ValueOf(m), out)
				},
			},
		},
	})
}
