package hit

import "github.com/Eun/go-hit/internal"

type IClearExpectBodyJSON interface {
	IStep
	// Equal removes all Expect().Body().JSON().Equal() steps, if you specify the expression it will only remove
	// the Expect.Body().JSON().Equal() steps with the matching expression.
	// Examples:
	//           Expect().Body().JSON().Equal()            // will remove all Equal() steps
	//           Expect().Body().JSON().Equal("data")      // will only remove Equal("data", ...) steps
	Equal(...string) IStep

	// NotEqual removes all Expect().Body().JSON().NotEqual() steps, if you specify the expression it will only remove
	// the Expect.Body().JSON().NotEqual() steps with the matching expression.
	// Examples:
	//           Expect().Body().JSON().NotEqual()            // will remove all NotEqual() steps
	//           Expect().Body().JSON().NotEqual("data")      // will only remove NotEqual("data", ...) steps
	NotEqual(...string) IStep

	// Contains removes all Expect().Body().JSON().Contains() steps, if you specify the expression it will only remove
	// the Expect.Body().JSON().Contains() steps with the matching expression.
	// Examples:
	//           Expect().Body().JSON().Contains()            // will remove all Contains() steps
	//           Expect().Body().JSON().Contains("data")      // will only remove Contains("data", ...) steps
	Contains(...string) IStep

	// NotContains removes all Expect().Body().JSON().NotContains() steps, if you specify the expression it will only remove
	// the Expect.Body().JSON().NotContains() steps with the matching expression.
	// Examples:
	//           Expect().Body().JSON().NotContains()            // will remove all NotContains() steps
	//           Expect().Body().JSON().NotContains("data")      // will only remove NotContains("data", ...) steps
	NotContains(...string) IStep
}

type clearExpectBodyJSON struct {
	clearExpectBody IClearExpectBody
	cleanPath       clearPath
}

func newClearExpectBodyJSON(body IClearExpectBody, cleanPath clearPath, params []interface{}) IClearExpectBodyJSON {
	if _, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalClearExpectBodyJSON{&hitStep{
			Trace:     ett.Prepare(),
			When:      CleanStep,
			ClearPath: nil,
			Exec: func(hit Hit) error {
				removeSteps(hit, cleanPath)
				return nil
			},
		}}
	}
	return &clearExpectBodyJSON{
		clearExpectBody: body,
		cleanPath:       cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect().Body().JSON()

func (jsn *clearExpectBodyJSON) when() StepTime {
	return CleanStep
}

// exec contains the logic for Clear().Expect().Body().JSON(), it will remove all Expect().Body().JSON() and Expect().Body().JSON().* Steps
func (jsn *clearExpectBodyJSON) exec(hit Hit) error {
	removeSteps(hit, jsn.cleanPath)
	return nil
}

func (jsn *clearExpectBodyJSON) clearPath() clearPath {
	return jsn.cleanPath
}

// Equal removes all Expect().Body().JSON().Equal() steps, if you specify the expression it will only remove
// the Expect.Body().JSON().Equal() steps with the matching expression.
// Examples:
//           Expect().Body().JSON().Equal()            // will remove all Equal() steps
//           Expect().Body().JSON().Equal("data")      // will only remove Equal("data", ...) steps
func (jsn *clearExpectBodyJSON) Equal(expression ...string) IStep {
	args := make([]interface{}, len(expression))
	for i := range expression {
		args[i] = expression[i]
	}
	return removeStep(jsn.cleanPath.Push("Equal", args))
}

// NotEqual removes all Expect().Body().JSON().NotEqual() steps, if you specify the expression it will only remove
// the Expect.Body().JSON().NotEqual() steps with the matching expression.
// Examples:
//           Expect().Body().JSON().NotEqual()            // will remove all NotEqual() steps
//           Expect().Body().JSON().NotEqual("data")      // will only remove NotEqual("data", ...) steps
func (jsn *clearExpectBodyJSON) NotEqual(expression ...string) IStep {
	args := make([]interface{}, len(expression))
	for i := range expression {
		args[i] = expression[i]
	}
	return removeStep(jsn.cleanPath.Push("NotEqual", args))
}

// Contains removes all Expect().Body().JSON().Contains() steps, if you specify the expression it will only remove
// the Expect.Body().JSON().Contains() steps with the matching expression.
// Examples:
//           Expect().Body().JSON().Contains()            // will remove all Contains() steps
//           Expect().Body().JSON().Contains("data")      // will only remove Contains("data", ...) steps
func (jsn *clearExpectBodyJSON) Contains(expression ...string) IStep {
	args := make([]interface{}, len(expression))
	for i := range expression {
		args[i] = expression[i]
	}
	return removeStep(jsn.cleanPath.Push("Contains", args))
}

// NotContains removes all Expect().Body().JSON().NotContains() steps, if you specify the expression it will only remove
// the Expect.Body().JSON().NotContains() steps with the matching expression.
// Examples:
//           Expect().Body().JSON().NotContains()            // will remove all NotContains() steps
//           Expect().Body().JSON().NotContains("data")      // will only remove NotContains("data", ...) steps
func (jsn *clearExpectBodyJSON) NotContains(expression ...string) IStep {
	args := make([]interface{}, len(expression))
	for i := range expression {
		args[i] = expression[i]
	}
	return removeStep(jsn.cleanPath.Push("NotContains", args))
}

type finalClearExpectBodyJSON struct {
	IStep
}

func (finalClearExpectBodyJSON) Equal(...string) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) NotEqual(...string) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) Contains(...string) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}

func (finalClearExpectBodyJSON) NotContains(...string) IStep {
	panic("only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)")
}
