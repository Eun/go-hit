package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

// IClearExpect provides a clear functionality to remove previous steps from running in the Expect() scope
type IClearExpect interface {
	IStep
	// Body removes all previous Expect().Body() steps and all steps chained to Expect().Body() e.g. Expect().Body().Equal("Hello World").
	//
	// If you specify an argument it will only remove the Expect().Body() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body()                      // will remove all Expect().Body() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
	//     Clear().Expect().Body("Hello World")         // will remove all Expect().Body("Hello World") steps
	//     Clear().Expect().Body().Equal()              // will remove all Expect().Body().Equal() steps
	//     Clear().Expect().Body().Equal("Hello World") // will remove all Expect().Body().Equal("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Body("Hello Earth"),
	//         Clear().Expect().Body(),
	//         Expect().Body("Hello World"),
	//     )
	Body(value ...interface{}) IClearExpectBody

	// Interface removes all previous Expect().Interface() steps.
	//
	// If you specify an argument it will only remove the Expect().Interface() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Interface()              // will remove all Expect().Interface() steps
	//     Clear().Expect().Interface("Hello World") // will remove all Expect().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Interface("Hello Earth"),
	//         Clear().Expect().Interface(),
	//         Expect().Interface("Hello World"),
	//     )
	Interface(value ...interface{}) IStep

	// Custom removes all previous Expect().Custom() steps.
	//
	// If you specify an argument it will only remove the Expect().Custom() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Custom(fn) // will remove all Expect().Custom(fn) steps
	//     Clear().Expect().Custom()   // will remove all Expect().Custom() steps
	//
	// Example:
	//     MustDo(
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
	Custom(fn ...Callback) IStep

	// Header removes all previous Expect().Header() steps and all steps chained to Expect().Header()
	// e.g. Expect().Header("Content-Type").Equal("Content-Type").
	//
	// If you specify an argument it will only remove the Expect().Header() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header()                       // will remove all Expect().Header() steps and all chained steps to Expect().Header() e.g Expect().Header("Content-Type").Equal("Content-Type")
	//     Clear().Expect().Header("Content-Type")         // will remove all Expect().Header("Content-Type") steps and all chained steps to Expect().Header("Content-Type") e.g. Expect().Header("Content-Type").Equal("application/json")
	//     Clear().Expect().Header("Content-Type").Equal() // will remove all Expect().Header("Content-Type").Equal() steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Equal("application/json"),
	//         Clear().Expect().Header(),
	//         Expect().Header("Content-Type").Equal("application/octet-stream"),
	//     )
	Header(headerName ...string) IClearExpectHeader

	// Status removes all previous Expect().Status() steps and all steps chained to Expect().Status()
	// e.g. Expect().Status().Equal(http.StatusOK).
	//
	// If you specify an argument it will only remove the Expect().Status() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Status()              // will remove all Expect().Status() steps and all chained steps to Expect().Status() e.g Expect().Status().Equal(http.StatusOK)
	//     Clear().Expect().Status(http.StatusOK) // will remove all Expect().Status(http.StatusOK) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Status(http.StatusNotFound),
	//         Clear().Expect().Status(),
	//         Expect().Status(http.StatusOK),
	//     )
	Status(code ...int) IClearExpectStatus
}

type clearExpect struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearExpect(cleanPath clearPath, params []interface{}) IClearExpect {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect(something)
		return &finalClearExpect{
			removeStep(cleanPath),
			"only usable with Clear().Expect() not with Clear().Expect(value)",
		}
	}
	return &clearExpect{
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}
}

func (exp *clearExpect) when() StepTime {
	return CleanStep
}

func (exp *clearExpect) exec(hit Hit) error {
	// this runs if we called Clear().Expect()
	if err := removeSteps(hit, exp.clearPath()); err != nil {
		return exp.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (exp *clearExpect) clearPath() clearPath {
	return exp.cleanPath
}

func (exp *clearExpect) Body(value ...interface{}) IClearExpectBody {
	return newClearExpectBody(exp, exp.clearPath().Push("Body", value), value)
}

func (exp *clearExpect) Interface(value ...interface{}) IStep {
	return removeStep(exp.clearPath().Push("Interface", value))
}

func (exp *clearExpect) Custom(fn ...Callback) IStep {
	args := make([]interface{}, len(fn))
	for i := range fn {
		args[i] = fn[i]
	}
	return removeStep(exp.clearPath().Push("Custom", args))
}

func (exp *clearExpect) Header(headerName ...string) IClearExpectHeader {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}
	return newClearExpectHeader(exp.clearPath().Push("Header", args))
}

func (exp *clearExpect) Status(code ...int) IClearExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newClearExpectStatus(exp.clearPath().Push("Status", args), code)
}

type finalClearExpect struct {
	IStep
	message string
}

func (exp *finalClearExpect) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(exp.message)
		},
	}
}

func (exp *finalClearExpect) Body(...interface{}) IClearExpectBody {
	return &finalClearExpectBody{
		exp.fail(),
		exp.message,
	}
}

func (exp *finalClearExpect) Interface(...interface{}) IStep {
	return exp.fail()
}

func (exp *finalClearExpect) Custom(...Callback) IStep {
	return exp.fail()
}

func (exp *finalClearExpect) Header(...string) IClearExpectHeader {
	return &finalClearExpectHeader{
		exp.fail(),
		exp.message,
	}
}

func (exp *finalClearExpect) Status(...int) IClearExpectStatus {
	return &finalClearExpectStatus{
		exp.fail(),
		exp.message,
	}
}
