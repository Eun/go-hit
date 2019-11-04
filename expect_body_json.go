package hit

import (
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectBodyJSON interface {
	IStep
	Equal(expression string, data interface{}) IStep
	Contains(expression string, data interface{}) IStep
}

type expectBodyJSON struct {
	expect IExpect
}

func newExpectBodyJSON(expect IExpect) IExpectBodyJSON {
	return &expectBodyJSON{expect}
}

func (jsn *expectBodyJSON) when() StepTime {
	return jsn.expect.when()
}

func (jsn *expectBodyJSON) exec(hit Hit) error {
	return jsn.expect.exec(hit)
}

// Compare functions

// Equal expects the json body to be equal the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().Equal("", map[string]interface{}{"ID": 10, "Name": "Joe"})
//           Expect().Body().JSON().Equal("ID", 10)
func (jsn *expectBodyJSON) Equal(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		minitest.NoError(err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if v == nil || data == nil {
			// will fail
			minitest.Equal(data, v)
		}

		compareData := data
		err = converter.Convert(v, &compareData)
		minitest.NoError(err)
		minitest.Equal(data, compareData)
	})
}

// Contains expects the json body to contain the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().Contains("", "ID")
//           Expect().Body().JSON().Contains("Name", "J")
func (jsn *expectBodyJSON) Contains(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		minitest.NoError(err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if !internal.Contains(v, data) {
			minitest.Errorf("%#v does not contain %#v", v, data)
		}
	})
}

type finalExpectBodyJSON struct {
	IStep
}

func (finalExpectBodyJSON) Equal(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) Contains(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}
