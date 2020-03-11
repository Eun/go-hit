package hit

import (
	"github.com/Eun/go-hit/internal"
)

// IClearExpectStatus provides a clear functionality to remove previous steps from running in the Expect().Status() scope
type IClearExpectStatus interface {
	IStep
	// Equal removes all previous Expect().Status().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().Equal() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().Equal()              // will remove all Expect().Status().Equal() steps
	//     Clear().Expect().Status().Equal(http.StatusOK) // will remove all Expect().Status().Equal(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().Equal(http.StatusNotFound),
	//         Clear().Expect().Status().Equal(),
	//         Expect().Status().Equal(http.StatusOK),
	//     )
	Equal(...int) IStep

	// NotEqual removes all previous Expect().Status().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotEqual() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().NotEqual()              // will remove all Expect().Status().NotEqual() steps
	//     Clear().Expect().Status().NotEqual(http.StatusOK) // will remove all Expect().Status().NotEqual(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().NotEqual(http.StatusOK),
	//         Clear().Expect().Status().NotEqual(),
	//         Expect().Status().NotEqual(http.StatusNotFound),
	//     )
	NotEqual(...int) IStep

	// OneOf removes all previous Expect().Status().OneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().OneOf() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().OneOf()              // will remove all Expect().Status().OneOf() steps
	//     Clear().Expect().Status().OneOf(http.StatusOK) // will remove all Expect().Status().OneOf(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().OneOf(http.StatusOK, http.StatusNoContent),
	//         Clear().Expect().Status().OneOf(),
	//         Expect().Status().OneOf(http.StatusOK),
	//     )
	OneOf(...int) IStep

	// NotOneOf removes all previous Expect().Status().NotOneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotOneOf() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().NotOneOf()              // will remove all Expect().Status().NotOneOf() steps
	//     Clear().Expect().Status().NotOneOf(http.StatusOK) // will remove all Expect().Status().NotOneOf(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().NotOneOf(http.StatusOK, http.StatusNoContent),
	//         Clear().Expect().Status().NotOneOf(),
	//         Expect().Status().NotOneOf(http.StatusNotFound),
	//     )
	NotOneOf(...int) IStep

	// GreaterThan removes all previous Expect().Status().GreaterThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterThan() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().GreaterThan()              // will remove all Expect().Status().GreaterThan() steps
	//     Clear().Expect().Status().GreaterThan(http.StatusOK) // will remove all Expect().Status().GreaterThan(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterThan(http.StatusContinue),
	//         Clear().Expect().Status().GreaterThan(),
	//         Expect().Status().GreaterThan(http.StatusOK),
	//     )
	GreaterThan(...int) IStep

	// LessThan removes all previous Expect().Status().LessThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().LessThan() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().LessThan()              // will remove all Expect().Status().LessThan() steps
	//     Clear().Expect().Status().LessThan(http.StatusOK) // will remove all Expect().Status().LessThan(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().LessThan(http.StatusBadRequest),
	//         Clear().Expect().Status().LessThan(),
	//         Expect().Status().LessThan(http.StatusInternalServerError),
	//     )
	LessThan(...int) IStep

	// GreaterOrEqualThan removes all previous Expect().Status().GreaterOrEqualThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterOrEqualThan() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().GreaterOrEqualThan()              // will remove all Expect().Status().GreaterOrEqualThan() steps
	//     Clear().Expect().Status().GreaterOrEqualThan(http.StatusOK) // will remove all Expect().Status().GreaterOrEqualThan(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterOrEqualThan(http.StatusBadRequest),
	//         Clear().Expect().Status().GreaterOrEqualThan(),
	//         Expect().Status().GreaterOrEqualThan(http.StatusOK),
	//     )
	GreaterOrEqualThan(...int) IStep

	// GreaterOrEqualThan removes all previous Expect().Status().GreaterOrEqualThan() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().GreaterOrEqualThan() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().GreaterOrEqualThan()              // will remove all Expect().Status().GreaterOrEqualThan() steps
	//     Clear().Expect().Status().GreaterOrEqualThan(http.StatusOK) // will remove all Expect().Status().GreaterOrEqualThan(http.StatusOK) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().GreaterOrEqualThan(http.StatusBadRequest),
	//         Clear().Expect().Status().GreaterOrEqualThan(),
	//         Expect().Status().GreaterOrEqualThan(http.StatusOK),
	//     )
	LessOrEqualThan(...int) IStep

	// Between removes all previous Expect().Status().Between() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().Between() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().Between()              // will remove all Expect().Status().Between() steps
	//     Clear().Expect().Status().Between(http.StatusOK) // will remove all Expect().Status().Between(http.StatusOK) steps
	//     Clear().Expect().Status().Between(http.StatusOK, http.StatusAccepted) // will remove all Expect().Status().Between(http.StatusOK, http.StatusAccepted) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().Between(http.StatusOK, http.StatusAccepted),
	//         Clear().Expect().Status().Between(),
	//         Expect().Status().Between(http.StatusBadRequest, http.StatusUnavailableForLegalReasons),
	//     )
	Between(...int) IStep

	// NotBetween removes all previous Expect().Status().NotBetween() steps.
	//
	// If you specify an argument it will only remove the Expect().Status().NotBetween() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Status().NotBetween()              // will remove all Expect().Status().NotBetween() steps
	//     Clear().Expect().Status().NotBetween(http.StatusOK) // will remove all Expect().Status().NotBetween(http.StatusOK) steps
	//     Clear().Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted) // will remove all Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Status().NotBetween(http.StatusOK, http.StatusAccepted),
	//         Clear().Expect().Status().NotBetween(),
	//         Expect().Status().NotBetween(http.StatusBadRequest, http.StatusUnavailableForLegalReasons),
	//     )
	NotBetween(...int) IStep
}

type clearExpectStatus struct {
	cleanPath clearPath
}

func newClearExpectStatus(cleanPath clearPath, params []int) IClearExpectStatus {
	if _, ok := internal.GetLastIntArgument(params); ok {
		// this runs if we called Clear().Expect().Status(something)
		return finalClearExpectStatus{removeStep(cleanPath)}
	}
	return &clearExpectStatus{
		cleanPath: cleanPath,
	}
}

func (status *clearExpectStatus) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Status()
	removeSteps(hit, status.cleanPath)
	return nil
}

func (*clearExpectStatus) when() StepTime {
	return CleanStep
}

func (status *clearExpectStatus) clearPath() clearPath {
	return status.cleanPath
}

func (status *clearExpectStatus) Equal(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("Equal", args))
}

func (status *clearExpectStatus) NotEqual(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("NotEqual", args))
}

func (status *clearExpectStatus) OneOf(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("OneOf", args))
}

func (status *clearExpectStatus) NotOneOf(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("NotOneOf", args))
}

func (status *clearExpectStatus) GreaterThan(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("GreaterThan", args))
}

func (status *clearExpectStatus) LessThan(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("LessThan", args))
}

func (status *clearExpectStatus) GreaterOrEqualThan(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("GreaterOrEqualThan", args))
}

func (status *clearExpectStatus) LessOrEqualThan(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("LessOrEqualThan", args))
}

func (status *clearExpectStatus) Between(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("Between", args))
}

func (status *clearExpectStatus) NotBetween(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(status.cleanPath.Push("NotBetween", args))
}

type finalClearExpectStatus struct {
	IStep
}

func (finalClearExpectStatus) Equal(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) NotEqual(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) OneOf(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) NotOneOf(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) GreaterThan(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) LessThan(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) GreaterOrEqualThan(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) LessOrEqualThan(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) Between(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
func (finalClearExpectStatus) NotBetween(...int) IStep {
	panic("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)")
}
