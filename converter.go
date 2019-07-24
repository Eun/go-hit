package hit

import (
	"net/http"
	"strings"

	"github.com/Eun/go-convert"
)

var converter *convert.Converter

// create a new converter with custom converts
func init() {
	converter = convert.New(
		convert.Options.CustomConverter(func(from http.Header) (map[string]string, error) {
			m := make(map[string]string)
			for key, value := range from {
				m[key] = strings.Join(value, "")
			}
			return m, nil
		}),
	)
}
