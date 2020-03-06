package hit

import (
	"errors"

	"github.com/Eun/go-hit/internal"
)

type ISend interface {
	IStep
	Body(data ...interface{}) ISendBody
	Custom(f Callback) IStep
	JSON(data interface{}) IStep
	Headers() ISendHeaders
	// Header sets the specified header
	// Examples:
	//           Send().Header("Content-Type").Set("application/json")
	//           Send().Header("Content-Type").Delete()
	Header(name string) ISendSpecificHeader
	Interface(data interface{}) IStep
}

type send struct {
	body      ISendBody
	cleanPath CleanPath
}

func newSend(cleanPath CleanPath, params []interface{}) ISend {
	snd := &send{
		cleanPath: cleanPath,
	}

	if param, ok := internal.GetLastArgument(params); ok {
		// default action is Interface()
		return finalSend{wrap(wrapStep{
			IStep:     snd.Interface(param),
			CleanPath: cleanPath,
		})}
	}
	return snd
}

func (*send) When() StepTime {
	return SendStep
}

// Exec contains the logic for Send(...)
func (*send) Exec(hit Hit) error {
	return errors.New("unsupported")
}

func (snd *send) CleanPath() CleanPath {
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
	return custom(Step{
		When:      SendStep,
		CleanPath: snd.cleanPath.Push("Custon", []interface{}{f}),
		Exec:      f,
	})
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *send) JSON(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: snd.cleanPath.Push("JSON", []interface{}{data}),
		Exec: func(hit Hit) {
			panic("TODO")
			// hit.Send().Body().JSON(data)
		},
	})
}

// Headers sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Headers().Set("Content-Type", "application/json")
//           Send().Headers().Delete("Content-Type")
func (snd *send) Headers() ISendHeaders {
	return newSendHeaders(snd.cleanPath.Push("Headers", nil))
}

// Header sets the specified header
// Examples:
//           Send().Header("Content-Type").Set("application/json")
//           Send().Header("Content-Type").Delete()
func (snd *send) Header(name string) ISendSpecificHeader {
	return newSendSpecificHeader(snd.cleanPath.Push("Header", []interface{}{name}), name)
}

func (snd *send) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return custom(Step{
			When:      SendStep,
			CleanPath: snd.cleanPath.Push("Interface", []interface{}{data}),
			Exec:      x,
		})
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return custom(Step{
				When:      SendStep,
				CleanPath: snd.cleanPath.Push("Interface", []interface{}{data}),
				Exec: func(hit Hit) {
					internal.CallGenericFunc(f)
				},
			})
		}
		return wrap(wrapStep{
			IStep:     snd.Body(data),
			CleanPath: snd.cleanPath.Push("Interface", []interface{}{data}),
		})
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

func (d finalSend) Headers() ISendHeaders {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Header(string) ISendSpecificHeader {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Clear() IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Interface(interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}
