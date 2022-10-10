package hit

import (
	"reflect"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectBodyJSONDasel provides assertions on the http response json body.
type IExpectBodyJSONDasel interface {
	// Equal expects the json body to be equal to the specified value.
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Dasel(".").Equal(map[string]interface{}{"ID": 10, "Name": "Joe", "Roles": []string{"Admin", "User"}})
	//     Expect().Body().JSON().Dasel(".Name").Equal("Joe")
	//     Expect().Body().JSON().Dasel(".Roles").Equal([]string{"Admin", "User"})
	//     Expect().Body().JSON().Dasel(".Roles.[0]").Equal("Admin")
	//
	// Example:
	//     // given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Body().JSON().Dasel(".").Equal(map[string]interface{}{"ID": 10, "Name": "Joe", "Roles": []string{"Admin", "User"}}),
	//         Expect().Body().JSON().Dasel(".Name").Equal("Joe"),
	//         Expect().Body().JSON().Dasel(".Roles").Equal([]string{"Admin", "User"}),
	//         Expect().Body().JSON().Dasel(".Roles.[0]").Equal("Admin"),
	//     )
	Equal(data interface{}) IStep

	// NotEqual expects the json body to be equal to the specified value.
	//
	// see Equal() for usage and examples
	NotEqual(data ...interface{}) IStep

	// Contains expects the json body to be equal to the specified value.
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Dasel(".").Contains("ID")
	//     Expect().Body().JSON().Dasel(".Name").Contains("J")
	//     Expect().Body().JSON().Dasel(".Roles").Contains("Admin")
	//
	// Example:
	//     // given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Body().JSON().Dasel(".").Contains("ID"),
	//         Expect().Body().JSON().Dasel(".Name").Contains("J"),
	//         Expect().Body().JSON().Dasel(".Roles").Contains("Admin"),
	//     )
	Contains(data ...interface{}) IStep

	// NotContains expects the json body to be equal to the specified value.
	//
	// see Contains() for usage and examples
	NotContains(data ...interface{}) IStep

	// Len provides assertions to the json object size
	//
	// Usage:
	//     Expect().Body().JSON().Dasel(".Name").Len().Equal(10)
	Len() IExpectInt

	// JQ runs an additional jq expression
	JQ(expression ...string) IExpectBodyJSONJQ

	// Dasel runs an additional dasel expression
	Dasel(expression ...string) IExpectBodyJSONDasel
}

type expectBodyJSONDasel struct {
	expectBody IExpectBody
	cleanPath  callPath
	expression []string
}

func newExpectBodyJSONDasel(expectBody IExpectBody, cleanPath callPath, expression []string) IExpectBodyJSONDasel {
	return &expectBodyJSONDasel{
		expectBody: expectBody,
		cleanPath:  cleanPath,
		expression: expression,
	}
}

func (v *expectBodyJSONDasel) Equal(data interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Equal", []interface{}{data}),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Dasel(&obj, v.expression...); err != nil {
				return err
			}
			return minitest.Equal(obj, data)
		},
	}
}

func (v *expectBodyJSONDasel) NotEqual(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotEqual", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Dasel(&obj, v.expression...); err != nil {
				return err
			}
			return minitest.NotEqual(obj, data...)
		},
	}
}

func (v *expectBodyJSONDasel) Contains(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Contains", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Dasel(&obj, v.expression...); err != nil {
				return err
			}
			return minitest.Contains(obj, data...)
		},
	}
}

func (v *expectBodyJSONDasel) NotContains(data ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotContains", data),
		Exec: func(hit *hitImpl) error {
			var obj interface{}
			if err := hit.Response().body.JSON().Dasel(&obj, v.expression...); err != nil {
				return err
			}
			return minitest.NotContains(obj, data...)
		},
	}
}

func (v *expectBodyJSONDasel) Len() IExpectInt {
	return newExpectInt(v.cleanPath.Push("Len", nil), func(hit Hit) int {
		var obj interface{}
		hit.Response().body.JSON().MustDasel(&obj, v.expression...)
		if obj == nil {
			return 0
		}
		rv := reflect.ValueOf(obj)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			return rv.Len()
		default:
			panic(xerrors.Errorf("cannot get len for %#v", obj))
		}
	})
}

func (v *expectBodyJSONDasel) JQ(expression ...string) IExpectBodyJSONJQ {
	return newExpectBodyJSONJQ(
		v.expectBody,
		v.cleanPath.Push("JQ", stringSliceToInterfaceSlice(expression)), append(v.expression, expression...),
	)
}

func (v *expectBodyJSONDasel) Dasel(expression ...string) IExpectBodyJSONDasel {
	return newExpectBodyJSONDasel(
		v.expectBody,
		v.cleanPath.Push("Dasel", stringSliceToInterfaceSlice(expression)), append(v.expression, expression...),
	)
}
