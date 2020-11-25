package hit

import (
	"reflect"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectBodyJSON provides assertions on the http response json body.
type IExpectBodyJSON interface {
	// Equal expects the json body to be equal to the specified value.
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Equal(map[string]interface{}{"ID": 10, "Name": "Joe", "Roles": []string{"Admin", "User"}})
	Equal(data interface{}) IStep

	// NotEqual expects the json body to be not equal to the specified value.
	//
	// see Equal() for usage and examples
	NotEqual(data ...interface{}) IStep

	// Contains expects the json body to contain the specified value.
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Contains("ID")
	Contains(data ...interface{}) IStep

	// NotContains expects the json body to not contain the specified value.
	//
	// see Contains() for usage and examples
	NotContains(data ...interface{}) IStep

	// JQ runs an jq expression on the JSON body, the result can be asserted
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().JQ(".Name").Equal("Joe")
	//     Expect().Body().JSON().JQ(".Roles").Equal([]string{"Admin", "User"})
	//     Expect().Body().JSON().JQ(".Roles.[0]").Equal("Admin")
	JQ(expression ...string) IExpectBodyJSONJQ

	// Len provides assertions to the json body size
	//
	// Usage:
	//     Expect().Body().JSON().Len().Equal(10)
	Len() IExpectInt
}

type expectBodyJSON struct {
	expectBody IExpectBody
	cleanPath  callPath
}

func newExpectBodyJSON(expectBody IExpectBody, cleanPath callPath) IExpectBodyJSON {
	return &expectBodyJSON{
		expectBody: expectBody,
		cleanPath:  cleanPath,
	}
}

func (v *expectBodyJSON) Equal(data interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Equal", []interface{}{data}),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Decode(&obj); err != nil {
				return err
			}

			return minitest.Equal(obj, data)
		},
	}
}

func (v *expectBodyJSON) NotEqual(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotEqual", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Decode(&obj); err != nil {
				return err
			}
			return minitest.NotEqual(obj, data...)
		},
	}
}

func (v *expectBodyJSON) Contains(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Contains", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Decode(&obj); err != nil {
				return err
			}
			return minitest.Contains(obj, data...)
		},
	}
}

func (v *expectBodyJSON) NotContains(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotContains", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Decode(&obj); err != nil {
				return err
			}
			return minitest.NotContains(obj, data...)
		},
	}
}

func (v *expectBodyJSON) JQ(expression ...string) IExpectBodyJSONJQ {
	return newExpectBodyJSONJQ(v.expectBody, v.cleanPath.Push("JQ", stringSliceToInterfaceSlice(expression)), expression)
}

func (v *expectBodyJSON) Len() IExpectInt {
	return newExpectInt(v.cleanPath.Push("Len", nil), func(hit Hit) int {
		var obj interface{}
		hit.Response().body.JSON().MustDecode(&obj)
		rv := reflect.ValueOf(obj)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			return rv.Len()
		default:
			panic(xerrors.Errorf("cannot get len for %#v", obj))
		}
	})
}
