package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type ISendBody interface {
	IStep
	// JSON sets the request body to the specified json value.
	//
	// Usage:
	//     Send().Body().JSON(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//     )
	JSON(value interface{}) IStep
	// Interface sets the request body to the specified json value.
	//
	// Usage:
	//     Send().Body().Interface("Hello World")
	//     Send().Body().Interface(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().Interface("Hello World"),
	//     )
	Interface(value interface{}) IStep
}

type sendBody struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newSendBody(clearPath clearPath, params []interface{}) ISendBody {
	snd := &sendBody{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}

	if param, ok := internal.GetLastArgument(params); ok {
		return finalSendBody{&hitStep{
			Trace:     snd.trace,
			When:      SendStep,
			ClearPath: clearPath,
			Exec:      snd.Interface(param).exec,
		}}
	}

	return snd
}

func (*sendBody) when() StepTime {
	return SendStep
}

func (body *sendBody) exec(hit Hit) error {
	return body.trace.Format(hit.Description(), "unable to run Send().Body() without an argument or without a chain. Please use Send().Body(something) or Send().Body().Something")
}

func (body *sendBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *sendBody) JSON(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: body.cleanPath.Push("JSON", []interface{}{value}),
		Exec: func(hit Hit) error {
			hit.Request().Body().JSON().Set(value)
			return nil
		},
	}
}

func (body *sendBody) Interface(value interface{}) IStep {
	switch x := value.(type) {
	case func(e Hit):
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec: func(hit Hit) error {
				x(hit)
				return nil
			},
		}
	case func(e Hit) error:
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec:      x,
		}
	default:
		if f := internal.GetGenericFunc(value); f.IsValid() {
			return &hitStep{
				Trace:     ett.Prepare(),
				When:      SendStep,
				ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
				Exec: func(hit Hit) error {
					internal.CallGenericFunc(f)
					return nil
				},
			}
		}
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{value}),
			Exec: func(hit Hit) error {
				hit.Request().Body().Set(value)
				return nil
			},
		}
	}
}

type finalSendBody struct {
	IStep
}

func (d finalSendBody) JSON(interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send().Body() not with Send().Body(value)")
		},
	}
}

func (d finalSendBody) Interface(interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Send().Body() not with Send().Body(value)")
		},
	}
}
