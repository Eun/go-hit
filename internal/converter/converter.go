// Package converter contains a convert.Converter for the hit package with some standard recipes.
package converter

import (
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

// Convert converts src to dst using the options, it uses the created converter in the init func.
func Convert(src, dst interface{}, options ...convert.Options) error {
	return converter.Convert(src, dst, options...)
}
