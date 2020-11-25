package hit

import (
	"net/http"
)

// IRequest provides methods to set request url parameters.
type IRequest interface {
	// Set sets the request to the specified value.
	Set(request *http.Request) IStep

	// URL provides methods to set request url parameters.
	URL() IRequestURL

	// Headers sets the specified request header to the specified value(s).
	//
	// Usage:
	//     Request().Headers("Content-Type").Add("application/json")
	//     Request().Headers("Set-Cookie").Add("user=joe")
	//     Request().Headers("Set-Cookie").Add("id=joe")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Request().Headers("Content-Type").Add("application/json"),
	//     )
	Headers(name string) ISendHeaders

	// Trailers sets the specified request trailer to the specified value(s).
	//
	// Usage:
	//     Request().Trailers("Content-Type").Add("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Request().Trailers("Content-Type").Add("application/json"),
	//     )
	Trailers(name string) ISendHeaders

	// Host sets the request host to the specified value.
	//
	// Usage:
	//     Request().Host("example.com")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Request().Host("example.com"),
	//     )
	Host(name string) IStep

	// Method sets the request method to the specified value.
	//
	// Usage:
	//     Request().Method(http.MethodPost)
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodPost),
	//         Request().URL().Scheme("https"),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//     )
	Method(method string) IStep
}

type request struct {
	cleanPath callPath
}

func newRequest(clearPath callPath) IRequest {
	return &request{
		cleanPath: clearPath,
	}
}

func (req *request) Set(r *http.Request) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: req.cleanPath.Push("Set", []interface{}{r}),
		Exec: func(hit *hitImpl) error {
			hit.request = newHTTPRequest(hit, r)
			return nil
		},
	}
}

func (req *request) URL() IRequestURL {
	return newRequestURL(req.cleanPath.Push("URL", nil))
}

func (req *request) Headers(name string) ISendHeaders {
	return newSendHeaders(req.cleanPath.Push("Headers", []interface{}{name}), func(hit *hitImpl) http.Header {
		return hit.request.Header
	}, name)
}

func (req *request) Trailers(name string) ISendHeaders {
	return newSendHeaders(req.cleanPath.Push("Trailers", []interface{}{name}), func(hit *hitImpl) http.Header {
		if hit.request.Trailer == nil {
			hit.request.Trailer = make(http.Header)
		}
		return hit.request.Trailer
	}, name)
}

func (req *request) Host(name string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: req.cleanPath.Push("Host", []interface{}{name}),
		Exec: func(hit *hitImpl) error {
			hit.request.Host = name
			return nil
		},
	}
}

func (req *request) Method(method string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: req.cleanPath.Push("Method", []interface{}{method}),
		Exec: func(hit *hitImpl) error {
			hit.request.Method = method
			return nil
		},
	}
}
