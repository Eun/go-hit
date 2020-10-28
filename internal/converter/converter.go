package converter

import (
	"reflect"

	"net/url"

	"github.com/Eun/go-convert"
)

//nolint:gochecknoglobals
var converter convert.Converter

//nolint:gochecknoinits
// create a new Converter with custom converts.
func init() {
	converter = convert.New(convert.Options{
		Recipes: convert.MustMakeRecipes(
			func(_ convert.Converter, in url.Userinfo, out *url.Userinfo) error {
				if pass, ok := in.Password(); ok {
					*out = *url.UserPassword(in.Username(), pass)
				} else {
					*out = *url.User(in.Username())
				}
				return nil
			},
		),
	})
}

func Convert(src, dst interface{}, options ...convert.Options) error {
	return converter.Convert(src, dst, options...)
}
func MustConvert(src, dst interface{}, options ...convert.Options) {
	converter.MustConvert(src, dst, options...)
}
func ConvertReflectValue(src, dst reflect.Value, options ...convert.Options) error {
	return converter.ConvertReflectValue(src, dst, options...)
}
func MustConvertReflectValue(src, dst reflect.Value, options ...convert.Options) {
	converter.MustConvertReflectValue(src, dst, options...)
}
