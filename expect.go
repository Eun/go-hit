package hit

import (
	"io"
	"io/ioutil"
)

// IExpect provides assertions on the http response.
type IExpect interface {
	// Custom can be used to expect a custom behavior.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Custom(func(hit Hit) error {
	//               if hit.Response().StatusCode != 200 {
	//                   return errors.New("Expected 200")
	//               }
	//               return nil
	//         }),
	//     )
	Custom(fn Callback) IStep

	// Body provides assertions for the body
	//
	// Usage:
	//     Expect().Body().String().Equal("Hello World")
	//     Expect().Body().JSON().Equal(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().String().Contains("Hello World"),
	//     )
	Body() IExpectBody

	// Headers provides assertions to one specific response headers.
	//
	// If you specify the argument you can directly assert a specific header
	//
	// Usage:
	//     Expect().Headers("Content-Type").NotEmpty()
	//     Expect().Headers("Content-Type").Equal("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").NotEmpty(),
	//         Expect().Headers("Content-Type").Equal("text/plain; charset=utf-8"),
	//     )
	Headers(headerName string) IExpectHeaders

	// Status provides assertions to the response status code
	//
	// Usage:
	//     Expect().Status().Equal(http.StatusOK)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Status().LessThan(http.StatusBadRequest),
	//     )
	Status() IExpectInt64

	// Trailers provides assertions to the response trailers.
	//
	// If you specify the argument you can directly assert a specific trailer
	//
	// Usage:
	//     Expect().Trailers("X-CustomTrailer").NotEmpty()
	//     Expect().Trailers("X-CustomTrailer").Equal("Foo")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailers("X-CustomTrailer").Empty(),
	//         Expect().Trailers("X-CustomTrailer").NotEqual("Foo"),
	//     )
	Trailers(trailerName string) IExpectHeaders
}

type expect struct {
	cleanPath callPath
}

func newExpect(cleanPath callPath) IExpect {
	return &expect{
		cleanPath: cleanPath,
	}
}

func (exp *expect) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: exp.cleanPath.Push("Custom", []interface{}{fn}),
		Exec: func(hit *hitImpl) error {
			return fn(hit)
		},
	}
}

func (exp *expect) Body() IExpectBody {
	return newExpectBody(exp, exp.cleanPath.Push("Body", nil))
}

func (exp *expect) Headers(headerName string) IExpectHeaders {
	return newExpectHeader(exp.cleanPath.Push("Headers", []interface{}{headerName}), func(hit Hit) []string {
		return hit.Response().Header.Values(headerName)
	})
}

func (exp *expect) Status() IExpectInt64 {
	return newExpectInt64(exp.cleanPath.Push("Status", nil), func(hit Hit) int64 {
		return int64(hit.Response().StatusCode)
	})
}

func (exp *expect) Trailers(trailerName string) IExpectHeaders {
	return newExpectHeader(exp.cleanPath.Push("Trailers", []interface{}{trailerName}), func(hit Hit) []string {
		// we have to read the body to get the trailers
		_, _ = io.Copy(ioutil.Discard, hit.Response().Body().Reader())
		return hit.Response().Trailer.Values(trailerName)
	})
}
