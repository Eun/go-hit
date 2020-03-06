package hit

import (
	"errors"

	"github.com/Eun/go-hit/errortrace"
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
	hit       Hit
	body      ISendBody
	call      Callback
	et        *errortrace.ErrorTrace
	param     []interface{}
	cleanPath CleanPath
}

func newSend(hit Hit, cleanPath CleanPath, param []interface{}) ISend {
	return &send{
		hit:       hit,
		cleanPath: cleanPath,
		param:     param,
	}
}

func (snd *send) When() StepTime {
	return SendStep
}

// Exec contains the logic for Send(...)
func (snd *send) Exec(hit Hit) error {
	param, ok := internal.GetLastArgument(snd.param)
	if !ok {
		return errors.New("invalid argument")
	}
	return snd.Interface(param).Exec(hit)
}

func (snd *send) CleanPath() CleanPath {
	return snd.cleanPath
}

func (snd *send) Body(data ...interface{}) ISendBody {
	if snd.body == nil {
		snd.body = newSendBody(snd, snd.hit, snd.cleanPath.Push("Body", data), data)
	}
	return snd.body
}

// custom can be used to send a custom behaviour
func (snd *send) Custom(f Callback) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: snd.cleanPath.Push("custom"),
		Instance:  snd.hit,
		Exec:      f,
	})
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *send) JSON(data interface{}) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: snd.cleanPath.Push("JSON"),
		Instance:  snd.hit,
		Exec: func(hit Hit) {
			hit.Send().Body().JSON(data)
		},
	})
}

// Headers sets the specified header, omit the parameter to get all headers
// Examples:
//           Send().Headers().Set("Content-Type", "application/json")
//           Send().Headers().Delete("Content-Type")
func (snd *send) Headers() ISendHeaders {
	return newSendHeaders(snd, snd.hit, snd.cleanPath.Push("Headers"))
}

// Header sets the specified header
// Examples:
//           Send().Header("Content-Type").Set("application/json")
//           Send().Header("Content-Type").Delete()
func (snd *send) Header(name string) ISendSpecificHeader {
	return newSendSpecificHeader(snd, snd.hit, snd.cleanPath.Push("Header", name), name)
}

func (snd *send) Interface(data interface{}) IStep {
	switch x := data.(type) {
	case func(e Hit):
		return custom(Step{
			When:      SendStep,
			CleanPath: snd.cleanPath.Push("Interface"),
			Instance:  snd.hit,
			Exec:      x,
		})
	default:
		if f := internal.GetGenericFunc(data); f.IsValid() {
			return custom(Step{
				When:      SendStep,
				CleanPath: snd.cleanPath.Push("Interface"),
				Instance:  snd.hit,
				Exec: func(hit Hit) {
					internal.CallGenericFunc(f)
				},
			})
		}
		return custom(Step{
			When:      SendStep,
			CleanPath: snd.cleanPath.Push("Body").Push("Interface"),
			Instance:  snd.hit,
			Exec: func(hit Hit) {
				hit.Send().Body().Interface(x)
			},
		})
	}
}
