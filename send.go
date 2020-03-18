package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type ISend interface {
	IStep
	// Body sets the request body to the specified value.
	//
	// If you omit the argument you can fine tune the send value.
	//
	// Usage:
	//     Send().Body("Hello World")
	//     Send().Body().Interface("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body("Hello World"),
	//     )
	Body(value ...interface{}) ISendBody

	// Interface sets the request body to the specified value.
	//
	// Usage:
	//     Send().Interface("Hello World")
	//     Send().Interface(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Interface("Hello World"),
	//     )
	Interface(value interface{}) IStep

	// JSON sets the request body to the specified json value.
	//
	// Usage:
	//     Send().JSON(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().JSON(map[string]interface{}{"Name": "Joe"}),
	//     )
	JSON(value interface{}) IStep

	// Header sets the specified request header to the specified value
	//
	// Usage:
	//     Send().Header("Content-Type", "application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Header("Content-Type", "application/json"),
	//     )
	Header(name string, value interface{}) IStep

	// Custom can be used to send a custom behaviour.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Custom(func(hit Hit) {
	//                hit.Request().Body().SetString("Hello World")
	//         }),
	//     )
	Custom(fn Callback) IStep
}

type send struct {
	body      ISendBody
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newSend(clearPath clearPath, params []interface{}) ISend {
	snd := &send{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalSend{&hitStep{
			Trace:     snd.trace,
			When:      SendStep,
			ClearPath: clearPath,
			Exec:      snd.Interface(param).exec,
		}}
	}
	return snd
}

func (*send) when() StepTime {
	return SendStep
}

func (snd *send) exec(hit Hit) error {
	return snd.trace.Format(hit.Description(), "unable to run Send() without an argument or without a chain. Please use Send(something) or Send().Something")
}

func (snd *send) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *send) Body(value ...interface{}) ISendBody {
	if snd.body == nil {
		snd.body = newSendBody(snd.cleanPath.Push("Body", value), value)
	}
	return snd.body
}

func (snd *send) Interface(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Interface", []interface{}{value}),
		Exec:      snd.Body(value).exec,
	}
}

func (snd *send) JSON(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("JSON", []interface{}{data}),
		Exec:      snd.Body().JSON(data).exec,
	}
}

func (snd *send) Header(name string, value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Header", []interface{}{name, value}),
		Exec: func(hit Hit) error {
			var s string
			if err := converter.Convert(value, &s); err != nil {
				return err
			}
			hit.Request().Header.Set(name, s)
			return nil
		},
	}
}

func (snd *send) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Custom", []interface{}{fn}),
		Exec: func(hit Hit) error {
			fn(hit)
			return nil
		},
	}
}

type finalSend struct {
	IStep
}

func (d finalSend) Body(...interface{}) ISendBody {
	return finalSendBody{&hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send() not with Send(value)")
		},
	}}
}

func (d finalSend) Custom(Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send() not with Send(value)")
		},
	}
}

func (d finalSend) JSON(interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send() not with Send(value)")
		},
	}
}

func (d finalSend) Header(string, interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send() not with Send(value)")
		},
	}
}

func (d finalSend) Interface(interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send() not with Send(value)")
		},
	}
}
