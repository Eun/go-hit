package hit

import (
	"io"
	"strconv"
)

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
	//     Debug().Request().URL("Path") // print only the URL.Path
	URL(expression ...string) IStep

	// Proto prints the Request's Proto
	Proto() IStep

	// ProtoMajor prints the Request's ProtoMajor
	ProtoMajor() IStep

	// ProtoMinor prints the Request's ProtoMinor
	ProtoMinor() IStep

	// ContentLength prints the Request's ContentLength
	ContentLength() IStep

	// TransferEncoding prints the Request's TransferEncoding
	//
	// The argument can be used to narrow down the print path
	TransferEncoding(expression ...string) IStep

	// Host prints the Request's Host
	Host() IStep

	// Header prints the Request's Header
	//
	// If you omit the argument it will print all headers.
	//
	// Usage:
	//     Debug().Request().Header()               // all headers
	//     Debug().Request().Header("Content-Type") // only Content-Type
	Header(expression ...string) IStep

	// Trailer prints the Request's Trailer
	//
	// If you omit the argument it will print all trailers.
	//
	// Usage:
	//     Debug().Request().Trailer()               // all trailers
	//     Debug().Request().Trailer("Content-Type") // only Content-Type
	Trailer(expression ...string) IStep

	// Body prints the Request's Body
	//
	// The argument can be used to narrow down the print path
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Debug().Request().Body()          // will print the whole Body
	//     Debug().Request().Body("ID")      // will print 10
	//     Debug().Request().Body("Name")    // will print "Name"
	//     Debug().Request().Body("Roles")   // will print ["Admin", "User"]
	//     Debug().Request().Body("Roles.0") // will print "Admin"
	Body(expression ...string) IStep
}

type debugRequest struct {
	debug *debug
}

func newDebugRequest() *debugRequest {
	return &debugRequest{}
}

func (*debugRequest) when() StepTime {
	return BeforeExpectStep
}

func (d *debugRequest) data(hit Hit) map[string]interface{} {
	return map[string]interface{}{
		"Method":           hit.Request().Method,
		"URL":              hit.Request().URL,
		"Proto":            hit.Request().Proto,
		"ProtoMajor":       hit.Request().ProtoMajor,
		"ProtoMinor":       hit.Request().ProtoMinor,
		"ContentLength":    hit.Request().ContentLength,
		"TransferEncoding": hit.Request().TransferEncoding,
		"Host":             hit.Request().Host,
		"Header":           d.debug.getMap(hit.Request().Header),
		"Trailer":          d.debug.getMap(hit.Request().Trailer),
		"Body":             d.debug.getBody(hit.Request().Body()),
	}
}

func (d *debugRequest) exec(hit Hit) error {
	return d.debug.print(hit, d.data(hit))
}

func (*debugRequest) clearPath() clearPath {
	return nil // not clearable
}

func (d *debugRequest) Method() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), hit.Request().Method)
			return err
		},
	}
}

func (d *debugRequest) URL(expression ...string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return d.debug.printWithExpression(hit, hit.Request().URL, expression)
		},
	}
}

func (d *debugRequest) Proto() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), hit.Request().Proto)
			return err
		},
	}
}

func (d *debugRequest) ProtoMajor() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.Itoa(hit.Request().ProtoMajor))
			return err
		},
	}
}

func (d *debugRequest) ProtoMinor() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.Itoa(hit.Request().ProtoMinor))
			return err
		},
	}
}

func (d *debugRequest) ContentLength() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.FormatInt(hit.Request().ContentLength, 10))
			return err
		},
	}
}

func (d *debugRequest) TransferEncoding(expression ...string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return d.debug.printWithExpression(hit, hit.Request().TransferEncoding, expression)
		},
	}
}

func (d *debugRequest) Host() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), hit.Request().Host)
			return err
		},
	}
}

func (d *debugRequest) Header(expression ...string) IStep {
	return newDebugHeader(d.debug, debugHeaderRequest, expression)
}

func (d *debugRequest) Trailer(expression ...string) IStep {
	return newDebugHeader(d.debug, debugTrailerRequest, expression)
}

func (d *debugRequest) Body(expression ...string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return d.debug.printWithExpression(hit, d.debug.getBody(hit.Request().Body()), expression)
		},
	}
}
