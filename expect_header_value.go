package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectHeaderValue provides assertions on the http response header value.
type IExpectHeaderValue interface {
	// Contains expects the specific header value to contain all of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().Contains("application/json")
	//     Expect().Trailers("X-Trailer1").First().Contains("secret")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").First().Contains("application/json"),
	//     )
	Contains(values ...interface{}) IStep

	// NotContains expects the specified header value to not contain all of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().NotContains("application/json")
	//     Expect().Trailers("X-Trailer").First().NotContains("secret")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").First().NotContains("json"),
	//     )
	NotContains(values ...interface{}) IStep

	// OneOf expects the specified header value to contain one of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().OneOf("application/json", "text/json")
	//     Expect().Trailers("X-Trailer1").First().OneOf("foo", "bar")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").First().OneOf("application/json", "text/json"),
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf expects the specified header value to not contain one of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().NotOneOf("application/json", "text/json")
	//     Expect().Trailers("X-Trailer1").First().NotOneOf("foo", "bar")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").First().NotOneOf("application/json", "text/json"),
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty expects the specified header value to be empty.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().Empty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Authentication").First().Empty(),
	//     )
	Empty() IStep

	// Empty expects the specified header value to not be empty.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().NotEmpty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").First().NotEmpty(),
	//     )
	NotEmpty() IStep

	// Len expects the specified header value to be the same length then specified.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().Len().GreaterThan(0)
	//     Expect().Trailers("X-Trailer1").First().Len().GreaterThan(0)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").First().Len().GreaterThan(0),
	//     )
	Len() IExpectInt

	// Equal expects the specified header value to be equal the specified value.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().Equal("application/json")
	//     Expect().Headers("Content-Length").First().Equal(11)
	//     Expect().Trailers("X-Trailer1").First().Equal("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").First().Equal("application/json"),
	//     )
	Equal(value interface{}) IStep

	// NotEqual expects the specified header value to be not equal to one of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").First().NotEqual("application/json")
	//     Expect().Headers("Content-Length").First().NotEqual(11)
	//     Expect().Trailers("X-Trailer1").First().NotEqual("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").First().NotEqual("application/json"),
	//     )
	NotEqual(value ...interface{}) IStep
}

type expectHeaderValueValueCallback func(hit Hit) *string

type expectHeaderValue struct {
	cleanPath     callPath
	valueCallback expectHeaderValueValueCallback
}

func newExpectHeaderValue(cleanPath callPath, valueCallback expectHeaderValueValueCallback) IExpectHeaderValue {
	return &expectHeaderValue{
		cleanPath:     cleanPath,
		valueCallback: valueCallback,
	}
}

func (hdr *expectHeaderValue) Contains(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Contains", values),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.Contains(nil, values...)
			}
			return minitest.Contains(*v, values...)
		},
	}
}

func (hdr *expectHeaderValue) NotContains(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotContains", values),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.NotContains(nil, values...)
			}
			return minitest.NotContains(*v, values...)
		},
	}
}

func (hdr *expectHeaderValue) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.OneOf(nil, values...)
			}
			return minitest.OneOf(*v, values...)
		},
	}
}

func (hdr *expectHeaderValue) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.NotOneOf(nil, values...)
			}
			return minitest.NotOneOf(*v, values...)
		},
	}
}

func (hdr *expectHeaderValue) Empty() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.Empty(nil)
			}
			return minitest.Empty(*v)
		},
	}
}

func (hdr *expectHeaderValue) NotEmpty() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotEmpty", nil),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.NotEmpty(nil)
			}
			return minitest.NotEmpty(*v)
		},
	}
}

func (hdr *expectHeaderValue) Len() IExpectInt {
	return newExpectInt(hdr.cleanPath.Push("Len", nil), func(hit Hit) int {
		v := hdr.valueCallback(hit)
		if v == nil {
			return 0
		}
		return len(*v)
	})
}

func (hdr *expectHeaderValue) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.Equal(nil, value)
			}
			return minitest.Equal(*v, value)
		},
	}
}

func (hdr *expectHeaderValue) NotEqual(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotEqual", values),
		Exec: func(hit *hitImpl) error {
			v := hdr.valueCallback(hit)
			if v == nil {
				return minitest.NotEqual(nil, values...)
			}
			return minitest.NotEqual(*v, values...)
		},
	}
}
