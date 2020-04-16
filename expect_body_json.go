package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal/minitest"
	"github.com/Eun/go-hit/internal/misc"
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
	trace      *errortrace.ErrorTrace
}

func newExpectBodyJSON(expectBody IExpectBody, cleanPath clearPath, params []interface{}) IExpectBodyJSON {
	jsn := &expectBodyJSON{
		expectBody: expectBody,
		cleanPath:  cleanPath,
		trace:      ett.Prepare(),
	}

	if param, ok := misc.GetLastArgument(params); ok {
		return &finalExpectBodyJSON{
			&hitStep{
				Trace:     jsn.trace,
				When:      ExpectStep,
				ClearPath: jsn.cleanPath,
				Exec:      jsn.Equal("", param).exec,
			},
			"only usable with Expect().Body().JSON() not with Expect().Body().JSON(value)",
		}
	}
	return jsn
}

func (jsn *expectBodyJSON) exec(hit Hit) error {
	return jsn.trace.Format(hit.Description(), "unable to run Expect().Body().JSON() without an argument or without a chain. Please use Expect().Body().JSON(something) or Expect().Body().JSON().Something")
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
		ClearPath: jsn.clearPath().Push("Equal", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				return nil
			}

			if v == nil || data == nil {
				// will fail
				return minitest.Error.Equal(data, v)
			}

			compareData, err := makeCompareable(v, data)
			if err != nil {
				return nil
			}
			return minitest.Error.Equal(data, compareData)
		},
	}
}

func (jsn *expectBodyJSON) NotEqual(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.clearPath().Push("NotEqual", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			v := hit.Response().body.JSON().Get(expression)
			if v == nil && data == nil {
				return minitest.Error.Errorf("should not be %s", minitest.PrintValue(v))
			}

			if v == nil || data == nil {
				return minitest.Error.NotEqual(data, v)
			}

			compareData, err := makeCompareable(v, data)
			if err != nil {
				return err
			}
			return minitest.Error.NotEqual(data, compareData)
		},
	}
}

func (jsn *expectBodyJSON) Contains(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.clearPath().Push("Contains", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			return minitest.Error.Contains(hit.Response().body.JSON().Get(expression), data)
		},
	}
}

func (jsn *expectBodyJSON) NotContains(expression string, data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: jsn.clearPath().Push("NotContains", []interface{}{expression, data}),
		Exec: func(hit Hit) error {
			return minitest.Error.NotContains(hit.Response().body.JSON().Get(expression), data)
		},
	}
}

type finalExpectBodyJSON struct {
	IStep
	message string
}

func (jsn *finalExpectBodyJSON) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(jsn.message)
		},
	}
}

func (jsn *finalExpectBodyJSON) Equal(string, interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalExpectBodyJSON) NotEqual(string, interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalExpectBodyJSON) Contains(string, interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalExpectBodyJSON) NotContains(string, interface{}) IStep {
	return jsn.fail()
}
