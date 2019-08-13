package hit

import (
	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal"
)

type expectBodyJSON struct {
	*expectBody
}

func newExpectBodyJSON(body *expectBody) *expectBodyJSON {
	return &expectBodyJSON{
		expectBody: body,
	}
}

// Compare functions

// Equal expects the json body to be equal the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           IExpect().Body().JSON().Equal("", map[string]interface{}{"ID": 10, "Name": "Joe"})
//           IExpect().Body().JSON().Equal("ID", 10)
func (jsn *expectBodyJSON) Equal(expression string, data interface{}) Step {
	et := errortrace.Prepare()
	return jsn.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		et.NoError(err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if v == nil || data == nil {
			// will fail
			et.Equal(data, v)
		}

		compareData, err := converter.Convert(v, data, convert.Options.ConvertEmbeddedStructToParentType())
		et.NoError(err)
		et.Equal(data, compareData)
	})
}

// Contains expects the json body to contain the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           IExpect().Body().JSON().Contains("", "ID")
//           IExpect().Body().JSON().Contains("Name", "J")
func (jsn *expectBodyJSON) Contains(expression string, data interface{}) Step {
	et := errortrace.Prepare()
	return jsn.Custom(func(hit Hit) {
		v, ok, err := expr.GetValue(hit.Response().body.JSON().Get(), expression, expr.IgnoreCase)
		et.NoError(err)
		if !ok {
			v = nil
		}

		if v == nil && data == nil {
			return
		}

		if !internal.Contains(v, data) {
			et.Errorf("%#v does not contain %#v", v, data)
		}
	})
}
