package hit

import "github.com/Eun/go-hit/internal"

// IClearExpectBody provides a clear functionality to remove previous steps from running in the Expect().Body() scope
type IClearExpectBody interface {
	IStep
	// JSON removes all previous Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON()
	// e.g. Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}).
	//
	// If you specify an argument it will only remove the Expect().Body().JSON() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().JSON()                                              // will remove all Expect().Body().JSON() steps and all chained steps to JSON() e.g. Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"})
	//     Clear().Expect().Body().JSON(map[string]interface{}{"Name": "Joe"})         // will remove all Expect().Body().JSON(map[string]interface{}{"Name": "Joe"}) steps
	//     Clear().Expect().Body().JSON().Equal()                                      // will remove all Expect().Body().JSON().Equal() steps
	//     Clear().Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}) // will remove all Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"}) steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Expect().Body().JSON(),
	//         Expect().Body().JSON(map[string]interface{}{"Name": "Alice"}),
	//     )

	JSON(...interface{}) IClearExpectBodyJSON
	// Equal removes all previous Expect().Body().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Equal() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().Equal()              // will remove all Expect().Body().Equal() steps
	//     Clear().Expect().Body().Equal("Hello World") // will remove all Expect().Body().Equal("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().Equal("Hello Earth"),
	//         Clear().Expect().Body().Equal(),
	//         Expect().Body().Equal("Hello World"),
	//     )
	Equal(...interface{}) IStep

	// NotEqual removes all previous Expect().Body().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotEqual() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().NotEqual()              // will remove all Expect().Body().NotEqual() steps
	//     Clear().Expect().Body().NotEqual("Hello World") // will remove all Expect().Body().NotEqual("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().NotEqual("Hello World"),
	//         Clear().Expect().Body().NotEqual(),
	//         Expect().Body().NotEqual("Hello Earth"),
	//     )
	NotEqual(...interface{}) IStep

	// Contains removes all previous Expect().Body().Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Contains() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().Contains()              // will remove all Expect().Body().Contains() steps
	//     Clear().Expect().Body().Contains("Hello World") // will remove all Expect().Body().Contains("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().Contains("Hello Earth"),
	//         Clear().Expect().Body().Contains(),
	//         Expect().Body().Contains("Hello World"),
	//     )
	Contains(...interface{}) IStep

	// NotContains removes all previous Expect().Body().NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotContains() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().NotContains()              // will remove all Expect().Body().NotContains() steps
	//     Clear().Expect().Body().NotContains("Hello World") // will remove all Expect().Body().NotContains("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().NotContains("Hello World"),
	//         Clear().Expect().Body().NotContains(),
	//         Expect().Body().NotContains("Hello Earth"),
	//     )
	NotContains(...interface{}) IStep
}

type clearExpectBody struct {
	clearExpect IClearExpect
	cleanPath   clearPath
}

func newClearExpectBody(exp IClearExpect, cleanPath clearPath, params []interface{}) IClearExpectBody {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect().Body(something)
		return finalClearExpectBody{removeStep(cleanPath)}
	}
	return &clearExpectBody{
		clearExpect: exp,
		cleanPath:   cleanPath,
	}
}

func (body *clearExpectBody) when() StepTime {
	return CleanStep
}

func (body *clearExpectBody) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Body()
	removeSteps(hit, body.cleanPath)
	return nil
}

func (body *clearExpectBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *clearExpectBody) JSON(data ...interface{}) IClearExpectBodyJSON {
	return newClearExpectBodyJSON(body, body.cleanPath.Push("JSON", data), data)
}

func (body *clearExpectBody) Equal(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Equal", data))
}

func (body *clearExpectBody) NotEqual(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("NotEqual", data))
}

func (body *clearExpectBody) Contains(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Contains", data))
}

func (body *clearExpectBody) NotContains(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("NotContains", data))
}

type finalClearExpectBody struct {
	IStep
}

func (finalClearExpectBody) JSON(...interface{}) IClearExpectBodyJSON {
	panic("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)")
}
func (finalClearExpectBody) Interface(...interface{}) IStep {
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
