package hit

import (
	"io"
	"strconv"

	"github.com/Eun/go-hit/errortrace"
)

// IDebugRequest defines the debug functions that are available for the http request.
type IDebugRequest interface {
	IStep

	// Method prints the Request's Method
	Method() IStep

	// URL prints the Request's URL
	//
	// The argument can be used to narrow down the print path
	//
	// Usage:
	//     Debug().Request().URL()       // print the whole URL struct
	URL() IStep

	// Proto prints the Request's Proto
	Proto() IStep

	// ProtoMajor prints the Request's ProtoMajor
	ProtoMajor() IStep

	// ProtoMinor prints the Request's ProtoMinor
	ProtoMinor() IStep

	// ContentLength prints the Request's ContentLength
	ContentLength() IStep

	// TransferEncoding prints the Request's TransferEncoding
	TransferEncoding() IStep

	// Host prints the Request's Host
	Host() IStep

	// Headers prints the Request's Headers
	//
	// If you omit the argument it will print all headers.
	//
	// Usage:
	//     Debug().Request().Headers()               // all headers
	//     Debug().Request().Headers("Content-Type") // only Content-Type
	Headers(header ...string) IStep

	// Trailers prints the Request's Trailers
	//
	// If you omit the argument it will print all trailers.
	//
	// Usage:
	//     Debug().Request().Trailers()               // all trailers
	//     Debug().Request().Trailers("Content-Type") // only Content-Type
	Trailers(header ...string) IStep

	// Body prints the Request's Body
	//
	// The argument can be used to narrow down the print path
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Debug().Request().Body()                         // will print the whole Body
	//     Debug().Request().Body().JSON()                  // will print the whole Body as JSON data
	//     Debug().Request().Body().JSON().JQ(".ID")        // will print 10
	//     Debug().Request().Body().JSON().JQ(".Name")      // will print "Name"
	//     Debug().Request().Body().JSON().JQ(".Roles")     // will print ["Admin", "User"]
	//     Debug().Request().Body().JSON().JQ(".Roles.0")   // will print "Admin"
	Body() IDebugBody
}

type debugRequest struct {
	cp    callPath
	debug *debug
}

func newDebugRequest(cp callPath, d *debug) *debugRequest {
	return &debugRequest{
		cp:    cp,
		debug: d,
	}
}

func (*debugRequest) trace() *errortrace.ErrorTrace {
	return nil
}

func (*debugRequest) when() StepTime {
	return BeforeExpectStep
}

func (d *debugRequest) data(hit Hit) map[string]interface{} {
	var urlData map[string]interface{}

	if u := hit.Request().URL; u != nil {
		urlData = make(map[string]interface{})
		urlData["Scheme"] = u.Scheme
		urlData["Opaque"] = u.Opaque

		urlData["Host"] = u.Host
		urlData["Path"] = u.Path
		urlData["Query"] = d.debug.getMap(u.Query())
		urlData["Hostname"] = u.Hostname()
		urlData["RequestURI"] = u.RequestURI()
		urlData["Port"] = u.Port()
		urlData["RawPath"] = u.RawPath
		urlData["EscapedPath"] = u.EscapedPath()
		urlData["ForceQuery"] = u.ForceQuery
		urlData["RawQuery"] = u.RawQuery
		urlData["Fragment"] = u.Fragment
		urlData["String"] = u.String()

		if u.User == nil {
			urlData["User"] = nil
		} else {
			password, _ := u.User.Password()
			urlData["User"] = map[string]interface{}{
				"Username": u.User.Username(),
				"Password": password,
			}
		}
	}

	return map[string]interface{}{
		"Method":           hit.Request().Method,
		"URL":              urlData,
		"Proto":            hit.Request().Proto,
		"ProtoMajor":       hit.Request().ProtoMajor,
		"ProtoMinor":       hit.Request().ProtoMinor,
		"ContentLength":    hit.Request().ContentLength,
		"TransferEncoding": hit.Request().TransferEncoding,
		"Host":             hit.Request().Host,
		"Headers":          d.debug.getMap(hit.Request().Header),
		"Trailers":         d.debug.getMap(hit.Request().Trailer),
		"Body":             hit.Request().Body().GetBestFittingObject(),
	}
}

func (d *debugRequest) exec(hit *hitImpl) error {
	return d.debug.print(d.debug.out(hit), d.data(hit))
}

func (d *debugRequest) callPath() callPath {
	return d.cp
}

func (d *debugRequest) Method() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("Method", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), hit.Request().Method)
			return err
		},
	}
}

func (d *debugRequest) URL() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("URL", nil),
		Exec: func(hit *hitImpl) error {
			return d.debug.print(d.debug.out(hit), hit.Request().URL)
		},
	}
}

func (d *debugRequest) Proto() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("Proto", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), hit.Request().Proto)
			return err
		},
	}
}

func (d *debugRequest) ProtoMajor() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ProtoMajor", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.Itoa(hit.Request().ProtoMajor))
			return err
		},
	}
}

func (d *debugRequest) ProtoMinor() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ProtoMinor", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.Itoa(hit.Request().ProtoMinor))
			return err
		},
	}
}

func (d *debugRequest) ContentLength() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ContentLength", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.FormatInt(hit.Request().ContentLength, 10))
			return err
		},
	}
}

func (d *debugRequest) TransferEncoding() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("TransferEncoding", nil),
		Exec: func(hit *hitImpl) error {
			te := hit.Request().TransferEncoding
			if te == nil {
				te = []string{}
			}
			return d.debug.print(d.debug.out(hit), te)
		},
	}
}

func (d *debugRequest) Host() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("Host", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), hit.Request().Host)
			return err
		},
	}
}

func (d *debugRequest) Headers(header ...string) IStep {
	v, _ := getLastStringArgument(header)
	return newDebugHeaders(d.cp.Push("Headers", stringSliceToInterfaceSlice(header)), d.debug, debugHeaderRequest, v)
}

func (d *debugRequest) Trailers(trailer ...string) IStep {
	v, _ := getLastStringArgument(trailer)
	return newDebugHeaders(d.cp.Push("Trailers", stringSliceToInterfaceSlice(trailer)), d.debug, debugTrailerRequest, v)
}

func (d *debugRequest) Body() IDebugBody {
	return newDebugBody(d.cp.Push("Body", nil), d.debug, debugBodyRequest)
}
