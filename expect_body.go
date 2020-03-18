package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
)

// IExpectBody provides assertions on the http response body
type IExpectBody interface {
	IStep

	// Interface expects the body to be equal the specified interface.
	//
	// Usage:
	//     Expect().Body().Interface("Hello World")
	//     Expect().Body().Interface(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().Interface("Hello World"),
	//     )
	Interface(data interface{}) IStep

	// JSON expects the body to be equal the specified value.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//           Expect().Body().JSON([]int{1, 2, 3})
	//           Expect().Body().JSON().Contains(1)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//     )
	JSON(value ...interface{}) IExpectBodyJSON

	// Equal expects the body to be equal to the specified value
	//
	// Usage:
	//     Expect().Body().Equal("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().Equal("Hello World"),
	//     )
	Equal(value interface{}) IStep

	// NotEqual expects the body to be not equal to the specified value
	//
	// Usage:
	//     Expect().Body().NotEqual("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().NotEqual("Hello World"),
	//     )
	NotEqual(value interface{}) IStep

	// Contains expects the body to contain the specified value
	//
	// Usage:
	//     Expect().Body().Contains("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().Contains("Hello World"),
	//     )
	Contains(value interface{}) IStep

	// NotContains expects the body to not contain the specified value
	//
	// Usage:
	//     Expect().Body().NotContains("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().NotContains("Hello World"),
	//     )
	NotContains(value interface{}) IStep
}

type expectBody struct {
	expect    IExpect
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newExpectBody(expect IExpect, cleanPath clearPath, params []interface{}) IExpectBody {
	body := &expectBody{
		expect:    expect,
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Equal()
		return finalExpectBody{&hitStep{
			Trace:     body.trace,
			When:      ExpectStep,
			ClearPath: cleanPath,
			Exec:      body.Interface(param).exec,
		}}
	}

	return body
}

func (body *expectBody) exec(hit Hit) error {
	return body.trace.Format(hit.Description(), "unable to run Expect().Body() without an argument or without a chain. Please use Expect().Body(something) or Expect().Body().Something")
}

func (*expectBody) when() StepTime {
	return ExpectStep
}

func (body *expectBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *expectBody) JSON(value ...interface{}) IExpectBodyJSON {
	return newExpectBodyJSON(body, body.cleanPath.Push("JSON", value), value)
}

func (body *expectBody) Interface(value interface{}) IStep {
	switch x := value.(type) {
	case func(e Hit):
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec: func(hit Hit) error {
				x(hit)
				return nil
			},
		}
	case func(e Hit) error:
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec:      x,
		}
	default:
		if f := internal.GetGenericFunc(value); f.IsValid() {
			return &hitStep{
				Trace:     ett.Prepare(),
				When:      ExpectStep,
				ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
				Exec: func(hit Hit) error {
					internal.CallGenericFunc(f)
					return nil
				},
			}
		}
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec:      body.Equal(value).exec,
		}
	}
}

func (body *expectBody) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit Hit) error {
			if hit.Response().body.equalOnlyNativeTypes(value, true) {
				return nil
			}
			return Expect().Body().JSON().Equal("", value).exec(hit)
		},
	}
}

func (body *expectBody) NotEqual(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit Hit) error {
			if hit.Response().body.equalOnlyNativeTypes(value, false) {
				return nil
			}
			return Expect().Body().JSON().NotEqual("", value).exec(hit)
		},
	}
}

func (body *expectBody) Contains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("Contains", []interface{}{value}),
		Exec: func(hit Hit) error {
			if hit.Response().body.containsOnlyNativeTypes(value, true) {
				return nil
			}
			return Expect().Body().JSON().Contains("", value).exec(hit)
		},
	}
}

func (body *expectBody) NotContains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit Hit) error {
			if hit.Response().body.containsOnlyNativeTypes(value, false) {
				return nil
			}
			return Expect().Body().JSON().NotContains("", value).exec(hit)
		},
	}
}

type finalExpectBody struct {
	IStep
}

func (finalExpectBody) JSON(...interface{}) IExpectBodyJSON {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (finalExpectBody) Interface(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (finalExpectBody) Equal(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (finalExpectBody) NotEqual(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (finalExpectBody) Contains(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (finalExpectBody) NotContains(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
