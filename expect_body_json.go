package hit

import (
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectBodyJSON interface {
	IStep
	Equal(expression string, data interface{}) IStep
	NotEqual(expression string, data interface{}) IStep
	Contains(expression string, data interface{}) IStep
	NotContains(expression string, data interface{}) IStep
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
//           Giving following response: { "ID": 10, "Name": "Joe" }
//           Expect().Body().JSON().Equal("", map[string]interface{}{"ID": 10, "Name": "Joe"})
//           Expect().Body().JSON().Equal("ID", 10)

//           Giving following response: { "ID": 10, "Name": "12", "Names": ["12"] }
//           Expect().Body().JSON().Equal("Names", []string{"12"})
func (jsn *expectBodyJSON) Equal(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v := hit.Response().body.JSON().Get(expression)
		if v == nil && data == nil {
			return
		}

		if v == nil || data == nil {
			// will fail
			minitest.Equal(data, v)
		}

		compareData, err := makeCompareable(v, data)
		minitest.NoError(err)
		minitest.Equal(data, compareData)
	})
}

// NotEqual expects the json body to be not equal the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().NotEqual("", map[string]interface{}{"ID": 10, "Name": "Joe"})
//           Expect().Body().JSON().NotEqual("ID", 10)
func (jsn *expectBodyJSON) NotEqual(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v := hit.Response().body.JSON().Get(expression)
		if v == nil && data == nil {
			minitest.Errorf("should not be %s", minitest.PrintValue(v))
		}

		if v == nil || data == nil {
			minitest.NotEqual(data, v)
		}

		compareData, err := makeCompareable(v, data)
		minitest.NoError(err)
		minitest.NotEqual(data, compareData)
	})
}

// Contains expects the json body to contain the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().Contains("", "ID")
//           Expect().Body().JSON().Contains("Name", "J")
func (jsn *expectBodyJSON) Contains(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v := hit.Response().body.JSON().Get(expression)
		if v == nil && data == nil {
			return
		}

		if !internal.Contains(v, data) {
			minitest.Errorf("%s does not contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
		}
	})
}

// NotContains expects the json body to not contain the specified value, the first parameter can be used to narrow down the search path
// Example:
//           Giving following response: { "ID": 10, "Name": Joe }
//           Expect().Body().JSON().NotContains("", "ID")
//           Expect().Body().JSON().NotContains("Name", "J")
func (jsn *expectBodyJSON) NotContains(expression string, data interface{}) IStep {
	return jsn.expect.Custom(func(hit Hit) {
		v := hit.Response().body.JSON().Get(expression)
		if v == nil && data == nil {
			minitest.Errorf("%s does contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
		}

		if internal.Contains(v, data) {
			minitest.Errorf("%s does contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
		}
	})
}

type finalExpectBodyJSON struct {
	IStep
}

func (finalExpectBodyJSON) Equal(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) NotEqual(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) Contains(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) NotContains(expression string, data interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}
