package hit

import (
	"fmt"

	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
)

type IExpect interface {
	IStep
	// Body expects the body to be equal the specified value, omit the parameter to get more options
	// Examples:
	//           Expect().Body("Hello World")
	//           Expect().Body().Contains("Hello World")
	Body(data ...interface{}) IExpectBody

	// Interface expects the specified interface
	Interface(interface{}) IStep

	// Custom can be used to expect a custom behaviour
	// Example:
	//           Expect().Custom(func(hit Hit) {
	//               if hit.Response().StatusCode != 200 {
	//                   panic("Expected 200")
	//               }
	//           })
	Custom(f Callback) IStep

	// Headers gets all headers
	// Examples:
	//           Expect().Headers().Contains("Content-Type")
	//           Expect().Headers().Get("Content-Type").Contains("json")
	Headers() IExpectHeaders

	// Header gets the specified header
	// Example:
	//           Expect().Headers("Content-Type").Equal("application/json")
	Header(name string) IExpectSpecificHeader

	// Status expects the status to be the specified code, omit the code to get more options
	// Examples:
	//           Expect().Status(200)
	//           Expect().Status().Equal(200)
	Status(code ...int) IExpectStatus

	// Clear removes all previous expect steps
	// Example:
	//           Do(
	//               Expect().Status(200),  // Will be ignored
	//               Expect().Clear(),
	//               Expect().Status(404),
	//           )
	Clear() IStep
}

type expect struct {
	executeNowContext Hit
	et                *errortrace.ErrorTrace
	call              Callback
}

func newExpect(executeNowContext Hit) IExpect {
	return &expect{
		executeNowContext: executeNowContext,
	}
}

func (exp *expect) when() StepTime {
	return ExpectStep
}

func (exp *expect) exec(hit Hit) (err error) {
	if exp.call == nil {
		return nil
	}
	defer func() {
		r := recover()
		if r != nil {
			err = exp.et.Format(fmt.Sprint(r))
		}
	}()
	exp.call(hit)
	return err
}

// Body expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body("Hello World")
//           Expect().Body().Contains("Hello World")
func (exp *expect) Body(data ...interface{}) IExpectBody {
	if arg, ok := getLastArgument(data); ok {
		return finalExpectBody{newExpectBody(exp).Equal(arg)}
	}
	return newExpectBody(exp)
}

// Custom can be used to expect a custom behaviour
// Example:
//           Expect().Custom(func(hit Hit){
//               if hit.Response().StatusCode != 200 {
//                   panic("Expected 200")
//               }
//           })
func (exp *expect) Custom(f Callback) IStep {
	if exp.executeNowContext != nil {
		f(exp.executeNowContext)
		return nil
	}
	exp.et = errortrace.Prepare()
	exp.call = f
	return exp
}

// Headers gets all headers
// Examples:
//           Expect().Headers().Contains("Content-Type")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (exp *expect) Headers() IExpectHeaders {
	return newExpectHeaders(exp)
}

// Header gets the specified header
// Example:
//           Expect().Header("Content-Type").Equal("application/json")
func (exp *expect) Header(name string) IExpectSpecificHeader {
	return newExpectSpecificHeader(exp, name)
}

// Status expects the status to be the specified code, omit the code to get more options
// Examples:
//           Expect().Status(200)
//           Expect().Status().Equal(200)
func (exp *expect) Status(code ...int) IExpectStatus {
	if size := len(code); size > 0 {
		return finalExpectStatus{newExpectStatus(exp).Equal(code[size-1])}
	}
	return newExpectStatus(exp)
}

// Clear removes all previous expect steps
// Example:
//           Do(
//               Expect().Status(200),  // Will be ignored
//               Expect().Clear(),
//               Expect().Status(404),
//           )
func (exp *expect) Clear() IStep {
	return Custom(ExpectStep|CleanStep, func(hit Hit) {})
}

// Interface expects the specified interface
func (exp *expect) Interface(data interface{}) IStep {
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

type finalExpect struct {
	IStep
}

func (f finalExpect) Body(data ...interface{}) IExpectBody {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Interface(interface{}) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Custom(Callback) IStep {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Headers() IExpectHeaders {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Header(name string) IExpectSpecificHeader {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Status(code ...int) IExpectStatus {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Clear() IStep {
	panic("only usable with Expect() not with Expect(value)")
}
