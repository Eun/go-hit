package httpbody

import (
	"encoding/xml"

	"github.com/Eun/go-hit/internal/minitest"
)

type HttpBodyXml struct {
	// hit.Hit
	body *HttpBody
}

func newHttpBodyXml(body *HttpBody) *HttpBodyXml {
	return &HttpBodyXml{
		body: body,
		// Hit:  body.hit,
	}
}

// Set sets the body to the specified json data
func (jsn *HttpBodyXml) Set(data interface{}) {
	buf, err := xml.Marshal(data)
	minitest.Panic.NoError(err)
	jsn.body.SetBytes(buf)
}
