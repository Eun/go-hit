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
	hit       Hit
	cleanPath CleanPath
	params    []interface{}
}

func newExpectBody(expect IExpect, hit Hit, path CleanPath, params []interface{}) IExpectBody {
	return &expectBody{
		expect:    expect,
		hit:       hit,
		cleanPath: path,
		params:    params,
	}
}

func (*expectBody) When() StepTime {
	return ExpectStep
}

// Exec contains the logic for Expect().Body(...)
func (exp *expectBody) Exec(hit Hit) error {
	param, ok := internal.GetLastArgument(exp.params)
	if !ok {
		return errors.New("invalid argument")
	}
	exp.Equal(param)
	return nil
}

func (body *expectBody) CleanPath() CleanPath {
	return body.cleanPath
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) IExpectBodyJSON {
	return newExpectBodyJSON(body, body.hit, body.cleanPath.Push("JSON", data), data)
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("Equal"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			if hit.Response().body.equalOnlyNativeTypes(data, true) {
				return
			}
			hit.Expect().Body().JSON().Equal("", data)
		},
	})
}

// NotEqual expects the body to be not equal to the specified value
// Example:
//           Expect().Body().NotEqual("Hello World")
func (body *expectBody) NotEqual(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("NotEqual"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			if hit.Response().body.equalOnlyNativeTypes(data, false) {
				return
			}
			hit.Expect().Body().JSON().NotEqual("", data)
		},
	})
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("Contains"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			if hit.Response().body.containsOnlyNativeTypes(data, true) {
				return
			}
			hit.Expect().Body().JSON().Contains("", data)
		},
	})
}

// NotContains expects the body to not contain the specified value
// Example:
//           Expect().Body().NotContains("Hello World")
func (body *expectBody) NotContains(data interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: body.cleanPath.Push("NotContains"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			if hit.Response().body.containsOnlyNativeTypes(data, false) {
				return
			}
			hit.Expect().Body().JSON().NotContains("", data)
		},
	})
}
