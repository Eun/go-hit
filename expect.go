package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/mohae/deepcopy"
)

// IExpect provides assertions on the http response
type IExpect interface {
	// Body expects the body to be equal the specified value
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Body("Hello World")
	//     Expect().Body().Contains("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().Contains("Hello World"),
	//     )
	Body(value ...interface{}) IExpectBody

	// Header provides assertions to one specific response header.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Header().Contains("Content-Type")
	//     Expect().Header("Content-Type").Equal("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header().Contains("Content-Type"),
	//         Expect().Header("Content-Type").Equal("application/json"),
	//     )
	Header(headerName ...string) IExpectHeader

	// Trailer provides assertions to one specific response header.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Trailer().Contains("Content-Type")
	//     Expect().Trailer("Content-Type").Equal("application/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer().Contains("Content-Type"),
	//         Expect().Trailer("Content-Type").Equal("application/json"),
	//     )
	Trailer(trailerName ...string) IExpectTrailer

	// Status expects the status to be the specified code.
	//
	// If you omit the argument you can fine tune the assertions.
	//
	// Usage:
	//     Expect().Status(200)
	//     Expect().Status().Equal(200)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Status().OneOf(http.StatusOk, http.StatusNoContent),
	//     )
	Status(code ...int) IExpectStatus

	// Custom can be used to expect a custom behaviour.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Custom(func(hit Hit) {
	//               if hit.Response().StatusCode != 200 {
	//                   panic("Expected 200")
	//               }
	//         }),
	//     )
	Custom(fn Callback) IStep
}

type expect struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newExpect(cleanPath clearPath) IExpect {
	return &expect{
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}
}

func (exp *expect) Body(value ...interface{}) IExpectBody {
	return newExpectBody(exp, exp.cleanPath.Push("Body", value), value)
}

func (exp *expect) Custom(fn Callback) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: exp.cleanPath.Push("Custom", []interface{}{fn}),
		Exec: func(hit Hit) error {
			fn(hit)
			return nil
		},
	}
}

func (exp *expect) Header(headerName ...string) IExpectHeader {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}

	return newExpectHeader(exp, exp.cleanPath.Push("Header", args), headerName...)
}

func (exp *expect) Trailer(trailerName ...string) IExpectTrailer {
	args := make([]interface{}, len(trailerName))
	for i := range trailerName {
		args[i] = trailerName[i]
	}

	return newExpectTrailer(exp, exp.cleanPath.Push("Trailer", args), trailerName...)
}

func (exp *expect) Status(code ...int) IExpectStatus {
	args := make([]interface{}, len(code))
	for i := range code {
		args[i] = code[i]
	}
	return newExpectStatus(exp, exp.cleanPath.Push("Status", args), code)
}

func makeCompareable(in, data interface{}) (interface{}, error) {
	compareData := deepcopy.Copy(data)
	err := converter.Convert(in, &compareData)
	if err != nil {
		return nil, err
	}

	return compareData, nil
}
