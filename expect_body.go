package hit

import (
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type IExpectBody interface {
	IStep
	JSON(data ...interface{}) IExpectBodyJSON
	Interface(data interface{}) IStep
	Equal(data interface{}) IStep
	NotEqual(data interface{}) IStep
	Contains(data interface{}) IStep
	NotContains(data interface{}) IStep
}

type expectBody struct {
	expect    IExpect
	cleanPath clearPath
}

func newExpectBody(expect IExpect, cleanPath clearPath, params []interface{}) IExpectBody {
	body := &expectBody{
		expect:    expect,
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Equal()
		return finalExpectBody{&hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: cleanPath,
			Exec:      body.Interface(param).exec,
		}}
	}

	return body
}

func (*expectBody) exec(Hit) error {
	return xerrors.New("unsupported")
}

func (*expectBody) when() StepTime {
	return ExpectStep
}

func (body *expectBody) clearPath() clearPath {
	return body.cleanPath
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) IExpectBodyJSON {
	return newExpectBodyJSON(body, body.cleanPath.Push("JSON", data), data)
}

// Interface expects the specified interface
func (body *expectBody) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec: func(hit Hit) error {
				x(hit)
				return nil
			},
		}
	case func(e Hit) error:
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec:      x,
		}
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return &hitStep{
				Trace:     ett.Prepare(),
				When:      ExpectStep,
				ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
				Exec: func(hit Hit) error {
					internal.CallGenericFunc(f)
					return nil
				},
			}
		}
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec:      body.Equal(data).exec,
		}
	}
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("Equal", []interface{}{data}),
		Exec: func(hit Hit) error {
			if hit.Response().body.equalOnlyNativeTypes(data, true) {
				return nil
			}
			return Expect().Body().JSON().Equal("", data).exec(hit)
		},
	}
}

// NotEqual expects the body to be not equal to the specified value
// Example:
//           Expect().Body().NotEqual("Hello World")
func (body *expectBody) NotEqual(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("NotEqual", []interface{}{data}),
		Exec: func(hit Hit) error {
			if hit.Response().body.equalOnlyNativeTypes(data, false) {
				return nil
			}
			return Expect().Body().JSON().NotEqual("", data).exec(hit)
		},
	}
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("Contains", []interface{}{data}),
		Exec: func(hit Hit) error {
			if hit.Response().body.containsOnlyNativeTypes(data, true) {
				return nil
			}
			return Expect().Body().JSON().Contains("", data).exec(hit)
		},
	}
}

// NotContains expects the body to not contain the specified value
// Example:
//           Expect().Body().NotContains("Hello World")
func (body *expectBody) NotContains(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: body.cleanPath.Push("NotContains", []interface{}{data}),
		Exec: func(hit Hit) error {
			if hit.Response().body.containsOnlyNativeTypes(data, false) {
				return nil
			}
			return Expect().Body().JSON().NotContains("", data).exec(hit)
		},
	}
}

type finalExpectBody struct {
	IStep
}

func (f finalExpectBody) JSON(data ...interface{}) IExpectBodyJSON {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Interface(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Equal(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotEqual(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Contains(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotContains(interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
