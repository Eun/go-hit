package hit

import (
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type ISend interface {
	IStep
	Body(data ...interface{}) ISendBody
	Custom(f Callback) IStep
	JSON(data interface{}) IStep
	// Header sets the specified header
	// Examples:
	//           Send().Header("Content-Type").Set("application/json")
	//           Send().Header("Content-Type").Delete()
	Header(name, value string) IStep
	Interface(data interface{}) IStep
}

type send struct {
	body      ISendBody
	cleanPath clearPath
}

func newSend(clearPath clearPath, params []interface{}) ISend {
	snd := &send{
		cleanPath: clearPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalSend{&hitStep{
			Trace:     ett.Prepare(),
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

// exec contains the logic for Send(...)
func (*send) exec(hit Hit) error {
	return xerrors.New("unsupported")
}

func (snd *send) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *send) Body(data ...interface{}) ISendBody {
	if snd.body == nil {
		snd.body = newSendBody(snd.cleanPath.Push("Body", data), data)
	}
	return snd.body
}

// custom can be used to send a custom behaviour
func (snd *send) Custom(f Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Custom", []interface{}{f}),
		Exec: func(hit Hit) error {
			f(hit)
			return nil
		},
	}
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *send) JSON(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("JSON", []interface{}{data}),
		Exec:      snd.Body().JSON(data).exec,
	}
}

// Header sets the specified header
// Examples:
//           Send().Header("Content-Type", "application/json")
func (snd *send) Header(name, value string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Header", []interface{}{name, value}),
		Exec: func(hit Hit) error {
			hit.Request().Header.Set(name, value)
			return nil
		},
	}
}

func (snd *send) Interface(data interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Interface", []interface{}{data}),
		Exec:      snd.Body(data).exec,
	}
}

type finalSend struct {
	IStep
}

func (d finalSend) Body(...interface{}) ISendBody {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Custom(Callback) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) JSON(interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Header(string, string) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Interface(interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}
