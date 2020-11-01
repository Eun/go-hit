package hit

import (
	"reflect"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/internal/converter"
)

// IStore provides a store functionality for the Request and Response.
type IStore interface {
	// Request stores the Request.
	//
	// Example:
	//     var body string
	//     var userName string
	//     var headers http.Header
	//     var contentType string
	//     MustDo(
	//         Post("https://example.com/json"),
	//         Send().Body().JSON(map[string]interface{}{"name": "Joe"}),
	//         Store().Request().Body().String().In(&body),               // store the body
	//         Store().Request().Body().JSON().JQ(".name").In(&userName), // parse body as json and store the data object into data variable
	//         Store().Request().Headers().In(&headers),                   // store all headers
	//         Store().Request().Headers("Content-Type").In(&contentType), // store the Content-Type header
	//     )
	Request() IStoreRequest

	// Response stores the Response.
	//
	// Example:
	//     var body string
	//     var userName string
	//     var headers http.Header
	//     var contentType string
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Store().Response().Body().String().In(&body),               // store the body
	//         Store().Response().Body().JSON().JQ(".name").In(&userName), // parse body as json and store the data object into data variable
	//         Store().Response().Headers().In(&headers),                   // store all headers
	//         Store().Response().Headers("Content-Type").In(&contentType), // store the Content-Type header
	//     )
	Response() IStoreResponse
}

type store struct {
}

func newStore() IStore {
	return &store{}
}

func (*store) Request() IStoreRequest {
	return newStoreRequest()
}

func (*store) Response() IStoreResponse {
	return newStoreResponse()
}

type storeFunc func(Hit, interface{}) error

// IStoreStep defines the In function for the Store() functionality.
type IStoreStep interface {
	// In can be used to to store the result into an existing variable
	//
	// Example:
	//     var body string
	//     MustDo(
	//         Get("https://example.com"),
	//         Store().Response().Body().String().In(&body),
	//     )
	In(interface{}) IStep
}

func newStoreStep(f storeFunc) IStoreStep {
	return &storeStep{
		f,
	}
}

type storeStep struct {
	f storeFunc
}

func (s *storeStep) In(v interface{}) IStep {
	return newStoreInStep(s.f, v)
}

func newStoreInStep(f storeFunc, v interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     AfterExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			return f(hit, v)
		},
	}
}

func storeStringSlice(in []string, out interface{}) error {
	if out == nil {
		return xerrors.New("destination type cannot be nil")
	}

	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return xerrors.New("destination must be a pointer")
	}
	if t.Elem().Kind() == reflect.Slice {
		return converter.Convert(in, out)
	}

	switch len(in) {
	case 0:
		return converter.Convert("", out)
	case 1:
		return converter.Convert(in[0], out)
	default:
		return xerrors.Errorf("could not put %#v into %T", in, out)
	}
}
