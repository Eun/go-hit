package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

type IClearExpectStatus interface {
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

type clearExpectStatus struct {
	expect    IClearExpect
	cleanPath CleanPath
}

func newClearExpectStatus(expect IClearExpect, cleanPath CleanPath, params []int) IClearExpectStatus {
	status := &clearExpectStatus{
		expect:    expect,
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastIntArgument(params); ok {
		return finalClearExpectStatus{status.Equal(param)}
	}

	return status
}

func (*clearExpectStatus) Exec(Hit) error {
	return errors.New("unsupported")
}

func (*clearExpectStatus) When() StepTime {
	return CleanStep
}

func (status *clearExpectStatus) CleanPath() CleanPath {
	return status.cleanPath
}

// Equal checks if the status is equal to the specified value
// Examples:
//           Expect().Status().Equal(200)
func (status *clearExpectStatus) Equal(statusCode int) IStep {
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
func (status *clearExpectStatus) NotEqual(statusCode int) IStep {
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
func (status *clearExpectStatus) OneOf(statusCodes ...int) IStep {
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
func (status *clearExpectStatus) NotOneOf(statusCodes ...int) IStep {
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
func (status *clearExpectStatus) GreaterThan(statusCode int) IStep {
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
func (status *clearExpectStatus) LessThan(statusCode int) IStep {
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
func (status *clearExpectStatus) GreaterOrEqualThan(statusCode int) IStep {
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
func (status *clearExpectStatus) LessOrEqualThan(statusCode int) IStep {
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
func (status *clearExpectStatus) Between(min, max int) IStep {
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
func (status *clearExpectStatus) NotBetween(min, max int) IStep {
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

type finalClearExpectStatus struct {
	IStep
}

func (f finalClearExpectStatus) Equal(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) NotEqual(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) OneOf(statusCodes ...int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) NotOneOf(statusCodes ...int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) GreaterThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) LessThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) GreaterOrEqualThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) LessOrEqualThan(statusCode int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) Between(min, max int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
func (f finalClearExpectStatus) NotBetween(min, max int) IStep {
	panic("only usable with Expect().Status() not with Expect().Status(value)")
}
