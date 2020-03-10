package hit

import (
	"github.com/Eun/go-hit/internal"
	"github.com/mohae/deepcopy"
	"golang.org/x/xerrors"
)

type IExpect interface {
	IStep
	// Body expects the body to be equal the specified value, omit the parameter to get more options
	// Examples:
	//           Expect().Body("Hello World")
	//           Expect().Body().Contains("Hello World")
	Body(data ...interface{}) IExpectBody

	// Interface expects the specified interface
	Interface(data interface{}) IStep

	// custom can be used to expect a custom behaviour
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
}

type expect struct {
	cleanPath clearPath
}

func newExpect(cleanPath clearPath, params []interface{}) IExpect {
	exp := &expect{
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalExpect{&hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: cleanPath,
			Exec:      exp.Interface(param).exec,
		}}
	}
	return exp
}

func (*expect) exec(hit Hit) error {
	return xerrors.New("unsupported")
}

func (*expect) when() StepTime {
	return ExpectStep
}

func (exp *expect) clearPath() clearPath {
	return exp.cleanPath
}

// Body expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body("Hello World")
//           Expect().Body().Contains("Hello World")
func (exp *expect) Body(data ...interface{}) IExpectBody {
	return newExpectBody(exp, exp.cleanPath.Push("Body", data), data)
}

// custom can be used to expect a custom behaviour
// Example:
//           Expect().custom(func(hit Hit){
//               if hit.Response().StatusCode != 200 {
//                   panic("Expected 200")
//               }
//           })
func (exp *expect) Custom(f Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.cleanPath.Push("Custom", []interface{}{f}),
		Exec: func(hit Hit) error {
			f(hit)
			return nil
		},
	}
}

// Headers gets all headers
// Examples:
//           Expect().Headers().Contains("Content-Type")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (exp *expect) Headers() IExpectHeaders {
	return newExpectHeaders(exp, exp.cleanPath.Push("Headers", nil))
}

// Header gets the specified header
// Example:
//           Expect().Header("Content-Type").Equal("application/json")
func (exp *expect) Header(name string) IExpectSpecificHeader {
	return newExpectSpecificHeader(exp, exp.cleanPath.Push("Header", []interface{}{name}), name)
}

// Status expects the status to be the specified code, omit the code to get more options
// Examples:
//           Expect().Status(200)
//           Expect().Status().Equal(200)
func (exp *expect) Status(code ...int) IExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newExpectStatus(exp, exp.cleanPath.Push("Status", args), code)
}

// Interface expects the specified interface
func (exp *expect) Interface(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.cleanPath.Push("Interface", []interface{}{data}),
		Exec:      exp.Body(data).exec,
	}
}

type finalExpect struct {
	IStep
}

func (f finalExpect) Body(...interface{}) IExpectBody {
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

func (f finalExpect) Header(string) IExpectSpecificHeader {
	panic("only usable with Expect() not with Expect(value)")
}

func (f finalExpect) Status(...int) IExpectStatus {
	panic("only usable with Expect() not with Expect(value)")
}

func makeCompareable(in, data interface{}) (interface{}, error) {
	compareData := deepcopy.Copy(data)
	err := converter.Convert(in, &compareData)
	if err != nil {
		return nil, err
	}

	return compareData, nil
}
