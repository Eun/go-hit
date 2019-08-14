package hit

type expectBody struct {
	*defaultExpect
}

func newExpectBody(expect *defaultExpect) *expectBody {
	return &expectBody{
		defaultExpect: expect,
	}
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) *expectBodyJSON {
	jsn := newExpectBodyJSON(body)
	if arg, ok := getLastArgument(data); ok {
		jsn.Equal("", arg)
	}
	return jsn
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) IStep {
	return body.defaultExpect.Custom(func(hit Hit) {
		if hit.Response().body.equalOnlyNativeTypes(data) {
			return
		}
		hit.AddSteps(body.JSON().Equal("", data))
		return
	})
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) IStep {
	return body.defaultExpect.Custom(func(hit Hit) {
		if hit.Response().body.containsOnlyNativeTypes(data) {
			return
		}
		hit.AddSteps(body.JSON().Contains("", data))
	})
}
