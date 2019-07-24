package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
)

type Send interface {
	Hit
	Body(data ...interface{}) *sendBody
	Custom(f Callback) Hit
	JSON(data interface{}) Hit
	Headers() *sendHeaders
	Clear() Hit
	Interface(data interface{}) Hit
	CollectedSends() []Callback
}

type defaultSend struct {
	Hit
	sendCalls []Callback
	headers   *sendHeaders
	body      *sendBody
}

func newSend(hit Hit) *defaultSend {
	snd := &defaultSend{
		Hit: hit,
	}
	snd.headers = newSendHeaders(snd)
	snd.body = newSendBody(snd)
	return snd
}

func (snd *defaultSend) Body(data ...interface{}) *sendBody {
	if arg := getLastArgument(data); arg != nil {
		snd.Interface(arg)
	}
	return snd.body
}

// Custom can be used to send a custom behaviour
func (snd *defaultSend) Custom(f Callback) Hit {
	switch snd.State() {
	case Done:
		errortrace.Panic.Errorf(snd.Hit.T(), "request already fired")
	case Working:
		f(snd.Hit)
	default:
		snd.sendCalls = append(snd.sendCalls, f)
	}
	return snd.Hit
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *defaultSend) JSON(data interface{}) Hit {
	return snd.body.JSON(data)
}

func (snd *defaultSend) Headers() *sendHeaders {
	switch snd.State() {
	case Done:
		errortrace.Panic.Errorf(snd.Hit.T(), "request already fired")
		return nil
	default: // Ready, Working
		return snd.headers
	}
}

func (snd *defaultSend) Clear() Hit {
	snd.sendCalls = nil
	return snd.Hit
}

func (snd *defaultSend) Interface(data interface{}) Hit {
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

func (snd *defaultSend) CollectedSends() []Callback {
	return snd.sendCalls
}

func (snd *defaultSend) copy(toHit Hit) *defaultSend {
	n := &defaultSend{
		Hit: toHit,
	}
	n.headers = newSendHeaders(n)
	n.body = newSendBody(n)
	// copy the send calls
	n.sendCalls = make([]Callback, len(snd.sendCalls))
	for i, v := range snd.sendCalls {
		n.sendCalls[i] = v
	}
	return n
}
