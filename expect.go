package hit

import "github.com/Eun/go-hit/internal"

type IExpect interface {
	Step
	// Body expects the body to be equal the specified value, omit the parameter to get more options
	// Examples:
	//           Expect().Body("Hello World")
	//           Expect().Body().Contains("Hello World")
	Body(data ...interface{}) *expectBody

	// Interface expects the specified interface
	Interface(interface{}) Step

	// Custom can be used to expect a custom behaviour
	// Example:
	//           Expect().Custom(func(hit Hit){
	//               if hit.Response().StatusCode != 200 {
	//                   hit.T().FailNow()
	//               }
	//           })
	Custom(f Callback) Step

	// Headers gets the specified header, omit the parameter to get all headers
	// Examples:
	//           Expect().Headers("Content-Type").Equal("application/json")
	//           Expect().Headers().Contains("Content-Type")
	Headers(name ...string) *expectHeaders

	// Status expects the status to be the specified code, omit the code to get more options
	// Examples:
	//           Expect().Status(200)
	//           Expect().Status().Equal(200)
	Status(code ...int) *expectStatus

	// Clear removes all previous expect steps
	// Example:
	//           Expect().Status(200).  // Will be ignored
	//           Expect().Clear().
	//           Expect().Status(404)
	Clear() Step
}

type defaultExpect struct {
	call Callback
}

func newExpect() *defaultExpect {
	return &defaultExpect{}
}

func (exp *defaultExpect) when() StepTime {
	return ExpectStep
}

func (exp *defaultExpect) exec(hit Hit) {
	exp.call(hit)
}

// Body expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body("Hello World")
//           Expect().Body().Contains("Hello World")
func (exp *defaultExpect) Body(data ...interface{}) *expectBody {
	body := newExpectBody(exp)
	if arg, ok := getLastArgument(data); ok {
		body.Equal(arg)
	}
	return body
}

// Custom can be used to expect a custom behaviour
// Example:
//           Expect().Custom(func(hit Hit){
//               if hit.Response().StatusCode != 200 {
//                   hit.T().FailNow()
//               }
//           })
func (exp *defaultExpect) Custom(f Callback) Step {
	exp.call = f
	return exp
}

// Headers gets the specified header, omit the parameter to get all headers
// Examples:
//           Expect().Headers("Content-Type").Equal("application/json")
//           Expect().Headers().Contains("Content-Type")
func (exp *defaultExpect) Headers(name ...string) *expectHeaders {
	if size := len(name); size > 0 {
		return newExpectHeaders(exp, name[size-1])
	}
	return newExpectHeaders(exp, "")
}

// Status expects the status to be the specified code, omit the code to get more options
// Examples:
//           Expect().Status(200)
//           Expect().Status().Equal(200)
func (exp *defaultExpect) Status(code ...int) *expectStatus {
	s := newExpectStatus(exp)
	if size := len(code); size > 0 {
		s.Equal(code[size-1])
	}
	return s
}

// Clear removes all previous expect steps
// Example:
//           Expect().Status(200).  // Will be ignored
//           Expect().Clear().
//           Expect().Status(404)
func (exp *defaultExpect) Clear() Step {
	return MakeStep(ExpectStep|CleanStep, func(hit Hit) {})
}

// Interface expects the specified interface
func (exp *defaultExpect) Interface(data interface{}) Step {
	switch x := data.(type) {
	case func(e Hit):
		return exp.Custom(x)
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return exp.Custom(func(hit Hit) {
				internal.CallGenericFunc(f)
			})
		}
		return exp.Body(data)
	}
}

type dummyExpect struct {
	Step
}

func (d *dummyExpect) Body(data ...interface{}) *expectBody {
	panic("implement me")
}

func (d *dummyExpect) Interface(interface{}) Step {
	panic("implement me")
}

func (d *dummyExpect) Custom(f Callback) Step {
	panic("implement me")
}

func (d *dummyExpect) Headers(name ...string) *expectHeaders {
	panic("implement me")
}

func (d *dummyExpect) Status(code ...int) *expectStatus {
	panic("implement me")
}

func (d *dummyExpect) Clear() Step {
	panic("implement me")
}
