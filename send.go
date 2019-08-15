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
	Body(data ...interface{}) *sendBody
	Custom(f Callback) IStep
	JSON(data interface{}) IStep
	Headers() *sendHeaders
	Header(name string) *sendSpecificHeader
	Clear() IStep
	Interface(data interface{}) IStep
}

type defaultSend struct {
	body *sendBody
	call Callback
}

func (exp *defaultSend) when() StepTime {
	return SendStep
}

func (exp *defaultSend) exec(hit Hit) {
	exp.call(hit)
}

func newSend() *defaultSend {
	snd := &defaultSend{}
	snd.body = newSendBody(snd)
	return snd
}

func (snd *defaultSend) Body(data ...interface{}) *sendBody {
	if arg, ok := getLastArgument(data); ok {
		snd.Interface(arg)
	}
	return snd.body
}

// Custom can be used to send a custom behaviour
func (snd *defaultSend) Custom(f Callback) IStep {
	snd.call = f
	return snd
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *defaultSend) JSON(data interface{}) IStep {
	return snd.body.JSON(data)
}

// Headers sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Headers().Set("Content-Type", "application/json")
//           Send().Headers().Delete("Content-Type")
func (snd *defaultSend) Headers() *sendHeaders {
	return newSendHeaders(snd)
}

// Header sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Header("Content-Type").Set("application/json")
//           Send().Header("Content-Type").Delete()
func (snd *defaultSend) Header(name string) *sendSpecificHeader {
	return newSendSpecificHeader(snd, name)
}

func (snd *defaultSend) Clear() IStep {
	return Custom(SendStep|CleanStep, func(hit Hit) {})
}

func (snd *defaultSend) Interface(data interface{}) IStep {
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

type dummySend struct {
	IStep
}

func (d dummySend) Body(data ...interface{}) *sendBody {
	panic("implement me")
}

func (d dummySend) Custom(f Callback) IStep {
	panic("implement me")
}

func (d dummySend) JSON(data interface{}) IStep {
	panic("implement me")
}

func (d dummySend) Headers() *sendHeaders {
	panic("implement me")
}

func (d dummySend) Header(name string) *sendSpecificHeader {
	panic("implement me")
}

func (d dummySend) Clear() IStep {
	panic("implement me")
}

func (d dummySend) Interface(data interface{}) IStep {
	panic("implement me")
}
