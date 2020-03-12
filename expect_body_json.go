package hit

import (
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
	"golang.org/x/xerrors"
)

// IExpectBodyJSON provides assertions on the http response json body
type IExpectBodyJSON interface {
	IStep
	// Equal expects the json body to be equal to the specified value.
	//
	// The first argument can be used to narrow down the compare path
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Equal("", map[string]interface{}{"ID": 10, "Name": "Joe", "Roles": []string{"Admin", "User"}})
	//     Expect().Body().JSON().Equal("ID", 10)
	//     Expect().Body().JSON().Equal("Name", "Joe")
	//     Expect().Body().JSON().Equal("Roles", []string{"Admin", "User"}),
	//     Expect().Body().JSON().Equal("Roles.0", "Admin"),
	//
	// Example:
	//     // given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().JSON().Equal("Name", "Joe"),
	//         Expect().Body().JSON().Equal("Roles", []string{"Admin", "User"}),
	//         Expect().Body().JSON().Equal("Roles.0", "Admin"),
	//     )
	Equal(expression string, data interface{}) IStep

	// NotEqual expects the json body to be equal to the specified value.
	//
	// The first argument can be used to narrow down the compare path
	//
	// see Equal() for usage and examples
	NotEqual(expression string, data interface{}) IStep

	// Contains expects the json body to be equal to the specified value.
	//
	// The first argument can be used to narrow down the compare path
	//
	// given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     Expect().Body().JSON().Contains("", "ID")
	//     Expect().Body().JSON().Contains("Name", "J")
	//     Expect().Body().JSON().Contains("Roles", "Admin"),
	//
	// Example:
	//     // given the following response: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().JSON().Contains("", "ID"),
	//     )
	Contains(expression string, data interface{}) IStep

	// NotContains expects the json body to be equal to the specified value.
	//
	// The first argument can be used to narrow down the compare path
	//
	// see Contains() for usage and examples
	NotContains(expression string, data interface{}) IStep
}

type expectBodyJSON struct {
	expectBody IExpectBody
	cleanPath  clearPath
}

func newExpectBodyJSON(expectBody IExpectBody, cleanPath clearPath, params []interface{}) IExpectBodyJSON {
	jsn := &expectBodyJSON{
		expectBody: expectBody,
		cleanPath:  cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		return finalExpectBodyJSON{&hitStep{
			Trace:     ett.Prepare(),
			When:      ExpectStep,
			ClearPath: jsn.cleanPath,
			Exec:      jsn.Equal("", param).exec,
		}}
	}
	return jsn
}

func (*expectBodyJSON) exec(Hit) error {
	return xerrors.New("unable to run Expect().Body().JSON() without an argument or without a chain. Please use Expect().Body().JSON(something) or Expect().Body().JSON().Something")
}

func (*expectBodyJSON) when() StepTime {
	return ExpectStep
}

func (jsn *expectBodyJSON) clearPath() clearPath {
	return jsn.cleanPath
}

func (jsn *expectBodyJSON) Equal(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.cleanPath.Push("Equal", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				return nil
			}

			if v == nil || data == nil {
				// will fail
				minitest.Equal(data, v)
			}

			compareData, err := makeCompareable(v, data)
			if err != nil {
				return nil
			}
			minitest.Equal(data, compareData)
			return nil
		},
	}
}

func (jsn *expectBodyJSON) NotEqual(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.cleanPath.Push("NotEqual", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				minitest.Errorf("should not be %s", minitest.PrintValue(v))
			}

			if v == nil || data == nil {
				minitest.NotEqual(data, v)
			}

			compareData, err := makeCompareable(v, data)
			if err != nil {
				return err
			}
			minitest.NotEqual(data, compareData)
			return nil
		},
	}
}

func (jsn *expectBodyJSON) Contains(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.cleanPath.Push("Contains", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				return nil
			}

			if !internal.Contains(v, data) {
				minitest.Errorf("%s does not contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
			}
			return nil
		},
	}
}

func (jsn *expectBodyJSON) NotContains(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.cleanPath.Push("NotContains", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				minitest.Errorf("%s does contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
			}

			if internal.Contains(v, data) {
				minitest.Errorf("%s does contain %s", minitest.PrintValue(v), minitest.PrintValue(data))
			}
			return nil
		},
	}
}

type finalExpectBodyJSON struct {
	IStep
}

func (finalExpectBodyJSON) Equal(string, interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) NotEqual(string, interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) Contains(string, interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}

func (finalExpectBodyJSON) NotContains(string, interface{}) IStep {
	panic("only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)")
}
