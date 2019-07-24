package hit

import "github.com/Eun/go-hit/internal"

type Expect interface {
	Hit
	// Body expects the body to be equal the specified value, omit the parameter to get more options
	// Examples:
	//           Expect().Body("Hello World")
	//           Expect().Body().Contains("Hello World")
	Body(data ...interface{}) *expectBody

	// Interface expects the specified interface
	Interface(interface{}) Hit

	// Custom can be used to expect a custom behaviour
	// Example:
	//           Expect().Custom(func(hit Hit){
	//               if hit.Response().StatusCode != 200 {
	//                   hit.T().FailNow()
	//               }
	//           })
	Custom(f Callback) Hit

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

	// Clear removes all previous expect calls
	// Example:
	//           Expect().Status(200).  // Will be ignored
	//           Clear().
	//           Expect().Status(404)
	Clear() Hit

	// CollectedExpects returns all the collected expect calls
	CollectedExpects() []Callback
}

type defaultExpect struct {
	Hit
	expectCalls []Callback
}

func newExpect(hit Hit) *defaultExpect {
	return &defaultExpect{
		Hit: hit,
	}
}

// Body expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body("Hello World")
//           Expect().Body().Contains("Hello World")
func (exp *defaultExpect) Body(data ...interface{}) *expectBody {
	body := newExpectBody(exp)
	if arg := getLastArgument(data); arg != nil {
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
func (exp *defaultExpect) Custom(f Callback) Hit {
	switch exp.Hit.State() {
	case Done, Working:
		f(exp.Hit)
	default: // ready
		exp.expectCalls = append(exp.expectCalls, f)
	}
	return exp.Hit
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

// Clear removes all previous expect calls
// Example:
//           Expect().Status(200).  // Will be ignored
//           Clear().
//           Expect().Status(404)
func (exp *defaultExpect) Clear() Hit {
	exp.expectCalls = nil
	return exp.Hit
}

// Interface expects the specified interface
func (exp *defaultExpect) Interface(data interface{}) Hit {
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

// CollectedExpects returns all the collected expect calls
func (exp *defaultExpect) CollectedExpects() []Callback {
	return exp.expectCalls
}

func (exp *defaultExpect) copy(toHit Hit) *defaultExpect {
	n := &defaultExpect{
		Hit: toHit,
	}
	// copy the expect calls
	n.expectCalls = make([]Callback, len(exp.expectCalls))
	for i, v := range exp.expectCalls {
		n.expectCalls[i] = v
	}
	return n
}
