package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
)

type IExpectBody interface {
	IStep
	JSON(data ...interface{}) IExpectBodyJSON
	Equal(data interface{}) IStep
	NotEqual(data interface{}) IStep
	Contains(data interface{}) IStep
	NotContains(data interface{}) IStep
}

type expectBody struct {
	expect    IExpect
	cleanPath CleanPath
}

func newExpectBody(expect IExpect, cleanPath CleanPath, params []interface{}) IExpectBody {
	body := &expectBody{
		expect:    expect,
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Equal()
		return finalExpectBody{wrap(wrapStep{
			IStep:     body.Equal(param),
			CleanPath: cleanPath,
		})}
	}

	return body
}

func (*expectBody) Exec(Hit) error {
	return errors.New("unsupported")
}

func (*expectBody) When() StepTime {
	return ExpectStep
}

func (body *expectBody) CleanPath() CleanPath {
	return body.cleanPath
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) IExpectBodyJSON {
	return newExpectBodyJSON(body, body.cleanPath.Push("JSON", data), data)
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("Equal", []interface{}{data}),
		Exec: func(hit Hit) {
			if hit.Response().body.equalOnlyNativeTypes(data, true) {
				return
			}
			hit.RunSteps(Expect().Body().JSON().Equal("", data))
		},
	})
}

// NotEqual expects the body to be not equal to the specified value
// Example:
//           Expect().Body().NotEqual("Hello World")
func (body *expectBody) NotEqual(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("NotEqual", []interface{}{data}),
		Exec: func(hit Hit) {
			if hit.Response().body.equalOnlyNativeTypes(data, false) {
				return
			}
			panic("TODO")
			//hit.Expect().Body().JSON().NotEqual("", data)
		},
	})
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("Contains", []interface{}{data}),
		Exec: func(hit Hit) {
			if hit.Response().body.containsOnlyNativeTypes(data, true) {
				return
			}
			panic("TODO")
			//hit.Expect().Body().JSON().Contains("", data)
		},
	})
}

// NotContains expects the body to not contain the specified value
// Example:
//           Expect().Body().NotContains("Hello World")
func (body *expectBody) NotContains(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("NotContains", []interface{}{data}),
		Exec: func(hit Hit) {
			if hit.Response().body.containsOnlyNativeTypes(data, false) {
				return
			}
			panic("TODO")
			//hit.Expect().Body().JSON().NotContains("", data)
		},
	})
}

type finalExpectBody struct {
	IStep
}

func (f finalExpectBody) JSON(data ...interface{}) IExpectBodyJSON {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Equal(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotEqual(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Contains(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotContains(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
