package hit

import (
	"encoding/json"

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
func (jsn *HTTPJson) Get() (container interface{}) {
	minitest.NoError(json.NewDecoder(jsn.body.Reader()).Decode(&container))
	return container
}

// GetAs returns the body as the specified interface type
func (jsn *HTTPJson) GetAs(container interface{}) interface{} {
	minitest.NoError(json.NewDecoder(jsn.body.Reader()).Decode(container))
	return container
}

// Set sets the body to the specified json data
func (jsn *HTTPJson) Set(data interface{}) {
	buf, err := json.Marshal(data)
	minitest.NoError(err)
	jsn.body.SetBytes(buf)
}
