package hit

import (
	"github.com/Eun/go-hit/internal"
)

// import (
// 	"github.com/Eun/go-hit/errortrace"
// 	"github.com/Eun/go-hit/internal"
// )
//
type ISend interface {
	IStep
	Body(data ...interface{}) ISendBody
	Custom(f Callback) IStep
	JSON(data interface{}) IStep
	Headers() ISendHeaders
	Header(name string) ISendSpecificHeader
	Clear() IStep
	Interface(data interface{}) IStep
}

type send struct {
	body ISendBody
	call Callback
}

func newSend() ISend {
	snd := &send{}
	snd.body = newSendBody(snd)
	return snd
}

func (exp *send) when() StepTime {
	return SendStep
}

func (exp *send) exec(hit Hit) {
	if exp.call != nil {
		exp.call(hit)
	}
}

func (snd *send) Body(data ...interface{}) ISendBody {
	if arg, ok := getLastArgument(data); ok {
		return finalSendBody{snd.Interface(arg)}
	}
	return snd.body
}

// Custom can be used to send a custom behaviour
func (snd *send) Custom(f Callback) IStep {
	snd.call = f
	return snd
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *send) JSON(data interface{}) IStep {
	return snd.body.JSON(data)
}

// Headers sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Headers().Set("Content-Type", "application/json")
//           Send().Headers().Delete("Content-Type")
func (snd *send) Headers() ISendHeaders {
	return newSendHeaders(snd)
}

// Header sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Header("Content-Type").Set("application/json")
//           Send().Header("Content-Type").Delete()
func (snd *send) Header(name string) ISendSpecificHeader {
	return newSendSpecificHeader(snd, name)
}

func (snd *send) Clear() IStep {
	return Custom(SendStep|CleanStep, func(hit Hit) {})
}

func (snd *send) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return snd.Custom(x)
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return snd.Custom(func(hit Hit) {
				internal.CallGenericFunc(f)
			})
		}
		return snd.body.Interface(x)
	}
}

type finalSend struct {
	IStep
}

func (d finalSend) Body(data ...interface{}) ISendBody {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Custom(f Callback) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) JSON(data interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Headers() ISendHeaders {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Header(name string) ISendSpecificHeader {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Clear() IStep {
	panic("only usable with Send() not with Send(value)")
}

func (d finalSend) Interface(data interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}
