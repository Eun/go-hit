package httpbody

import (
	"encoding/json"

	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal/minitest"
)

type HttpBodyJson struct {
	// hit.Hit
	body *HttpBody
}

func newHttpBodyJson(body *HttpBody) *HttpBodyJson {
	return &HttpBodyJson{
		body: body,
		// Hit:  body.hit,
	}
}

// Get decodes the body as JSON, evaluates the expression and returns the result
//
// Example:
//     // given following response: [{"Name": "Joe"}, {"Name": "Alice"}}
//     Get("0.Name") will return "Joe"
func (jsn *HttpBodyJson) Get(expression string) interface{} {
	var container interface{}
	minitest.Panic.NoError(json.NewDecoder(jsn.body.Reader()).Decode(&container))
	v, ok, err := expr.GetValue(container, expression, expr.IgnoreCase)
	minitest.Panic.NoError(err)
	if !ok {
		v = nil
	}
	return v
}

// GetAs decodes the body as JSON, evaluates the expression and puts the result into the container and returns the container
//
// Deprecated, use Decode() or MustDecode()
func (jsn *HttpBodyJson) GetAs(expression string, container interface{}) interface{} {
	jsn.MustDecode(expression, container)
	return container
}

// Decode decodes the body as JSON, evaluates the expression and puts the result into the container
//
// Examples:
//     // given following response: [{"Name": "Joe"}, {"Name": "Alice"}}
//
//     var name string
//     Decode("0.Name", &string)
//
//     type User struct {
//         Name string
//     }
//     var user User
//     Decode("0", &user)
//
//     var users []User
//     Decode("", &users)
func (jsn *HttpBodyJson) Decode(expression string, container interface{}) error {
	var v interface{}
	if err := json.NewDecoder(jsn.body.Reader()).Decode(&v); err != nil {
		return err
	}
	v, ok, err := expr.GetValue(v, expression, expr.IgnoreCase)
	if err != nil {
		return err
	}
	if !ok {
		v = nil
	}
	return convert.Convert(v, container)
}

// MustDecode decodes the body as JSON, evaluates the expression and puts the result into the container
// it will panic if something goes wrong
func (jsn *HttpBodyJson) MustDecode(expression string, container interface{}) {
	minitest.Panic.NoError(jsn.Decode(expression, container))
}

// Set sets the body to the specified json data
func (jsn *HttpBodyJson) Set(data interface{}) {
	buf, err := json.Marshal(data)
	minitest.Panic.NoError(err)
	jsn.body.SetBytes(buf)
}
