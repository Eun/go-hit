package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
)

// IClearExpectBody provides a clear functionality to remove previous steps from running in the Expect().Body() scope
type IClearExpectBody interface {
	IStep

	// Interface removes all previous Expect().Body().Interface() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Interface() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().Interface()                                              // will remove all Expect().Body().Interface() steps and all chained steps to Interface() e.g. Expect().Body().Interface().Equal(map[string]interface{}{"Name": "Joe"})
	//     Clear().Expect().Body().Interface(map[string]interface{}{"Name": "Joe"})         // will remove all Expect().Body().Interface(map[string]interface{}{"Name": "Joe"}) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Body().Interface(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Expect().Body().Interface(),
	//         Expect().Body().Interface(map[string]interface{}{"Name": "Alice"}),
	//     )
	Interface(value ...interface{}) IStep

	// JSON removes all previous Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON()
	// e.g. Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}).
	//
	// If you specify an argument it will only remove the Expect().Body().JSON() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().JSON()                                              // will remove all Expect().Body().JSON() steps and all chained steps to JSON() e.g. Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"})
	//     Clear().Expect().Body().JSON(map[string]interface{}{"Name": "Joe"})         // will remove all Expect().Body().JSON(map[string]interface{}{"Name": "Joe"}) steps
	//     Clear().Expect().Body().JSON().Equal()                                      // will remove all Expect().Body().JSON().Equal() steps
	//     Clear().Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}) // will remove all Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Expect().Body().JSON(),
	//         Expect().Body().JSON(map[string]interface{}{"Name": "Alice"}),
	//     )
	JSON(value ...interface{}) IClearExpectBodyJSON

	// Equal removes all previous Expect().Body().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Equal() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().Equal()              // will remove all Expect().Body().Equal() steps
	//     Clear().Expect().Body().Equal("Hello World") // will remove all Expect().Body().Equal("Hello World") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().Equal("Hello Earth"),
	//         Clear().Expect().Body().Equal(),
	//         Expect().Body().Equal("Hello World"),
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual removes all previous Expect().Body().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotEqual() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().NotEqual()              // will remove all Expect().Body().NotEqual() steps
	//     Clear().Expect().Body().NotEqual("Hello World") // will remove all Expect().Body().NotEqual("Hello World") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().NotEqual("Hello World"),
	//         Clear().Expect().Body().NotEqual(),
	//         Expect().Body().NotEqual("Hello Earth"),
	//     )
	NotEqual(value ...interface{}) IStep

	// Contains removes all previous Expect().Body().Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Contains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().Contains()              // will remove all Expect().Body().Contains() steps
	//     Clear().Expect().Body().Contains("Hello World") // will remove all Expect().Body().Contains("Hello World") steps
	//
	// Examples:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().Contains("Hello Earth"),
	//         Clear().Expect().Body().Contains(),
	//         Expect().Body().Contains("Hello World"),
	//     )
	Contains(value ...interface{}) IStep

	// NotContains removes all previous Expect().Body().NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotContains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().NotContains()              // will remove all Expect().Body().NotContains() steps
	//     Clear().Expect().Body().NotContains("Hello World") // will remove all Expect().Body().NotContains("Hello World") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().NotContains("Hello World"),
	//         Clear().Expect().Body().NotContains(),
	//         Expect().Body().NotContains("Hello Earth"),
	//     )
	NotContains(value ...interface{}) IStep
}

type clearExpectBody struct {
	clearExpect IClearExpect
	cleanPath   clearPath
	trace       *errortrace.ErrorTrace
}

func newClearExpectBody(exp IClearExpect, cleanPath clearPath, params []interface{}) IClearExpectBody {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect().Body(something)
		return finalClearExpectBody{removeStep(cleanPath)}
	}
	return &clearExpectBody{
		clearExpect: exp,
		cleanPath:   cleanPath,
		trace:       ett.Prepare(),
	}
}

func (body *clearExpectBody) when() StepTime {
	return CleanStep
}

func (body *clearExpectBody) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Body()
	if err := removeSteps(hit, body.clearPath()); err != nil {
		return body.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (body *clearExpectBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *clearExpectBody) Interface(value ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Interface", value))
}

func (body *clearExpectBody) JSON(value ...interface{}) IClearExpectBodyJSON {
	return newClearExpectBodyJSON(body, body.cleanPath.Push("JSON", value), value)
}

func (body *clearExpectBody) Equal(value ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Equal", value))
}

func (body *clearExpectBody) NotEqual(value ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("NotEqual", value))
}

func (body *clearExpectBody) Contains(value ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Contains", value))
}

func (body *clearExpectBody) NotContains(value ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("NotContains", value))
}

type finalClearExpectBody struct {
	IStep
}

func (finalClearExpectBody) Interface(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) JSON(...interface{}) IClearExpectBodyJSON {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) Equal(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) NotEqual(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) Contains(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) NotContains(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
