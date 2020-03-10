package hit

import (
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type ISendBody interface {
	IStep
	JSON(data interface{}) IStep
	Interface(data interface{}) IStep
}

type sendBody struct {
	cleanPath clearPath
}

func newSendBody(clearPath clearPath, params []interface{}) ISendBody {
	snd := &sendBody{
		cleanPath: clearPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		return finalSendBody{&hitStep{
			Trace:     ett.Prepare(),
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

// exec contains the logic for Send().Body(...)
func (*sendBody) exec(Hit) error {
	return xerrors.New("unsupported")
}

func (body *sendBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *sendBody) JSON(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: body.cleanPath.Push("JSON", []interface{}{data}),
		Exec: func(hit Hit) error {
			hit.Request().Body().JSON().Set(data)
			return nil
		},
	}
}

func (body *sendBody) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec: func(hit Hit) error {
				x(hit)
				return nil
			},
		}
	case func(e Hit) error:
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec:      x,
		}
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return &hitStep{
				Trace:     ett.Prepare(),
				When:      SendStep,
				ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
				Exec: func(hit Hit) error {
					internal.CallGenericFunc(f)
					return nil
				},
			}
		}
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.cleanPath.Push("Interface", []interface{}{data}),
			Exec: func(hit Hit) error {
				hit.Request().Body().Set(data)
				return nil
			},
		}
	}
}

type finalSendBody struct {
	IStep
}

func (d finalSendBody) JSON(interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}

func (d finalSendBody) Interface(interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}
