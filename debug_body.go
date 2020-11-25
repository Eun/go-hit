package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/httpbody"
)

// IDebugBody defines the debug functions that are available for the http request/response body.
type IDebugBody interface {
	IStep
	// Bytes prints the body contents as a byte slice
	Bytes() IStep

	// Bool prints the body contents as a bool type
	Bool() IStep

	// Float32 prints the body contents as a float32 type
	Float32() IStep

	// Float64 prints the body contents as a float64 type
	Float64() IStep

	// FormValues prints the body contents as FormValues
	FormValues() IStep

	// Int prints the body contents as a int type
	Int() IStep

	// Int8 prints the body contents as a int8 type
	Int8() IStep

	// Int16 prints the body contents as a int16 type
	Int16() IStep

	// Int32 prints the body contents as a int32 type
	Int32() IStep

	// Int64 prints the body contents as a int64 type
	Int64() IStep

	// JSON prints the body contents in the JSON format
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Debug().Request().Body().JSON()              // print the whole body
	//     Debug().Response().Body().JSON().JQ(".Name") // print "Joe"
	JSON() IDebugBodyJSON

	// String prints the body contents as a string type
	String() IStep

	// Uint prints the body contents as a uint type
	Uint() IStep

	// Uint8 prints the body contents as a uint8 type
	Uint8() IStep

	// Uint16 prints the body contents as a uint16 type
	Uint16() IStep

	// Uint32 prints the body contents as a uint32 type
	Uint32() IStep

	// Uint64 prints the body contents as a uint64 type
	Uint64() IStep
}

type debugBodyMode int

const (
	debugBodyRequest debugBodyMode = iota
	debugBodyResponse
)

type debugBody struct {
	cp    callPath
	debug *debug
	mode  debugBodyMode
}

func newDebugBody(cp callPath, debug *debug, mode debugBodyMode) IDebugBody {
	return &debugBody{
		cp:    cp,
		debug: debug,
		mode:  mode,
	}
}

func (s *debugBody) body(hit Hit) *httpbody.HTTPBody {
	if s.mode == debugBodyRequest {
		return hit.Request().Body()
	}
	return hit.Response().Body()
}

func (*debugBody) trace() *errortrace.ErrorTrace {
	return nil
}

func (*debugBody) when() StepTime {
	return BeforeExpectStep
}

func (s *debugBody) callPath() callPath {
	return s.cp
}

func (s *debugBody) exec(hit *hitImpl) error {
	return s.debug.print(s.debug.out(hit), s.body(hit).GetBestFittingObject())
}

func (s *debugBody) Bool() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustBool())
		},
	}
}

func (s *debugBody) Bytes() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustBytes())
		},
	}
}

func (s *debugBody) Float32() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustFloat32())
		},
	}
}

func (s *debugBody) Float64() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustFloat64())
		},
	}
}

func (s *debugBody) FormValues() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustFormValues())
		},
	}
}

func (s *debugBody) Int() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustInt())
		},
	}
}

func (s *debugBody) Int8() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustInt8())
		},
	}
}

func (s *debugBody) Int16() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustInt16())
		},
	}
}

func (s *debugBody) Int32() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustInt32())
		},
	}
}

func (s *debugBody) Int64() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustInt64())
		},
	}
}

func (s *debugBody) JSON() IDebugBodyJSON {
	return newDebugBodyJSON(s.callPath().Push("JSON", nil), s.debug, s.mode)
}

func (s *debugBody) String() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustString())
		},
	}
}

func (s *debugBody) Uint() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustUint())
		},
	}
}

func (s *debugBody) Uint8() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustUint8())
		},
	}
}

func (s *debugBody) Uint16() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustUint16())
		},
	}
}

func (s *debugBody) Uint32() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustUint32())
		},
	}
}

func (s *debugBody) Uint64() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return s.debug.print(s.debug.out(hit), s.body(hit).MustUint64())
		},
	}
}
