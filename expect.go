package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"github.com/mohae/deepcopy"
	"golang.org/x/xerrors"
)

// IExpect provides assertions on the http response
type IExpect interface {
	IStep
	// Body expects the body to be equal the specified value
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Body("Hello World")
	//     Expect().Body().Contains("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().Contains("Hello World"),
	//     )
	Body(value ...interface{}) IExpectBody

	// Interface expects the body to be equal the specified interface.
	//
	// Usage:
	//     Expect().Interface("Hello World")
	//     Expect().Interface(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Interface("Hello World"),
	//     )
	Interface(value interface{}) IStep

	// Header provides assertions to one specific response header.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Header().Contains("Content-Type")
	//     Expect().Header("Content-Type").Equal("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header().Contains("Content-Type"),
	//         Expect().Header("Content-Type").Equal("application/json"),
	//     )
	Header(headerName ...string) IExpectHeader

	// Status expects the status to be the specified code.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Status(200)
	//     Expect().Status().Equal(200)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Status().OneOf(http.StatusOk, http.StatusNoContent),
	//     )
	Status(code ...int) IExpectStatus

	// Custom can be used to expect a custom behaviour.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Custom(func(hit Hit) {
	//               if hit.Response().StatusCode != 200 {
	//                   panic("Expected 200")
	//               }
	//         }),
	//     )
	Custom(fn Callback) IStep
}

type expect struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newExpect(cleanPath clearPath, params []interface{}) IExpect {
	exp := &expect{
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return &finalExpect{
			&hitStep{
				Trace:     exp.trace,
				When:      ExpectStep,
				ClearPath: cleanPath,
				Exec:      exp.Interface(param).exec,
			},
			"only usable with Expect() not with Expect(value)",
		}
	}
	return exp
}

func (exp *expect) exec(hit Hit) error {
	return exp.trace.Format(hit.Description(), "unable to run Expect() without an argument or without a chain. Please use Expect(something) or Expect().Something")
}

func (*expect) when() StepTime {
	return ExpectStep
}

func (exp *expect) clearPath() clearPath {
	return exp.cleanPath
}

func (exp *expect) Body(value ...interface{}) IExpectBody {
	return newExpectBody(exp, exp.clearPath().Push("Body", value), value)
}

func (exp *expect) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.clearPath().Push("Custom", []interface{}{fn}),
		Exec: func(hit Hit) error {
			fn(hit)
			return nil
		},
	}
}

func (exp *expect) Header(headerName ...string) IExpectHeader {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}

	return newExpectHeader(exp, exp.clearPath().Push("Header", args), headerName...)
}

func (exp *expect) Status(code ...int) IExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newExpectStatus(exp, exp.clearPath().Push("Status", args), code)
}

func (exp *expect) Interface(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.clearPath().Push("Interface", []interface{}{value}),
		Exec:      exp.Body(value).exec,
	}
}

type finalExpect struct {
	IStep
	message string
}

func (exp *finalExpect) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(exp.message)
		},
	}
}

func (exp *finalExpect) Body(...interface{}) IExpectBody {
	return &finalExpectBody{
		exp.fail(),
		exp.message,
	}
}

func (exp *finalExpect) Interface(interface{}) IStep {
	return exp.fail()
}

func (exp *finalExpect) Custom(Callback) IStep {
	return exp.fail()
}

func (exp *finalExpect) Header(...string) IExpectHeader {
	return &finalExpectHeader{
		exp.fail(),
		exp.message,
	}
}

func (exp *finalExpect) Status(...int) IExpectStatus {
	return &finalExpectStatus{
		exp.fail(),
		exp.message,
	}
}

func makeCompareable(in, data interface{}) (interface{}, error) {
	compareData := deepcopy.Copy(data)
	err := converter.Convert(in, &compareData)
	if err != nil {
		return nil, err
	}

	return compareData, nil
}
