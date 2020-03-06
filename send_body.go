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
	send      ISend
	hit       Hit
	cleanPath CleanPath
	params    []interface{}
}

func newSendBody(send ISend, hit Hit, cleanPath CleanPath, params []interface{}) ISendBody {
	return &sendBody{
		send:      send,
		hit:       hit,
		cleanPath: cleanPath,
		params:    params,
	}
}

func (body *sendBody) When() StepTime {
	return body.send.When()
}

// Exec contains the logic for Send().Body(...)
func (body *sendBody) Exec(hit Hit) error {
	param, ok := internal.GetLastArgument(body.params)
	if !ok {
		return errors.New("invalid argument")
	}
	return body.Interface(param).Exec(hit)
}

func (body *sendBody) CleanPath() CleanPath {
	return body.cleanPath
}

func (body *sendBody) JSON(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: body.cleanPath.Push("JSON"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			hit.Request().Body().JSON().Set(data)
		},
	})
}

func (body *sendBody) Interface(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: body.cleanPath.Push("Interface"),
		Instance:  body.hit,
		Exec: func(hit Hit) {
			hit.Request().Body().Set(data)
		},
	})
}
