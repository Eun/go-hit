package hit

import (
	"io"
	"strconv"

	"github.com/Eun/go-hit/errortrace"
)

// IDebugResponse defines the debug functions that are available for the http response.
type IDebugResponse interface {
	IStep

	// Proto will print the Response's Proto
	Proto() IStep

	// ProtoMajor will print the Response's ProtoMajor
	ProtoMajor() IStep

	// ProtoMinor will print the Response's ProtoMinor
	ProtoMinor() IStep

	// ContentLength will print the Response's ContentLength
	ContentLength() IStep

	// TransferEncoding will print the Response's TransferEncoding
	TransferEncoding() IStep

	// Status will print the Response's Status
	Status() IStep

	// StatusCode will print the Response's StatusCode
	StatusCode() IStep

	// Headers prints the Response's Headers
	//
	// If you omit the header it will print all headers.
	//
	// Usage:
	//     Debug().Response().Headers()               // all headers
	//     Debug().Response().Headers("Content-Type") // only Content-Type
	Headers(header ...string) IStep

	// Trailers prints the Response's Trailers
	//
	// If you omit the trailer it will print all trailers.
	//
	// Usage:
	//     Debug().Response().Trailers()               // all trailers
	//     Debug().Response().Trailers("Content-Type") // only Content-Type
	Trailers(trailer ...string) IStep

	// Body prints the Response's Body
	//
	// The argument can be used to narrow down the print path
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Debug().Request().Body()                       // will print the whole Body
	//     Debug().Request().Body().JSON()                // will print the whole Body
	//     Debug().Request().Body().JSON().JQ(".ID")      // will print 10
	//     Debug().Request().Body().JSON().JQ(".Name")    // will print "Name"
	//     Debug().Request().Body().JSON().JQ(".Roles")   // will print ["Admin", "User"]
	//     Debug().Request().Body().JSON().JQ(".Roles.0") // will print "Admin"
	Body() IDebugBody
}

type debugResponse struct {
	cp    callPath
	debug *debug
}

func newDebugResponse(cp callPath, d *debug) *debugResponse {
	return &debugResponse{
		cp:    cp,
		debug: d,
	}
}

func (*debugResponse) trace() *errortrace.ErrorTrace {
	return nil
}

func (*debugResponse) when() StepTime {
	return BeforeExpectStep
}

func (d *debugResponse) data(hit Hit) map[string]interface{} {
	// we have to read body before trailer
	body := hit.Response().body.GetBestFittingObject()
	return map[string]interface{}{
		"Proto":            hit.Response().Proto,
		"ProtoMajor":       hit.Response().ProtoMajor,
		"ProtoMinor":       hit.Response().ProtoMinor,
		"ContentLength":    hit.Response().ContentLength,
		"TransferEncoding": hit.Response().TransferEncoding,
		"Status":           hit.Response().Status,
		"StatusCode":       hit.Response().StatusCode,
		"Headers":          d.debug.getMap(hit.Response().Header),
		"Trailers":         d.debug.getMap(hit.Response().Trailer),
		"Body":             body,
	}
}

func (d *debugResponse) exec(hit *hitImpl) error {
	return d.debug.print(d.debug.out(hit), d.data(hit))
}

func (d *debugResponse) callPath() callPath {
	return d.cp
}

func (d *debugResponse) Proto() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("Proto", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), hit.Response().Proto)
			return err
		},
	}
}

func (d *debugResponse) ProtoMajor() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ProtoMajor", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.Itoa(hit.Response().ProtoMajor))
			return err
		},
	}
}

func (d *debugResponse) ProtoMinor() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ProtoMinor", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.Itoa(hit.Response().ProtoMinor))
			return err
		},
	}
}

func (d *debugResponse) ContentLength() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("ContentLength", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.FormatInt(hit.Response().ContentLength, 10))
			return err
		},
	}
}

func (d *debugResponse) TransferEncoding() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("TransferEncoding", nil),
		Exec: func(hit *hitImpl) error {
			te := hit.Response().TransferEncoding
			if te == nil {
				te = []string{}
			}
			return d.debug.print(d.debug.out(hit), te)
		},
	}
}

func (d *debugResponse) Status() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("Status", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), hit.Response().Status)
			return err
		},
	}
}

func (d *debugResponse) StatusCode() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: d.cp.Push("StatusCode", nil),
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.debug.out(hit), strconv.Itoa(hit.Response().StatusCode))
			return err
		},
	}
}

func (d *debugResponse) Headers(header ...string) IStep {
	v, _ := getLastStringArgument(header)
	return newDebugHeaders(d.cp.Push("Headers", stringSliceToInterfaceSlice(header)), d.debug, debugHeaderResponse, v)
}

func (d *debugResponse) Trailers(trailer ...string) IStep {
	v, _ := getLastStringArgument(trailer)
	return newDebugHeaders(d.cp.Push("Trailers", stringSliceToInterfaceSlice(trailer)), d.debug, debugTrailerResponse, v)
}

func (d *debugResponse) Body() IDebugBody {
	return newDebugBody(d.cp.Push("Body", nil), d.debug, debugBodyResponse)
}
