package hit

import (
	"net/http"

	"strings"

	"github.com/Eun/go-hit/internal/converter"
)

// ISendHeaders provides methods to send header/trailer.
type ISendHeaders interface {
	// Add adds the specified value to the specified request header.
	//
	// Usage:
	//     Send().Headers("Set-Cookie").Add("foo=bar")
	//     Send().Headers("Set-Cookie").Add("foo=bar", "baz=taz")
	Add(value ...interface{}) IStep
}

type sendHeadersValueCallback func(hit *hitImpl) http.Header

type sendHeader struct {
	cleanPath     callPath
	valueCallback sendHeadersValueCallback
	name          string
}

func newSendHeaders(clearPath callPath, valueCallback sendHeadersValueCallback, name string) ISendHeaders {
	return &sendHeader{
		cleanPath:     clearPath,
		valueCallback: valueCallback,
		name:          name,
	}
}

func (hdr *sendHeader) Add(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: hdr.cleanPath.Push("Add", values),
		Exec: func(hit *hitImpl) error {
			for _, value := range values {
				var s string
				if err := converter.Convert(value, &s); err != nil {
					return err
				}
				if strings.EqualFold(hdr.name, "host") {
					hit.Request().Host = s
					return nil
				}
				hdr.valueCallback(hit).Add(hdr.name, s)
			}
			return nil
		},
	}
}
