package hit

import (
	"net/http"
	"reflect"

	"strings"

	"net/url"

	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/httpbody"
)

//nolint:gochecknoglobals
var converter convert.Converter

//nolint:gochecknoinits
// create a new converter with custom converts
func init() {
	converter = convert.New(convert.Options{
		Recipes: []convert.Recipe{
			// convert http.Header to different map types (map[string]string, map[string]interface{}, ...)
			convert.MustMakeRecipe(func(c convert.Converter, in http.Header, out convert.MapValue) error {
				m := make(map[string]string)
				for key, value := range in {
					m[key] = strings.Join(value, "")
				}
				return c.ConvertReflectValue(reflect.ValueOf(m), out.Value)
			}),
			// convert different map types to http.Header
			convert.MustMakeRecipe(func(c convert.Converter, in convert.MapValue, out *http.Header) error {
				(*out) = make(map[string][]string)

				iter := in.MapRange()
				for iter.Next() {
					var key string
					var value string
					if err := c.ConvertReflectValue(iter.Key(), reflect.ValueOf(&key)); err != nil {
						return err
					}
					if err := c.ConvertReflectValue(iter.Value(), reflect.ValueOf(&value)); err != nil {
						return err
					}
					(*out)[key] = []string{value}
				}
				return nil
			}),
			convert.MustMakeRecipe(func(_ convert.Converter, in url.Userinfo, out *url.Userinfo) error {
				if pass, ok := in.Password(); ok {
					*out = *url.UserPassword(in.Username(), pass)
				} else {
					*out = *url.User(in.Username())
				}
				return nil
			}),

			convert.MustMakeRecipe(func(_ convert.Converter, in httpbody.HttpBody, out *httpbody.HttpBody) error {
				// just copy
				*out = in
				return nil
			}),
		},
	})
}
