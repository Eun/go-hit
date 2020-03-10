package hit

import (
	"github.com/Eun/go-hit/internal"
)

type IClearSendBody interface {
	IStep
	JSON(...interface{}) IStep
	Interface(...interface{}) IStep
}

type clearSendBody struct {
	cleanPath clearPath
}

func newClearSendBody(clearPath clearPath, params []interface{}) IClearSendBody {
	if _, ok := internal.GetLastArgument(params); ok {
		return finalClearSendBody{&hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: clearPath,
			Exec: func(hit Hit) error {
				removeSteps(hit, clearPath)
				return nil
			},
		}}
	}

	return &clearSendBody{
		cleanPath: clearPath,
	}
}

func (*clearSendBody) when() StepTime {
	return SendStep
}

// exec contains the logic for Send().Body(...)
func (body *clearSendBody) exec(hit Hit) error {
	removeSteps(hit, body.cleanPath)
	return nil
}

func (body *clearSendBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *clearSendBody) JSON(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("JSON", data))
}

func (body *clearSendBody) Interface(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Interface", data))
}

type finalClearSendBody struct {
	IStep
}

func (finalClearSendBody) JSON(...interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}

func (finalClearSendBody) Interface(...interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}
