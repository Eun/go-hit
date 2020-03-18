package hit

import (
	"github.com/Eun/go-convert"
)

var converter convert.Converter

// create a new converter with custom converts
func init() {
	converter = convert.New(convert.Options{
		Recipes: convert.MustMakeRecipes(
		// func(_ convert.Converter, in http.Header, out *map[string]string) error {
		// 	*out = make(map[string]string)
		// 	for key, value := range in {
		// 		(*out)[key] = strings.Join(value, "")
		// 	}
		// 	return nil
		// },
		),
	})
}
