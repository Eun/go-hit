package hit

type expectBody struct {
	Hit
	expect *defaultExpect
}

func newExpectBody(expect *defaultExpect) *expectBody {
	return &expectBody{
		Hit:    expect.Hit,
		expect: expect,
	}
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) *expectBodyJSON {
	jsn := newExpectBodyJSON(body)
	if arg := getLastArgument(data); arg != nil {
		jsn.Equal("", arg)
	}
	return jsn
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) Hit {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.equalOnlyNativeTypes(data) {
			return
		}
		body.JSON().Equal("", data)
	})
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) Hit {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.containsOnlyNativeTypes(data) {
			return
		}
		body.JSON().Contains("", data)
	})
}
