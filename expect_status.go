package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
	"golang.org/x/xerrors"
)

// IExpectSpecificHeader provides assertions on the http response code
type IExpectStatus interface {
	IStep
	// Equal expects the status to be equal to the specified code.
	//
	// Usage:
	//     Expect().Status().Equal(http.StatusOK)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().Equal(http.StatusOK),
	//     )
	Equal(statusCode int) IStep

	// NotEqual expects the status to be not equal to the specified code.
	//
	// Usage:
	//     Expect().Status().NotEqual(http.StatusOK)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotEqual(http.StatusOK),
	//     )
	NotEqual(statusCode int) IStep

	// OneOf expects the status to be equal to one of the specified codes.
	//
	// Usage:
	//     Expect().Status().OneOf(http.StatusOK, http.StatusNoContent)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().OneOf(http.StatusOK, http.StatusNoContent),
	//     )
	OneOf(statusCodes ...int) IStep

	// NotOneOf expects the status to be not equal to one of the specified codes.
	//
	// Usage:
	//     Expect().Status().NotOneOf(http.StatusUnauthorized, http.StatusForbidden)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotOneOf(http.StatusUnauthorized, http.StatusForbidden),
	//     )
	NotOneOf(statusCodes ...int) IStep

	// GreaterThan expects the status to be not greater than the specified code.
	//
	// Usage:
	//     Expect().Status().GreaterThan(http.StatusContinue)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterThan(http.StatusContinue),
	//     )
	GreaterThan(statusCode int) IStep

	// LessThan expects the status to be less than the specified code.
	//
	// Usage:
	//     Expect().Status().LessThan(http.StatusInternalServerError)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().LessThan(http.StatusInternalServerError),
	//     )
	LessThan(statusCode int) IStep

	// GreaterOrEqualThan expects the status to be greater or equal than the specified code.
	//
	// Usage:
	//     Expect().Status().GreaterOrEqualThan(http.StatusOK)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterOrEqualThan(http.StatusOK),
	//     )
	GreaterOrEqualThan(statusCode int) IStep

	// LessOrEqualThan expects the status to be less or equal than the specified code.
	//
	// Usage:
	//     Expect().Status().LessOrEqualThan(http.StatusOK)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().LessOrEqualThan(http.StatusOK),
	//     )
	LessOrEqualThan(statusCode int) IStep

	// Between expects the status to be between the specified min and max value. (inclusive)
	//
	// Usage:
	//     Expect().Status().Between(http.StatusOK, http.StatusAccepted)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().Between(http.StatusOK, http.StatusAccepted),
	//     )
	Between(min, max int) IStep

	// NotBetween expects the status to be not between the specified min and max value. (inclusive)
	//
	// Usage:
	//     Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted)
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted),
	//     )
	NotBetween(min, max int) IStep
}

type expectStatus struct {
	expect    IExpect
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newExpectStatus(expect IExpect, cleanPath clearPath, params []int) IExpectStatus {
	status := &expectStatus{
		expect:    expect,
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}

	if param, ok := internal.GetLastIntArgument(params); ok {
		return finalExpectStatus{&hitStep{
			Trace:     status.trace,
			When:      ExpectStep,
			ClearPath: cleanPath,
			Exec:      status.Equal(param).exec,
		}}
	}

	return status
}

func (status *expectStatus) exec(hit Hit) error {
	return status.trace.Format(hit.Description(), "unable to run Expect().Status() without an argument or without a chain. Please use Expect().Status(something) or Expect().Status().Something")
}

func (*expectStatus) when() StepTime {
	return ExpectStep
}

func (status *expectStatus) clearPath() clearPath {
	return status.cleanPath
}

func (status *expectStatus) Equal(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("Equal", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode != statusCode {
				minitest.Errorf("Expected status code to be %d but was %d instead", statusCode, hit.Response().StatusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) NotEqual(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("NotEqual", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode == statusCode {
				minitest.Errorf("Expected status code not to be %d", statusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) OneOf(statusCodes ...int) IStep {
	args := make([]interface{}, len(statusCodes))
	for i := range statusCodes {
		args[i] = statusCodes[i]
	}
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("OneOf", args),
		Exec: func(hit Hit) error {
			minitest.Contains(statusCodes, hit.Response().StatusCode)
			return nil
		},
	}
}

func (status *expectStatus) NotOneOf(statusCodes ...int) IStep {
	args := make([]interface{}, len(statusCodes))
	for i := range statusCodes {
		args[i] = statusCodes[i]
	}
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("NotOneOf", args),
		Exec: func(hit Hit) error {
			minitest.NotContains(statusCodes, hit.Response().StatusCode)
			return nil
		},
	}
}

func (status *expectStatus) GreaterThan(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("GreaterThan", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode <= statusCode {
				minitest.Errorf("expected %d to be greater than %d", hit.Response().StatusCode, statusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) LessThan(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("LessThan", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode >= statusCode {
				minitest.Errorf("expected %d to be less than %d", hit.Response().StatusCode, statusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) GreaterOrEqualThan(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("GreaterOrEqualThan", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode < statusCode {
				minitest.Errorf("expected %d to be greater or equal than %d", hit.Response().StatusCode, statusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) LessOrEqualThan(statusCode int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("LessOrEqualThan", []interface{}{statusCode}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode > statusCode {
				minitest.Errorf("expected %d to be less or equal than %d", hit.Response().StatusCode, statusCode)
			}
			return nil
		},
	}
}

func (status *expectStatus) Between(min, max int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("Between", []interface{}{min, max}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode < min || hit.Response().StatusCode > max {
				minitest.Errorf("expected %d to be between %d and %d", hit.Response().StatusCode, min, max)
			}
			return nil
		},
	}
}

func (status *expectStatus) NotBetween(min, max int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: status.clearPath().Push("NotBetween", []interface{}{min, max}),
		Exec: func(hit Hit) error {
			if hit.Response().StatusCode >= min && hit.Response().StatusCode <= max {
				minitest.Errorf("expected %d not to be between %d and %d", hit.Response().StatusCode, min, max)
			}
			return nil
		},
	}
}

type finalExpectStatus struct {
	IStep
}

func (finalExpectStatus) Equal(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) NotEqual(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) OneOf(...int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) NotOneOf(...int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) GreaterThan(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) LessThan(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) GreaterOrEqualThan(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) LessOrEqualThan(int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) Between(int, int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
func (finalExpectStatus) NotBetween(int, int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Expect().Status() not with Expect().Status(value)")
		},
	}
}
