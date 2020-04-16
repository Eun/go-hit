package hit

import (
	"io"
	"strconv"
)

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
	//
	// The argument can be used to narrow down the print path
	TransferEncoding(expression ...string) IStep

	// Status will print the Response's Status
	Status() IStep

	// StatusCode will print the Response's StatusCode
	StatusCode() IStep

	// Header prints the Response's Header
	//
	// If you omit the argument it will print all headers.
	//
	// Usage:
	//     Debug().Response().Header()               // all headers
	//     Debug().Response().Header("Content-Type") // only Content-Type
	Header(expression ...string) IStep

	// Trailer prints the Response's Trailer
	//
	// If you omit the argument it will print all trailers.
	//
	// Usage:
	//     Debug().Response().Trailer()               // all trailers
	//     Debug().Response().Trailer("Content-Type") // only Content-Type
	Trailer(expression ...string) IStep

	// Body prints the Response's Body
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

type debugResponse struct {
	debug *debug
}

func newDebugResponse(d *debug) *debugResponse {
	return &debugResponse{d}
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
		"Header":           d.debug.getMap(hit.Response().Header),
		"Trailer":          d.debug.getMap(hit.Response().Trailer),
		"Body":             body,
	}
}

func (d *debugResponse) exec(hit Hit) error {
	return d.debug.print(hit, d.data(hit))
}

func (*debugResponse) clearPath() clearPath {
	return nil // not clearable
}

func (d *debugResponse) Proto() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), hit.Response().Proto)
			return err
		},
	}
}

func (d *debugResponse) ProtoMajor() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.Itoa(hit.Response().ProtoMajor))
			return err
		},
	}
}

func (d *debugResponse) ProtoMinor() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.Itoa(hit.Response().ProtoMinor))
			return err
		},
	}
}

func (d *debugResponse) ContentLength() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.FormatInt(hit.Response().ContentLength, 10))
			return err
		},
	}
}

func (d *debugResponse) TransferEncoding(expression ...string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			te := hit.Response().TransferEncoding
			if te == nil {
				te = []string{}
			}
			return d.debug.printWithExpression(hit, te, expression)
		},
	}
}

func (d *debugResponse) Status() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), hit.Response().Status)
			return err
		},
	}
}

func (d *debugResponse) StatusCode() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), strconv.Itoa(hit.Response().StatusCode))
			return err
		},
	}
}

func (d *debugResponse) Header(expression ...string) IStep {
	return newDebugHeader(d.debug, debugHeaderResponse, expression)
}

func (d *debugResponse) Trailer(expression ...string) IStep {
	return newDebugHeader(d.debug, debugTrailerResponse, expression)
}

func (d *debugResponse) Body(expression ...string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return d.debug.printWithExpression(hit, hit.Response().Body().GetBestFittingObject(), expression)
		},
	}
}
