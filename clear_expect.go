package hit

import (
	"github.com/Eun/go-hit/internal"
)

type IClearExpect interface {
	IStep
	// Body removes all Expect().Body() steps and all steps chained to Expect().Body(), e.g. Expect().Body().JSON()
	// Examples:
	//           Clear().Expect().Body()
	//           Clear().Expect().Body().JSON()
	Body(data ...interface{}) IClearExpectBody

	// Interface removes all Expect().Interface() steps
	Interface(data ...interface{}) IStep

	// custom removes all Expect().custom() steps
	// Example:
	//           Clear().Expect().custom()
	Custom(f ...Callback) IStep

	// Headers removes all Expect().Headers() steps and all steps chained to Expect().Headers(), e.g. Expect().Headers().Contains()
	// Examples:
	//           Clear().Expect().Headers()
	//           Clear().Expect().Headers().Contains()
	Headers() IClearExpectHeaders

	// Header removes all Expect().Header() steps and all steps chained to Expect().Header(), e.g. Expect().Header().Contains()
	// Examples:
	//           Clear().Expect().Header()
	//           Clear().Expect().Header().Equal()
	Header(header ...string) IClearExpectSpecificHeader

	// Status removes all Expect().Status() steps and all steps chained to Expect().Status(), e.g. Expect().Status().Equal()
	// Examples:
	//           Clear(),Expect().Status()
	//           Clear().Expect().Status().Equal()
	Status(code ...int) IClearExpectStatus
}

type clearExpect struct {
	clear     IClear
	cleanPath CleanPath
	params    []interface{}
}

func newClearExpect(clear IClear, cleanPath CleanPath, params []interface{}) IClearExpect {
	if _, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalClearExpect{custom(Step{
			When:      CleanStep,
			CleanPath: nil,
			Exec: func(hit Hit) {
				removeSteps(hit, cleanPath)
			},
		})}
	}
	return &clearExpect{
		clear:     clear,
		cleanPath: cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect()

func (exp *clearExpect) When() StepTime {
	return CleanStep
}

func (exp *clearExpect) Exec(hit Hit) error {
	removeSteps(hit, exp.cleanPath)
	return nil
}

func (exp *clearExpect) CleanPath() CleanPath {
	return exp.cleanPath
}

// Body removes all Expect().Body() steps and all steps chained to Expect().Body(), e.g. Expect().Body().JSON()
// Examples:
//           Clear().Expect().Body()
//           Clear().Expect().Body().JSON()
func (exp *clearExpect) Body(data ...interface{}) IClearExpectBody {
	return newClearExpectBody(exp, exp.cleanPath.Push("Body", data))
}

// Interface removes all Expect().Interface() steps
func (exp *clearExpect) Interface(data ...interface{}) IStep {
	return removeStep(exp.cleanPath.Push("Interface", data))
}

// custom removes all Expect().custom() steps
// Example:
//           Clear().Expect().custom()
func (exp *clearExpect) Custom(f ...Callback) IStep {
	args := make([]interface{}, len(f))
	for i := range f {
		args[i] = f[i]
	}
	return removeStep(exp.cleanPath.Push("Custom", args))
}

// Headers removes all Expect().Headers() steps and all steps chained to Expect().Headers(), e.g. Expect().Headers().Contains()
// Examples:
//           Clear().Expect().Headers()
//           Clear().Expect().Headers().Contains()
func (exp *clearExpect) Headers() IClearExpectHeaders {
	return newClearExpectHeaders(exp, exp.cleanPath.Push("Headers", nil))
}

// Header removes all Expect().Header() steps and all steps chained to Expect().Header(), e.g. Expect().Header().Contains()
// Examples:
//           Clear().Expect().Header()
//           Clear().Expect().Header().Equal()
func (exp *clearExpect) Header(name ...string) IClearExpectSpecificHeader {
	args := make([]interface{}, len(name))
	for i := range name {
		args[i] = name[i]
	}
	return newClearExpectSpecificHeader(exp, exp.cleanPath.Push("Header", args))
}

// Status removes all Expect().Status() steps and all steps chained to Expect().Status(), e.g. Expect().Status().Equal()
// Examples:
//           Clear(),Expect().Status()
//           Clear().Expect().Status().Equal()
func (exp *clearExpect) Status(code ...int) IClearExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newClearExpectStatus(exp, exp.cleanPath.Push("Status", args), code)
}

type finalClearExpect struct {
	IStep
}

func (f finalClearExpect) Body(...interface{}) IClearExpectBody {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalClearExpect) Interface(...interface{}) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalClearExpect) Custom(...Callback) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalClearExpect) Headers() IClearExpectHeaders {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalClearExpect) Header(...string) IClearExpectSpecificHeader {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalClearExpect) Status(...int) IClearExpectStatus {
	panic("only usable with Expect() not with Expect(value)")
}
