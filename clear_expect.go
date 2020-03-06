package hit

type IClearExpect interface {
	// Body removes all Expect().Body() steps and all steps chained to Expect().Body(), e.g. Expect().Body().JSON()
	// Examples:
	//           Clear().Expect().Body()
	//           Clear().Expect().Body().JSON()
	Body() IClearExpectBody

	// Interface removes all Expect().Interface() steps
	Interface() IStep

	// custom removes all Expect().custom() steps
	// Example:
	//           Clear().Expect().custom()
	Custom() IStep

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
	// Status() IClearExpectStatus
}

type clearExpect struct {
	clear     IClear
	hit       Hit
	cleanPath CleanPath
}

func newClearExpect(clear IClear, hit Hit, cleanPath CleanPath) IClearExpect {
	return &clearExpect{
		clear:     clear,
		hit:       hit,
		cleanPath: cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect()

func (exp *clearExpect) When() StepTime {
	return CleanStep
}

// Exec contains the logic for Clear().Expect(), it will remove all Expect() and Expect().* Steps
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
func (exp *clearExpect) Body() IClearExpectBody {
	return newClearExpectBody(exp, exp.hit, exp.cleanPath.Push("Body"))
}

// Interface removes all Expect().Interface() steps
func (exp *clearExpect) Interface() IStep {
	return removeStep(exp.hit, exp.cleanPath.Push("Interface"))
}

// custom removes all Expect().custom() steps
// Example:
//           Clear().Expect().custom()
func (exp *clearExpect) Custom() IStep {
	return removeStep(exp.hit, exp.cleanPath.Push("custom"))
}

// Headers removes all Expect().Headers() steps and all steps chained to Expect().Headers(), e.g. Expect().Headers().Contains()
// Examples:
//           Clear().Expect().Headers()
//           Clear().Expect().Headers().Contains()
func (exp *clearExpect) Headers() IClearExpectHeaders {
	return newClearExpectHeaders(exp, exp.hit, exp.cleanPath.Push("Headers"))
}

// Header removes all Expect().Header() steps and all steps chained to Expect().Header(), e.g. Expect().Header().Contains()
// Examples:
//           Clear().Expect().Header()
//           Clear().Expect().Header().Equal()
func (exp *clearExpect) Header(name ...string) IClearExpectSpecificHeader {
	return newClearExpectSpecificHeader(exp, exp.hit, exp.cleanPath.Push("Header", name))
}

// // Status removes all Expect().Status() steps and all steps chained to Expect().Status(), e.g. Expect().Status().Equal()
// // Examples:
// //           Clear(),Expect().Status()
// //           Clear().Expect().Status().Equal()
// func (*clearExpect) Status() IClearExpectStatus {
// 	panic("implement me")
// }
