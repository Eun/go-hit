package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/httpbody"
)

// IDebugBodyJSON defines the debug functions that are available for the http request/response body in JSON format.
type IDebugBodyJSON interface {
	IStep
	// JQ runs an jq expression on the JSON body and prints the result
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Debug().Response().Body().JSON().JQ(".Name") // print "Joe"
	JQ(expression ...string) IStep
}

type debugBodyJSON struct {
	cp    callPath
	debug *debug
	mode  debugBodyMode
}

func newDebugBodyJSON(cp callPath, debug *debug, mode debugBodyMode) IDebugBodyJSON {
	return &debugBodyJSON{
		cp:    cp,
		debug: debug,
		mode:  mode,
	}
}

func (*debugBodyJSON) trace() *errortrace.ErrorTrace {
	return nil
}

func (s *debugBodyJSON) body(hit Hit) *httpbody.HTTPBody {
	if s.mode == debugBodyRequest {
		return hit.Request().Body()
	}
	return hit.Response().Body()
}

func (s *debugBodyJSON) when() StepTime {
	return BeforeExpectStep
}

func (s *debugBodyJSON) callPath() callPath {
	return s.cp
}

func (s *debugBodyJSON) exec(hit *hitImpl) error {
	var container interface{}
	s.body(hit).JSON().MustDecode(&container)
	return s.debug.print(s.debug.out(hit), container)
}

func (s *debugBodyJSON) JQ(expression ...string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			var container interface{}
			s.body(hit).JSON().MustJQ(&container, expression...)
			return s.debug.print(s.debug.out(hit), container)
		},
	}
}
