package hit

import (
	"github.com/Eun/go-hit/errortrace"
)

type expectStatus struct {
	*defaultExpect
}

func newExpectStatus(expect *defaultExpect) *expectStatus {
	return &expectStatus{
		defaultExpect: expect,
	}
}

// Equal checks if the status is equal to the specified value
// Examples:
//           Expect().Status().Equal(200)
func (status *expectStatus) Equal(statusCode int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode != statusCode {
			et.Errorf("Expected status code to be %d but was %d instead", statusCode, hit.Response().StatusCode)
		}
	})
}

// OneOf checks if the status is one of the specified values
// Examples:
//           Expect().Status().OneOf(200, 201)
func (status *expectStatus) OneOf(statusCodes ...int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		et.Contains(statusCodes, hit.Response().StatusCode)
	})
}

// GreaterThan checks if the status is greater than the specified value
// Examples:
//           Expect().Status().GreaterThan(400)
func (status *expectStatus) GreaterThan(statusCode int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode <= statusCode {
			et.Errorf("expected %d to be greater than %d", hit.Response().StatusCode, statusCode)
		}
	})
}

// LessThan checks if the status is less than the specified value
// Examples:
//           Expect().Status().LessThan(400)
func (status *expectStatus) LessThan(statusCode int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode >= statusCode {
			et.Errorf("expected %d to be less than %d", hit.Response().StatusCode, statusCode)
		}
	})
}

// GreaterOrEqualThan checks if the status is greater or equal than the specified value
// Examples:
//           Expect().Status().GreaterOrEqualThan(200)
func (status *expectStatus) GreaterOrEqualThan(statusCode int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode < statusCode {
			et.Errorf("expected %d to be greater or equal than %d", hit.Response().StatusCode, statusCode)
		}
	})
}

// LessOrEqualThan checks if the status is less or equal than the specified value
// Examples:
//           Expect().Status().LessOrEqualThan(200)
func (status *expectStatus) LessOrEqualThan(statusCode int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode > statusCode {
			et.Errorf("expected %d to be less or equal than %d", hit.Response().StatusCode, statusCode)
		}
	})
}

// Between checks if the status is between the specified value (inclusive)
// Examples:
//           Expect().Status().Between(200, 400)
func (status *expectStatus) Between(min, max int) IStep {
	et := errortrace.Prepare()
	return status.Custom(func(hit Hit) {
		if hit.Response().StatusCode < min || hit.Response().StatusCode > max {
			et.Errorf("expected %d to be between %d and %d", hit.Response().StatusCode, min, max)
		}
	})
}
