package hit

import "github.com/Eun/go-hit/internal"

type IClearExpectBody interface {
	IStep
	// JSON removes all Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON(), e.g. Expect().Body().JSON().Equal()
	// Examples:
	//           Clear().Expect().Body().JSON()
	//           Clear().Expect().Body().JSON().Equal()
	JSON(...interface{}) IClearExpectBodyJSON
	// Equal removes all Expect().Body().Equal() steps
	Equal(...interface{}) IStep
	// NotEqual removes all Expect().Body().NotEqual() steps
	NotEqual(...interface{}) IStep
	// Contains removes all Expect().Body().Contains() steps
	Contains(...interface{}) IStep
	// JSON removes all Expect().Body().NotContains() steps
	NotContains(...interface{}) IStep
}

type clearExpectBody struct {
	clearExpect IClearExpect
	cleanPath   clearPath
}

func newClearExpectBody(exp IClearExpect, cleanPath clearPath, params []interface{}) IClearExpectBody {
	if _, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalClearExpectBody{&hitStep{
			Trace:     ett.Prepare(),
			When:      CleanStep,
			ClearPath: nil,
			Exec: func(hit Hit) error {
				removeSteps(hit, cleanPath)
				return nil
			},
		}}
	}
	return &clearExpectBody{
		clearExpect: exp,
		cleanPath:   cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect().Body()

func (body *clearExpectBody) when() StepTime {
	return CleanStep
}

// exec contains the logic for Clear().Expect().Body(), it will remove all Expect().Body() and Expect().Body().* Steps
func (body *clearExpectBody) exec(hit Hit) error {
	removeSteps(hit, body.cleanPath)
	return nil
}

func (body *clearExpectBody) clearPath() clearPath {
	return body.cleanPath
}

// JSON removes all Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON(), e.g. Expect().Body().JSON().Equal()
// Examples:
//           Clear().Expect().Body().JSON()
//           Clear().Expect().Body().JSON().Equal()
func (body *clearExpectBody) JSON(data ...interface{}) IClearExpectBodyJSON {
	return newClearExpectBodyJSON(body, body.cleanPath.Push("JSON", data), data)
}

// Equal removes all Expect().Body().Equal() steps
func (body *clearExpectBody) Equal(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Equal", data))
}

// NotEqual removes all Expect().Body().NotEqual() steps
func (body *clearExpectBody) NotEqual(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("NotEqual", data))
}

// Contains removes all Expect().Body().Contains() steps
func (body *clearExpectBody) Contains(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Contains", data))
}

// NotContains removes all Expect().Body().NotContains() steps
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
