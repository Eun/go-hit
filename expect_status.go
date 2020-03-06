package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectStatus interface {
	IStep
	Equal(statusCode int) IStep
	NotEqual(statusCode int) IStep
	OneOf(statusCodes ...int) IStep
	NotOneOf(statusCodes ...int) IStep
	GreaterThan(statusCode int) IStep
	LessThan(statusCode int) IStep
	GreaterOrEqualThan(statusCode int) IStep
	LessOrEqualThan(statusCode int) IStep
	Between(min, max int) IStep
	NotBetween(min, max int) IStep
}

type expectStatus struct {
	expect    IExpect
	cleanPath CleanPath
}

func newExpectStatus(expect IExpect, cleanPath CleanPath, params []int) IExpectStatus {
	status := &expectStatus{
		expect:    expect,
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastIntArgument(params); ok {
		return finalExpectStatus{wrap(wrapStep{
			IStep:     status.Equal(param),
			CleanPath: cleanPath,
		})}
	}

	return status
}

func (*expectStatus) Exec(Hit) error {
	return errors.New("unsupported")
}

func (*expectStatus) When() StepTime {
	return ExpectStep
}

func (status *expectStatus) CleanPath() CleanPath {
	return status.cleanPath
}

// Equal checks if the status is equal to the specified value
// Examples:
//           Expect().Status().Equal(200)
func (status *expectStatus) Equal(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("Equal", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode != statusCode {
				minitest.Errorf("Expected status code to be %d but was %d instead", statusCode, hit.Response().StatusCode)
			}
		},
	})
}

// NotEqual checks if the status is equal to the specified value
// Examples:
//           Expect().Status().NotEqual(200)
func (status *expectStatus) NotEqual(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("NotEqual", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode == statusCode {
				minitest.Errorf("Expected status code not to be %d", statusCode)
			}
		},
	})
}

// OneOf checks if the status is one of the specified values
// Examples:
//           Expect().Status().OneOf(200, 201)
func (status *expectStatus) OneOf(statusCodes ...int) IStep {
	args := make([]interface{}, len(statusCodes))
	for i := range statusCodes {
		args[i] = statusCodes[i]
	}
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("OneOf", args),
		Exec: func(hit Hit) {
			minitest.Contains(statusCodes, hit.Response().StatusCode)
		},
	})
}

// NotOneOf checks if the status is none of the specified values
// Examples:
//           Expect().Status().NotOneOf(200, 201)
func (status *expectStatus) NotOneOf(statusCodes ...int) IStep {
	args := make([]interface{}, len(statusCodes))
	for i := range statusCodes {
		args[i] = statusCodes[i]
	}
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("NotOneOf", args),
		Exec: func(hit Hit) {
			minitest.NotContains(statusCodes, hit.Response().StatusCode)
		},
	})
}

// GreaterThan checks if the status is greater than the specified value
// Examples:
//           Expect().Status().GreaterThan(400)
func (status *expectStatus) GreaterThan(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("GreaterThan", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode <= statusCode {
				minitest.Errorf("expected %d to be greater than %d", hit.Response().StatusCode, statusCode)
			}
		},
	})
}

// LessThan checks if the status is less than the specified value
// Examples:
//           Expect().Status().LessThan(400)
func (status *expectStatus) LessThan(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("LessThan", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode >= statusCode {
				minitest.Errorf("expected %d to be less than %d", hit.Response().StatusCode, statusCode)
			}
		},
	})
}

// GreaterOrEqualThan checks if the status is greater or equal than the specified value
// Examples:
//           Expect().Status().GreaterOrEqualThan(200)
func (status *expectStatus) GreaterOrEqualThan(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("GreaterOrEqualThan", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode < statusCode {
				minitest.Errorf("expected %d to be greater or equal than %d", hit.Response().StatusCode, statusCode)
			}
		},
	})
}

// LessOrEqualThan checks if the status is less or equal than the specified value
// Examples:
//           Expect().Status().LessOrEqualThan(200)
func (status *expectStatus) LessOrEqualThan(statusCode int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("LessOrEqualThan", []interface{}{statusCode}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode > statusCode {
				minitest.Errorf("expected %d to be less or equal than %d", hit.Response().StatusCode, statusCode)
			}
		},
	})
}

// Between checks if the status is between the specified value (inclusive)
// Examples:
//           Expect().Status().Between(200, 400)
func (status *expectStatus) Between(min, max int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("Between", []interface{}{min, max}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode < min || hit.Response().StatusCode > max {
				minitest.Errorf("expected %d to be between %d and %d", hit.Response().StatusCode, min, max)
			}
		},
	})
}

// NotBetween checks if the status is not between the specified value (inclusive)
// Examples:
//           Expect().Status().NotBetween(200, 400)
func (status *expectStatus) NotBetween(min, max int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: status.CleanPath().Push("NotBetween", []interface{}{min, max}),
		Exec: func(hit Hit) {
			if hit.Response().StatusCode >= min && hit.Response().StatusCode <= max {
				minitest.Errorf("expected %d not to be between %d and %d", hit.Response().StatusCode, min, max)
			}
		},
	})
}

type finalExpectStatus struct {
	IStep
}

func (f finalExpectStatus) Equal(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) NotEqual(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) OneOf(statusCodes ...int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) NotOneOf(statusCodes ...int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) GreaterThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) LessThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) GreaterOrEqualThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) LessOrEqualThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) Between(min, max int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalExpectStatus) NotBetween(min, max int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
