package hit

import (
	"net/http"
)

// ISend provides methods to set request data, such as body or headers.
type ISend interface {
	// Body sets the request body to the specified value.
	//
	// Usage:
	//     Send().Body().String("Hello World")
	//     Send().Body().JSON("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().String("Hello World"),
	//     )
	Body() ISendBody

	// Headers sets the specified request header to the specified value(s).
	//
	// Usage:
	//     Send().Headers("Content-Type").Add("application/json")
	//     Send().Headers("Set-Cookie").Add("user=joe")
	//     Send().Headers("Set-Cookie").Add("id=joe")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Headers("Content-Type").Add("application/json"),
	//     )
	Headers(name string) ISendHeaders

	// Trailers sets the specified request trailer to the specified value(s).
	//
	// Usage:
	//     Send().Trailers("Content-Type").Add("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Trailers("Content-Type").Add("application/json"),
	//     )
	Trailers(name string) ISendHeaders

	// Custom can be used to send a custom behavior.
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
	cleanPath callPath
}

func newSend(clearPath callPath) ISend {
	return &send{
		cleanPath: clearPath,
	}
}

func (snd *send) Body() ISendBody {
	if snd.body == nil {
		snd.body = newSendBody(snd.cleanPath.Push("Body", nil))
	}
	return snd.body
}

func (snd *send) Headers(name string) ISendHeaders {
	return newSendHeaders(snd.cleanPath.Push("Headers", []interface{}{name}), func(hit Hit) http.Header {
		return hit.Request().Header
	}, name)
}
func (snd *send) Trailers(name string) ISendHeaders {
	return newSendHeaders(snd.cleanPath.Push("Trailers", []interface{}{name}), func(hit Hit) http.Header {
		if hit.Request().Trailer == nil {
			hit.Request().Trailer = make(http.Header)
		}
		return hit.Request().Trailer
	}, name)
}

func (snd *send) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: snd.cleanPath.Push("Custom", []interface{}{fn}),
		Exec: func(hit *hitImpl) error {
			return fn(hit)
		},
	}
}
