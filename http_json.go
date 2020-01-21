package hit

import (
	"encoding/json"

	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal/minitest"
)

type HTTPJson struct {
	Hit
	body *HTTPBody
}

func newHTTPJson(body *HTTPBody) *HTTPJson {
	return &HTTPJson{
		body: body,
		Hit:  body.hit,
	}
}

// Get returns the body as an interface type based on the underlying data
func (jsn *HTTPJson) Get(expression string) interface{} {
	var container interface{}
	minitest.NoError(json.NewDecoder(jsn.body.Reader()).Decode(&container))
	v, ok, err := expr.GetValue(container, expression, expr.IgnoreCase)
	minitest.NoError(err)
	if !ok {
		v = nil
	}
	return v
}

// GetAs returns the body as the specified interface type
func (jsn *HTTPJson) GetAs(expression string, container interface{}) interface{} {
	var v interface{}
	minitest.NoError(json.NewDecoder(jsn.body.Reader()).Decode(&v))
	v, ok, err := expr.GetValue(v, expression, expr.IgnoreCase)
	minitest.NoError(err)
	if !ok {
		v = nil
	}
	minitest.NoError(convert.Convert(v, container))

	return container
}

// Set sets the body to the specified json data
func (jsn *HTTPJson) Set(data interface{}) {
	buf, err := json.Marshal(data)
	minitest.NoError(err)
	jsn.body.SetBytes(buf)
}
