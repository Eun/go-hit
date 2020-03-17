package hit

import (
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
}

func newExpect(cleanPath clearPath, params []interface{}) IExpect {
	exp := &expect{
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalExpect{&hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: cleanPath,
			Exec:      exp.Interface(param).exec,
		}}
	}
	return exp
}

func (*expect) exec(Hit) error {
	return xerrors.New("unable to run Expect() without an argument or without a chain. Please use Expect(something) or Expect().Something")
}

func (*expect) when() StepTime {
	return ExpectStep
}

func (exp *expect) clearPath() clearPath {
	return exp.cleanPath
}

func (exp *expect) Body(value ...interface{}) IExpectBody {
	return newExpectBody(exp, exp.cleanPath.Push("Body", value), value)
}

func (exp *expect) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.cleanPath.Push("Custom", []interface{}{fn}),
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

	return newExpectHeader(exp, exp.cleanPath.Push("Header", args), headerName...)
}

func (exp *expect) Status(code ...int) IExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newExpectStatus(exp, exp.cleanPath.Push("Status", args), code)
}

func (exp *expect) Interface(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.cleanPath.Push("Interface", []interface{}{value}),
		Exec:      exp.Body(value).exec,
	}
}

type finalExpect struct {
	IStep
}

func (finalExpect) Body(...interface{}) IExpectBody {
	panic("only usable with Expect() not with Expect(value)")
}

func (finalExpect) Interface(interface{}) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (finalExpect) Custom(Callback) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (finalExpect) Header(...string) IExpectHeader {
	panic("only usable with Expect() not with Expect(value)")
}

func (finalExpect) Status(...int) IExpectStatus {
	panic("only usable with Expect() not with Expect(value)")
}

func makeCompareable(in, data interface{}) (interface{}, error) {
	compareData := deepcopy.Copy(data)
	err := converter.Convert(in, &compareData)
	if err != nil {
		return nil, err
	}

	return compareData, nil
}
