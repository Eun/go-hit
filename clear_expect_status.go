package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

// IClearExpectStatus provides a clear functionality to remove previous steps from running in the Expect().Status() scope
type IClearExpectStatus interface {
	IStep
	// Equal removes all previous Expect().Status().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().Equal() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().Equal()              // will remove all Expect().Status().Equal() steps
	//     Clear().Expect().Status().Equal(http.StatusOK) // will remove all Expect().Status().Equal(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().Equal(http.StatusNotFound),
	//         Clear().Expect().Status().Equal(),
	//         Expect().Status().Equal(http.StatusOK),
	//     )
	Equal(code ...int) IStep

	// NotEqual removes all previous Expect().Status().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotEqual() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().NotEqual()              // will remove all Expect().Status().NotEqual() steps
	//     Clear().Expect().Status().NotEqual(http.StatusOK) // will remove all Expect().Status().NotEqual(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotEqual(http.StatusOK),
	//         Clear().Expect().Status().NotEqual(),
	//         Expect().Status().NotEqual(http.StatusNotFound),
	//     )
	NotEqual(code ...int) IStep

	// OneOf removes all previous Expect().Status().OneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().OneOf() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().OneOf()              // will remove all Expect().Status().OneOf() steps
	//     Clear().Expect().Status().OneOf(http.StatusOK) // will remove all Expect().Status().OneOf(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().OneOf(http.StatusOK, http.StatusNoContent),
	//         Clear().Expect().Status().OneOf(),
	//         Expect().Status().OneOf(http.StatusOK),
	//     )
	OneOf(code ...int) IStep

	// NotOneOf removes all previous Expect().Status().NotOneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotOneOf() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().NotOneOf()              // will remove all Expect().Status().NotOneOf() steps
	//     Clear().Expect().Status().NotOneOf(http.StatusOK) // will remove all Expect().Status().NotOneOf(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotOneOf(http.StatusOK, http.StatusNoContent),
	//         Clear().Expect().Status().NotOneOf(),
	//         Expect().Status().NotOneOf(http.StatusNotFound),
	//     )
	NotOneOf(code ...int) IStep

	// GreaterThan removes all previous Expect().Status().GreaterThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterThan() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().GreaterThan()              // will remove all Expect().Status().GreaterThan() steps
	//     Clear().Expect().Status().GreaterThan(http.StatusOK) // will remove all Expect().Status().GreaterThan(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterThan(http.StatusContinue),
	//         Clear().Expect().Status().GreaterThan(),
	//         Expect().Status().GreaterThan(http.StatusOK),
	//     )
	GreaterThan(code ...int) IStep

	// LessThan removes all previous Expect().Status().LessThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().LessThan() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().LessThan()              // will remove all Expect().Status().LessThan() steps
	//     Clear().Expect().Status().LessThan(http.StatusOK) // will remove all Expect().Status().LessThan(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().LessThan(http.StatusBadRequest),
	//         Clear().Expect().Status().LessThan(),
	//         Expect().Status().LessThan(http.StatusInternalServerError),
	//     )
	LessThan(code ...int) IStep

	// GreaterOrEqualThan removes all previous Expect().Status().GreaterOrEqualThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterOrEqualThan() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().GreaterOrEqualThan()              // will remove all Expect().Status().GreaterOrEqualThan() steps
	//     Clear().Expect().Status().GreaterOrEqualThan(http.StatusOK) // will remove all Expect().Status().GreaterOrEqualThan(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterOrEqualThan(http.StatusBadRequest),
	//         Clear().Expect().Status().GreaterOrEqualThan(),
	//         Expect().Status().GreaterOrEqualThan(http.StatusOK),
	//     )
	GreaterOrEqualThan(code ...int) IStep

	// GreaterOrEqualThan removes all previous Expect().Status().GreaterOrEqualThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterOrEqualThan() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().GreaterOrEqualThan()              // will remove all Expect().Status().GreaterOrEqualThan() steps
	//     Clear().Expect().Status().GreaterOrEqualThan(http.StatusOK) // will remove all Expect().Status().GreaterOrEqualThan(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterOrEqualThan(http.StatusBadRequest),
	//         Clear().Expect().Status().GreaterOrEqualThan(),
	//         Expect().Status().GreaterOrEqualThan(http.StatusOK),
	//     )
	LessOrEqualThan(code ...int) IStep

	// Between removes all previous Expect().Status().Between() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().Between() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().Between()              // will remove all Expect().Status().Between() steps
	//     Clear().Expect().Status().Between(http.StatusOK) // will remove all Expect().Status().Between(http.StatusOK) steps
	//     Clear().Expect().Status().Between(http.StatusOK, http.StatusAccepted) // will remove all Expect().Status().Between(http.StatusOK, http.StatusAccepted) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().Between(http.StatusOK, http.StatusAccepted),
	//         Clear().Expect().Status().Between(),
	//         Expect().Status().Between(http.StatusBadRequest, http.StatusUnavailableForLegalReasons),
	//     )
	Between(code ...int) IStep

	// NotBetween removes all previous Expect().Status().NotBetween() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotBetween() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status().NotBetween()              // will remove all Expect().Status().NotBetween() steps
	//     Clear().Expect().Status().NotBetween(http.StatusOK) // will remove all Expect().Status().NotBetween(http.StatusOK) steps
	//     Clear().Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted) // will remove all Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted),
	//         Clear().Expect().Status().NotBetween(),
	//         Expect().Status().NotBetween(http.StatusBadRequest, http.StatusUnavailableForLegalReasons),
	//     )
	NotBetween(code ...int) IStep
}

type clearExpectStatus struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearExpectStatus(cleanPath clearPath, params []int) IClearExpectStatus {
	if _, ok := internal.GetLastIntArgument(params); ok {
		// this runs if we called Clear().Expect().Status(something)
		return &finalClearExpectStatus{
			removeStep(cleanPath),
			"only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)",
		}
	}
	return &clearExpectStatus{
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}
}

func (status *clearExpectStatus) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Status()
	if err := removeSteps(hit, status.clearPath()); err != nil {
		return status.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (*clearExpectStatus) when() StepTime {
	return CleanStep
}

func (status *clearExpectStatus) clearPath() clearPath {
	return status.cleanPath
}

func (status *clearExpectStatus) Equal(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("Equal", args))
}

func (status *clearExpectStatus) NotEqual(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("NotEqual", args))
}

func (status *clearExpectStatus) OneOf(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("OneOf", args))
}

func (status *clearExpectStatus) NotOneOf(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("NotOneOf", args))
}

func (status *clearExpectStatus) GreaterThan(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("GreaterThan", args))
}

func (status *clearExpectStatus) LessThan(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("LessThan", args))
}

func (status *clearExpectStatus) GreaterOrEqualThan(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("GreaterOrEqualThan", args))
}

func (status *clearExpectStatus) LessOrEqualThan(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("LessOrEqualThan", args))
}

func (status *clearExpectStatus) Between(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("Between", args))
}

func (status *clearExpectStatus) NotBetween(code ...int) IStep {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return removeStep(status.clearPath().Push("NotBetween", args))
}

type finalClearExpectStatus struct {
	IStep
	message string
}

func (status *finalClearExpectStatus) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(status.message)
		},
	}
}

func (status *finalClearExpectStatus) Equal(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) NotEqual(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) OneOf(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) NotOneOf(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) GreaterThan(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) LessThan(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) GreaterOrEqualThan(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) LessOrEqualThan(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) Between(...int) IStep {
	return status.fail()
}
func (status *finalClearExpectStatus) NotBetween(...int) IStep {
	return status.fail()
}
