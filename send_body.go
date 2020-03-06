package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
)

type ISendBody interface {
	IStep
	JSON(data interface{}) IStep
	Interface(data interface{}) IStep
}

type sendBody struct {
	cleanPath CleanPath
}

func newSendBody(cleanPath CleanPath, params []interface{}) ISendBody {
	snd := &sendBody{
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		return finalSendBody{wrap(wrapStep{
			IStep:     snd.Interface(param),
			CleanPath: cleanPath,
		})}
	}

	return snd
}

func (*sendBody) When() StepTime {
	return SendStep
}

// Exec contains the logic for Send().Body(...)
func (*sendBody) Exec(Hit) error {
	return errors.New("unsupported")
}

func (body *sendBody) CleanPath() CleanPath {
	return body.cleanPath
}

func (body *sendBody) JSON(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: body.cleanPath.Push("JSON", []interface{}{data}),
		Exec: func(hit Hit) {
			hit.Request().Body().JSON().Set(data)
		},
	})
}

func (body *sendBody) Interface(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: body.cleanPath.Push("Interface", []interface{}{data}),
		Exec: func(hit Hit) {
			hit.Request().Body().Set(data)
		},
	})
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
