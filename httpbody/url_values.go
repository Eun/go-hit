package httpbody

import (
	"net/url"

	"github.com/Eun/go-convert"
)

// URLValues is a wrapper for url.Values but it works with HTTPBody.
type URLValues struct {
	body   *HTTPBody
	values url.Values
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v *URLValues) Get(key string) string {
	return v.values.Get(key)
}

// Values returns all the values associated with the given key.
func (v *URLValues) Values(key string) []string {
	if v.values == nil {
		return nil
	}
	return v.values[key]
}

// Set sets the key to value. It replaces any existing
// values.
func (v *URLValues) Set(key, value string) {
	v.values.Set(key, value)
	v.body.SetString(v.values.Encode())
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v *URLValues) Add(key, value string) {
	v.values.Add(key, value)
	v.body.SetString(v.values.Encode())
}

// Del deletes the values associated with key.
func (v *URLValues) Del(key string) {
	v.values.Del(key)
	v.body.SetString(v.values.Encode())
}

// ConvertRecipes contains recipes for go-convert.
func (v *URLValues) ConvertRecipes() []convert.Recipe {
	return convert.MustMakeRecipes(func(_ convert.Converter, in *URLValues, out *url.Values) error {
		*out = make(url.Values)
		for k, v := range in.values {
			(*out)[k] = make([]string, len(v))
			copy((*out)[k], v)
		}
		return nil
	})
}

// ParseURLValues takes a HTTPBody and parses the url.Values, it returns a pointer to URLValues.
func ParseURLValues(body *HTTPBody) (*URLValues, error) {
	s, err := body.String()
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(s)
	if err != nil {
		return nil, err
	}
	return &URLValues{
		body:   body,
		values: v,
	}, nil
}
