package hit

import (
	"net/http"

	"strings"

	"github.com/Eun/go-hit/errortrace"
)

type ISend interface {
	// Body sets the request body to the specified value.
	//
	// If you omit the argument you can fine tune the send value.
	//
	// Usage:
	//     Send().Body("Hello World")
	//     Send().Body().Interface("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body("Hello World"),
	//     )
	Body(value ...interface{}) ISendBody

	// Header sets the specified request header to the specified value
	//
	// Usage:
	//     Send().Header("Content-Type", "application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Header("Content-Type", "application/json"),
	//     )
	Header(name string, value interface{}) IStep

	// Trailer sets the specified request trailer to the specified value
	//
	// Usage:
	//     Send().Trailer("Content-Type", "application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Trailer("Content-Type", "application/json"),
	//     )
	Trailer(name string, value interface{}) IStep

	// Custom can be used to send a custom behaviour.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Custom(func(hit Hit) {
	//                hit.Request().Body().SetString("Hello World")
	//         }),
	//     )
	Custom(fn Callback) IStep
}

type send struct {
	body      ISendBody
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newSend(clearPath clearPath) ISend {
	return &send{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}
}

func (snd *send) Body(value ...interface{}) ISendBody {
	if snd.body == nil {
		snd.body = newSendBody(snd.cleanPath.Push("Body", value), value)
	}
	return snd.body
}

func (snd *send) Header(name string, value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Header", []interface{}{name, value}),
		Exec: func(hit Hit) error {
			var s string
			if err := converter.Convert(value, &s); err != nil {
				return err
			}
			if strings.EqualFold(name, "host") {
				hit.Request().Host = s
			}
			hit.Request().Header.Set(name, s)
			return nil
		},
	}
}

func (snd *send) Trailer(name string, value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Trailer", []interface{}{name, value}),
		Exec: func(hit Hit) error {
			var s string
			if err := converter.Convert(value, &s); err != nil {
				return err
			}
			if hit.Request().Trailer == nil {
				hit.Request().Trailer = make(http.Header)
			}
			hit.Request().Trailer.Set(name, s)
			return nil
		},
	}
}

func (snd *send) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: snd.cleanPath.Push("Custom", []interface{}{fn}),
		Exec: func(hit Hit) error {
			fn(hit)
			return nil
		},
	}
}
