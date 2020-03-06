package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
	"github.com/mohae/deepcopy"
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

	// custom can be used to expect a custom behaviour
	// Example:
	//           Expect().custom(func(hit Hit) {
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
	hit       Hit
	cleanPath CleanPath
	params    []interface{}
}

func newExpect(hit Hit, cleanPath CleanPath, param []interface{}) IExpect {
	return &expect{
		hit:       hit,
		cleanPath: cleanPath,
		params:    param,
	}
}

func (*expect) When() StepTime {
	return ExpectStep
}

// Exec contains the logic for Expect(...)
func (exp *expect) Exec(hit Hit) error {
	param, ok := internal.GetLastArgument(exp.params)
	if !ok {
		return errors.New("invalid argument")
	}
	return exp.Interface(param).Exec(hit)
}

func (exp *expect) CleanPath() CleanPath {
	return exp.cleanPath
}

// Body expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body("Hello World")
//           Expect().Body().Contains("Hello World")
func (exp *expect) Body(data ...interface{}) IExpectBody {
	return newExpectBody(exp, exp.hit, exp.cleanPath.Push("Body", data), data)
}

// custom can be used to expect a custom behaviour
// Example:
//           Expect().custom(func(hit Hit){
//               if hit.Response().StatusCode != 200 {
//                   panic("Expected 200")
//               }
//           })
func (exp *expect) Custom(f Callback) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: exp.cleanPath.Push("custom"),
		Instance:  exp.hit,
		Exec:      f,
	})
}

// Headers gets all headers
// Examples:
//           Expect().Headers().Contains("Content-Type")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (exp *expect) Headers() IExpectHeaders {
	return newExpectHeaders(exp, exp.hit, exp.cleanPath.Push("Headers"))
}

// Header gets the specified header
// Example:
//           Expect().Header("Content-Type").Equal("application/json")
func (exp *expect) Header(name string) IExpectSpecificHeader {
	return newExpectSpecificHeader(exp, exp.hit, exp.cleanPath.Push("Header"), name)
}

// Status expects the status to be the specified code, omit the code to get more options
// Examples:
//           Expect().Status(200)
//           Expect().Status().Equal(200)
func (exp *expect) Status(code ...int) IExpectStatus {
	return newExpectStatus(exp, exp.hit, exp.cleanPath.Push("Status"), code)
}

// Interface expects the specified interface
func (exp *expect) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return custom(Step{
			When:      ExpectStep,
			CleanPath: exp.cleanPath.Push("Interface"),
			Instance:  exp.hit,
			Exec:      x,
		})
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return custom(Step{
				When:      ExpectStep,
				CleanPath: exp.cleanPath.Push("Interface"),
				Instance:  exp.hit,
				Exec: func(hit Hit) {
					internal.CallGenericFunc(f)
				},
			})
		}
		return custom(Step{
			When:      ExpectStep,
			CleanPath: exp.cleanPath.Push("Body"),
			Instance:  exp.hit,
			Exec: func(hit Hit) {
				hit.Expect().Body(data)
			},
		})
	}
}

func makeCompareable(in, data interface{}) (interface{}, error) {
	compareData := deepcopy.Copy(data)
	err := converter.Convert(in, &compareData)
	if err != nil {
		return nil, err
	}

	return compareData, nil
}
