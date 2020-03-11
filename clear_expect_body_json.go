package hit

import "github.com/Eun/go-hit/internal"

// IClearExpectBody provides a clear functionality to remove previous steps from running in the Expect().Body().JSON() scope
type IClearExpectBodyJSON interface {
	IStep
	// Equal removes all previous Expect().Body().JSON().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Equal() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().JSON().Equal()              // will remove all Expect().Body().JSON().Equal() steps
	//     Clear().Expect().Body().JSON().Equal("Name")        // will remove all Expect().Body().JSON().Equal("Name", ...) steps
	//     Clear().Expect().Body().JSON().Equal("Name", "Joe") // will remove all Expect().Body().JSON().Equal("Name", "Joe") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().Equal("Name", "Joe"),
	//         Expect().Body().JSON().Equal("Id", 10),
	//         Clear().Expect().Body().JSON().Equal("Name"),
	//         Clear().Expect().Body().JSON().Equal("Id", 10),
	//         Expect().Body().JSON().Equal("Name", "Alice"),
	//     )
	Equal(...interface{}) IStep

	// NotEqual removes all previous Expect().Body().JSON().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotEqual() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().JSON().NotEqual()              // will remove all Expect().Body().JSON().NotEqual() steps
	//     Clear().Expect().Body().JSON().NotEqual("Name")        // will remove all Expect().Body().JSON().NotEqual("Name") steps
	//     Clear().Expect().Body().JSON().NotEqual("Name", "Joe") // will remove all Expect().Body().JSON().NotEqual("Name", "Joe") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().NotEqual("Name", "Joe"),
	//         Expect().Body().JSON().NotEqual("Id", 10),
	//         Clear().Expect().Body().JSON().NotEqual("Name"),
	//         Clear().Expect().Body().JSON().NotEqual("Id", 10),
	//         Expect().Body().JSON().NotEqual("Name", "Alice"),
	//     )
	NotEqual(...interface{}) IStep

	// Contains removes all previous Expect().Body().JSON().Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Contains() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().JSON().Contains()              // will remove all Expect().Body().JSON().Contains() steps
	//     Clear().Expect().Body().JSON().Contains("Name")        // will remove all Expect().Body().JSON().Contains("Name") steps
	//     Clear().Expect().Body().JSON().Contains("Name", "Joe") // will remove all Expect().Body().JSON().Contains("Name", "Joe") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().Contains("Name", "Joe"),
	//         Expect().Body().JSON().Contains("Id", 10),
	//         Clear().Expect().Body().JSON().Contains("Name"),
	//         Clear().Expect().Body().JSON().Contains("Id", 10),
	//         Expect().Body().JSON().Contains("Name", "Alice"),
	//     )
	Contains(...interface{}) IStep

	// NotContains removes all previous Expect().Body().JSON().NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotContains() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body().JSON().NotContains()              // will remove all Expect().Body().JSON().NotContains() steps
	//     Clear().Expect().Body().JSON().NotContains("Name")        // will remove all Expect().Body().JSON().NotContains("Name") steps
	//     Clear().Expect().Body().JSON().NotContains("Name", "Joe") // will remove all Expect().Body().JSON().NotContains("Name", "Joe") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().NotContains("Name", "Joe"),
	//         Expect().Body().JSON().NotContains("Id", 10),
	//         Clear().Expect().Body().JSON().NotContains("Name"),
	//         Clear().Expect().Body().JSON().NotContains("Id", 10),
	//         Expect().Body().JSON().NotContains("Name", "Alice"),
	//     )
	NotContains(...interface{}) IStep
}

type clearExpectBodyJSON struct {
	clearExpectBody IClearExpectBody
	cleanPath       clearPath
}

func newClearExpectBodyJSON(body IClearExpectBody, cleanPath clearPath, params []interface{}) IClearExpectBodyJSON {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect().Body().JSON(something)
		return finalClearExpectBodyJSON{removeStep(cleanPath)}
	}
	return &clearExpectBodyJSON{
		clearExpectBody: body,
		cleanPath:       cleanPath,
	}
}

func (jsn *clearExpectBodyJSON) when() StepTime {
	return CleanStep
}

func (jsn *clearExpectBodyJSON) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Body().JSON()
	removeSteps(hit, jsn.cleanPath)
	return nil
}

func (jsn *clearExpectBodyJSON) clearPath() clearPath {
	return jsn.cleanPath
}

func (jsn *clearExpectBodyJSON) Equal(v ...interface{}) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(jsn.cleanPath.Push("Equal", args))
}

func (jsn *clearExpectBodyJSON) NotEqual(v ...interface{}) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(jsn.cleanPath.Push("NotEqual", args))
}

func (jsn *clearExpectBodyJSON) Contains(v ...interface{}) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(jsn.cleanPath.Push("Contains", args))
}

func (jsn *clearExpectBodyJSON) NotContains(v ...interface{}) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(jsn.cleanPath.Push("NotContains", args))
}

type finalClearExpectBodyJSON struct {
	IStep
}

func (finalClearExpectBodyJSON) Equal(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) NotEqual(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) Contains(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) NotContains(...interface{}) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}
