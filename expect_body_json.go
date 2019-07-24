package hit

import (
	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal"
)

type expectBodyJSON struct {
	Hit
	body *expectBody
}

func newExpectBodyJSON(body *expectBody) *expectBodyJSON {
	return &expectBodyJSON{
		Hit:  body.Hit,
		body: body,
	}
}

// Compare functions

// Equal expects the json body to be equal the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().Equal("", map[string]interface{}{"ID": 10, "Name": "Joe"})
//           Expect().Body().JSON().Equal("ID", 10)
func (jsn *expectBodyJSON) Equal(expression string, data interface{}) Hit {
	et := errortrace.Prepare()
	return jsn.body.expect.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		et.Panic.NoError(hit.T(), err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if v == nil || data == nil {
			// will fail
			et.Panic.Equal(hit.T(), data, v)
		}

		compareData, err := converter.Convert(v, data, convert.Options.ConvertEmbeddedStructToParentType())
		et.Panic.NoError(hit.T(), err)
		et.Panic.Equal(hit.T(), data, compareData)
	})
}

// Contains expects the json body to contain the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().Contains("", "ID")
//           Expect().Body().JSON().Contains("Name", "J")
func (jsn *expectBodyJSON) Contains(expression string, data interface{}) Hit {
	et := errortrace.Prepare()
	return jsn.body.expect.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		et.Panic.NoError(hit.T(), err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if !internal.Contains(v, data) {
			et.Panic.Errorf(hit.T(), "%#v does not contain %#v", v, data)
		}
	})
}
