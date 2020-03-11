package hit

import (
	"github.com/Eun/go-hit/internal"
)

// IClearExpect provides a clear functionality to remove previous steps from running in the Expect() scope
type IClearExpect interface {
	IStep
	// Body removes all previous Expect().Body() steps and all steps chained to Expect().Body() e.g. Expect().Body().Equal("Hello World").
	//
	// If you specify an argument it will only remove the Expect().Body() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Body()                      // will remove all Expect().Body() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
	//     Clear().Expect().Body("Hello World")         // will remove all Expect().Body("Hello World") steps
	//     Clear().Expect().Body().Equal()              // will remove all Expect().Body().Equal() steps
	//     Clear().Expect().Body().Equal("Hello World") // will remove all Expect().Body().Equal("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body("Hello Earth"),
	//         Clear().Expect().Body(),
	//         Expect().Body("Hello World"),
	//     )
	Body(...interface{}) IClearExpectBody

	// Interface removes all previous Expect().Interface() steps.
	//
	// If you specify an argument it will only remove the Expect().Interface() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Interface()              // will remove all Expect().Interface() steps
	//     Clear().Expect().Interface("Hello World") // will remove all Expect().Interface("Hello World") steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Interface("Hello Earth"),
	//         Clear().Expect().Interface(),
	//         Expect().Interface("Hello World"),
	//     )
	Interface(...interface{}) IStep

	// Custom removes all previous Expect().Custom() steps.
	//
	// If you specify an argument it will only remove the Expect().Custom() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Custom(fn) // will remove all Expect().Custom(fn) steps
	//     Clear().Expect().Custom()   // will remove all Expect().Custom() steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Custom(func(hit Hit) {
	//             if hit.Response().Body().String() != "Hello Earth" {
	//                 panic("Expected Hello Earth")
	//             }
	//         }),
	//         Clear().Expect().Custom(),
	//         Expect().Custom(func(hit Hit) {
	//             if hit.Response().Body().String() != "Hello World" {
	//                 panic("Expected Hello World")
	//             }
	//         }),
	//     )
	Custom(...Callback) IStep

	// Headers removes all previous Expect().Headers() steps and all steps chained to Expect().Headers()
	// e.g. Expect().Headers().Contains("Content-Type").
	//
	// Examples:
	//     Clear().Expect().Headers()                          // will remove all Expect().Headers() steps
	//     Clear().Expect().Headers().Contains()               // will remove all Expect().Headers().Contains() steps
	//     Clear().Expect().Headers().Contains("Content-Type") // will remove all Expect().Headers().Contains("Content-Type") steps
	Headers() IClearExpectHeaders

	// Header removes all previous Expect().Header() steps and all steps chained to Expect().Header()
	// e.g. Expect().Header("Content-Type").Equal("Content-Type").
	//
	// If you specify an argument it will only remove the Expect().Header() steps matching that argument.
	//
	// Examples:
	//     Clear().Expect().Header()                       // will remove all Expect().Header() steps and all chained steps to Expect().Header() e.g Expect().Header("Content-Type").Equal("Content-Type")
	//     Clear().Expect().Header("Content-Type")         // will remove all Expect().Header("Content-Type") steps and all chained steps to Expect().Header("Content-Type") e.g. Expect().Header("Content-Type").Equal("application/json")
	//     Clear().Expect().Header("Content-Type").Equal() // will remove all Expect().Header("Content-Type").Equal() steps
	//
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Equal("application/json"),
	//         Clear().Expect().Header(),
	//         Expect().Header("Content-Type").Equal("application/octet-stream"),
	//     )
	Header(...string) IClearExpectSpecificHeader

	// Status removes all previous Expect().Status() steps and all steps chained to Expect().Status()
	// e.g. Expect().Status().Equal(http.StatusOK).
	//
	// If you specify an argument it will only remove the Expect().Status() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Status()              // will remove all Expect().Status() steps and all chained steps to Expect().Status() e.g Expect().Status().Equal(http.StatusOK)
	//           Clear().Expect().Status(http.StatusOK) // will remove all Expect().Status(http.StatusOK) steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Status(http.StatusNotFound),
	//               Clear().Expect().Status(),
	//               Expect().Status(http.StatusOK),
	//           )
	Status(...int) IClearExpectStatus
}

type clearExpect struct {
	cleanPath clearPath
	params    []interface{}
}

func newClearExpect(cleanPath clearPath, params []interface{}) IClearExpect {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect(something)
		return finalClearExpect{removeStep(cleanPath)}
	}
	return &clearExpect{
		cleanPath: cleanPath,
	}
}

func (exp *clearExpect) when() StepTime {
	return CleanStep
}

func (exp *clearExpect) exec(hit Hit) error {
	// this runs if we called Clear().Expect()
	removeSteps(hit, exp.cleanPath)
	return nil
}

func (exp *clearExpect) clearPath() clearPath {
	return exp.cleanPath
}

func (exp *clearExpect) Body(data ...interface{}) IClearExpectBody {
	return newClearExpectBody(exp, exp.cleanPath.Push("Body", data), data)
}

func (exp *clearExpect) Interface(data ...interface{}) IStep {
	return removeStep(exp.cleanPath.Push("Interface", data))
}

func (exp *clearExpect) Custom(v ...Callback) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(exp.cleanPath.Push("Custom", args))
}

func (exp *clearExpect) Headers() IClearExpectHeaders {
	return newClearExpectHeaders(exp, exp.cleanPath.Push("Headers", nil))
}

func (exp *clearExpect) Header(name ...string) IClearExpectSpecificHeader {
	args := make([]interface{}, len(name))
	for i := range name {
		args[i] = name[i]
	}
	return newClearExpectSpecificHeader(exp, exp.cleanPath.Push("Header", args))
}

func (exp *clearExpect) Status(code ...int) IClearExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newClearExpectStatus(exp.cleanPath.Push("Status", args), code)
}

type finalClearExpect struct {
	IStep
}

func (finalClearExpect) Body(...interface{}) IClearExpectBody {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}

func (finalClearExpect) Interface(...interface{}) IStep {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}

func (finalClearExpect) Custom(...Callback) IStep {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}

func (finalClearExpect) Headers() IClearExpectHeaders {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}

func (finalClearExpect) Header(...string) IClearExpectSpecificHeader {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}

func (finalClearExpect) Status(...int) IClearExpectStatus {
	panic("only usable with Clear().Expect() not with Clear().Expect(value)")
}
