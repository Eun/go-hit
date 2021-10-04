package httpbody

import (
	"encoding/xml"
)

// HTTPBodyXML provides XML functions for the HTTPBody.
type HTTPBodyXML struct { //nolint:revive //ignore type name will be used as httpbody.HTTPBodyXML by other packages
	body *HTTPBody
}

func newHTTPBodyXML(body *HTTPBody) *HTTPBodyXML {
	return &HTTPBodyXML{
		body: body,
	}
}

// Set sets the body to the specified json data.
func (jsn *HTTPBodyXML) Set(data interface{}) error {
	buf, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	jsn.body.SetBytes(buf)
	return nil
}
