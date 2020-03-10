package hit

import (
	"github.com/Eun/go-hit/internal"
)

type IClearSend interface {
	IStep
	Body(...interface{}) IClearSendBody
	Interface(...interface{}) IStep
	Custom(...Callback) IStep
	JSON(...interface{}) IStep
	Header(...string) IStep
}

type clearSend struct {
	body      ISendBody
	cleanPath clearPath
}

func newClearSend(clearPath clearPath, params []interface{}) IClearSend {
	if _, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalClearSend{&hitStep{
			Trace:     ett.Prepare(),
			When:      CleanStep,
			ClearPath: clearPath,
			Exec: func(hit Hit) error {
				removeSteps(hit, clearPath)
				return nil
			},
		}}
	}
	return &clearSend{
		cleanPath: clearPath,
	}
}

func (*clearSend) when() StepTime {
	return CleanStep
}

// exec contains the logic for Send(...)
func (snd *clearSend) exec(hit Hit) error {
	removeSteps(hit, snd.cleanPath)
	return nil
}

func (snd *clearSend) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *clearSend) Body(data ...interface{}) IClearSendBody {
	return newClearSendBody(snd.cleanPath.Push("Body", data), data)
}

func (snd *clearSend) Interface(data ...interface{}) IStep {
	return removeStep(snd.cleanPath.Push("Interface", data))
}

// custom can be used to send a custom behaviour
func (snd *clearSend) Custom(f ...Callback) IStep {
	args := make([]interface{}, len(f))
	for i := range f {
		args[i] = f[i]
	}
	return removeStep(snd.cleanPath.Push("Custom", args))
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *clearSend) JSON(data ...interface{}) IStep {
	return removeStep(snd.cleanPath.Push("JSON", data))
}

// Header clears the specified header
// Examples:
//           Clear().Send().Header("Content-Type")
func (snd *clearSend) Header(name ...string) IStep {
	args := make([]interface{}, len(name))
	for i := range name {
		args[i] = name[i]
	}
	return removeStep(snd.cleanPath.Push("Header", args))
}

type finalClearSend struct {
	IStep
}

func (finalClearSend) Body(...interface{}) IClearSendBody {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Custom(...Callback) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) JSON(...interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Header(...string) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Interface(...interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}
