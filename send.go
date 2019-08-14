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
	Headers(name ...string) *sendHeaders
	Clear() IStep
	Interface(data interface{}) IStep
}

type defaultSend struct {
	headers *sendHeaders
	body    *sendBody
	call    Callback
}

func (exp *defaultSend) when() StepTime {
	return SendStep
}

func (exp *defaultSend) exec(hit Hit) {
	exp.call(hit)
}

func newSend() *defaultSend {
	snd := &defaultSend{}
	snd.headers = newSendHeaders(snd)
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

func (snd *defaultSend) Headers(name ...string) *sendHeaders {
	return snd.headers
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

func (d dummySend) Headers(name ...string) *sendHeaders {
	panic("implement me")
}

func (d dummySend) Clear() IStep {
	panic("implement me")
}

func (d dummySend) Interface(data interface{}) IStep {
	panic("implement me")
}
