package hit

type IClearExpectBody interface {
	IStep
	// JSON removes all Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON(), e.g. Expect().Body().JSON().Equal()
	// Examples:
	//           Clear().Expect().Body().JSON()
	//           Clear().Expect().Body().JSON().Equal()
	JSON() IClearExpectBodyJSON
	// Equal removes all Expect().Body().Equal() steps
	Equal() IStep
	// NotEqual removes all Expect().Body().NotEqual() steps
	NotEqual() IStep
	// Contains removes all Expect().Body().Contains() steps
	Contains() IStep
	// JSON removes all Expect().Body().NotContains() steps
	NotContains() IStep
}

type clearExpectBody struct {
	clearExpect IClearExpect
	hit         Hit
	cleanPath   CleanPath
}

func newClearExpectBody(exp IClearExpect, hit Hit, cleanPath CleanPath) IClearExpectBody {
	return &clearExpectBody{
		clearExpect: exp,
		hit:         hit,
		cleanPath:   cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect().Body()

func (body *clearExpectBody) When() StepTime {
	return CleanStep
}

// Exec contains the logic for Clear().Expect().Body(), it will remove all Expect().Body() and Expect().Body().* Steps
func (body *clearExpectBody) Exec(hit Hit) error {
	removeSteps(hit, body.cleanPath)
	return nil
}

func (body *clearExpectBody) CleanPath() CleanPath {
	return body.cleanPath
}

// JSON removes all Expect().Body().JSON() steps and all steps chained to Expect().Body().JSON(), e.g. Expect().Body().JSON().Equal()
// Examples:
//           Clear().Expect().Body().JSON()
//           Clear().Expect().Body().JSON().Equal()
func (body *clearExpectBody) JSON() IClearExpectBodyJSON {
	return newClearExpectBodyJSON(body, body.hit, body.cleanPath.Push("JSON"))
}

// Equal removes all Expect().Body().Equal() steps
func (body *clearExpectBody) Equal() IStep {
	return removeStep(body.hit, body.cleanPath.Push("Equal"))
}

// NotEqual removes all Expect().Body().NotEqual() steps
func (body *clearExpectBody) NotEqual() IStep {
	return removeStep(body.hit, body.cleanPath.Push("NotEqual"))
}

// Contains removes all Expect().Body().Contains() steps
func (body *clearExpectBody) Contains() IStep {
	return removeStep(body.hit, body.cleanPath.Push("Contains"))
}

// NotContains removes all Expect().Body().NotContains() steps
func (body *clearExpectBody) NotContains() IStep {
	return removeStep(body.hit, body.cleanPath.Push("NotContains"))
}
