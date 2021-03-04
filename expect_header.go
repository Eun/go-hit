package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

//nolint:dupl // the methods of IExpectFormValues and IExpectHeaders are the same however the comments are different.
// IExpectHeaders provides assertions on the http response header.
type IExpectHeaders interface {
	// Contains expects the specific header to contain all of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").Contains("application/json")
	//     Expect().Trailers("X-Trailer1").Contains("secret")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").Contains("application/json"),
	//     )
	Contains(values ...interface{}) IStep

	// NotContains expects the specified header to not contain all of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").NotContains("application/json")
	//     Expect().Trailers("X-Trailer").NotContains("secret")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").NotContains("json"),
	//     )
	NotContains(values ...interface{}) IStep

	// OneOf expects the specified header to contain one of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").OneOf("application/json", "text/json")
	//     Expect().Trailers("X-Trailer1").OneOf("foo", "bar")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").OneOf("application/json", "text/json"),
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf expects the specified header to not contain one of the specified values.
	//
	// Usage:
	//     Expect().Headers("Content-Type").NotOneOf("application/json", "text/json")
	//     Expect().Trailers("X-Trailer1").NotOneOf("foo", "bar")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").NotOneOf("application/json", "text/json"),
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty expects the specified header to be empty.
	//
	// Usage:
	//     Expect().Headers("Content-Type").Empty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Authentication").Empty(),
	//     )
	Empty() IStep

	// Empty expects the specified header to not be empty.
	//
	// Usage:
	//     Expect().Headers("Content-Type").NotEmpty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").NotEmpty(),
	//     )
	NotEmpty() IStep

	// Len provides assertions for the length of the specified header.
	//
	// Usage:
	//     Expect().Headers("Content-Type").Len().GreaterThan(0)
	//     Expect().Trailers("X-Trailer1").Len().GreaterThan(0)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").Len().GreaterThan(0),
	//     )
	Len() IExpectInt

	// Equal expects the specified header to be equal the specified value.
	//
	// Usage:
	//     Expect().Headers("Content-Type").Equal("application/json")
	//     Expect().Headers("Content-Length").Equal(11)
	//     Expect().Trailers("X-Trailer1").Equal("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Expect().Headers("Content-Type").Equal("application/json"),
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual expects the specified header to be not equal the specified value.
	//
	// Usage:
	//     Expect().Headers("Content-Type").NotEqual("application/json")
	//     Expect().Headers("Content-Length").NotEqual(11)
	//     Expect().Trailers("X-Trailer1").NotEqual("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").NotEqual("application/json"),
	//     )
	NotEqual(value ...interface{}) IStep

	// First provides assertions for the first value of the specified header.
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
	First() IExpectHeaderValue

	// Last provides assertions for the last value of the specified header.
	//
	// Usage:
	//     Expect().Headers("Content-Type").Last().NotEqual("application/json")
	//     Expect().Headers("Content-Length").Last().NotEqual(11)
	//     Expect().Trailers("X-Trailer1").Last().NotEqual("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").Last().NotEqual("application/json"),
	//     )
	Last() IExpectHeaderValue

	// Nth provides assertions for the nth value of the specified header. (0 = first value)
	//
	// Usage:
	//     Expect().Headers("Content-Type").Nth(0).NotEqual("application/json")
	//     Expect().Headers("Content-Length").Nth(0).NotEqual(11)
	//     Expect().Trailers("X-Trailer1").Nth(0).NotEqual("data")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers("Content-Type").Nth(0).NotEqual("application/json"),
	//     )
	Nth(n int) IExpectHeaderValue
}

type expectHeaderValueCallback func(hit Hit) []string

type expectHeader struct {
	cleanPath     callPath
	valueCallback expectHeaderValueCallback
}

func newExpectHeader(cleanPath callPath, valueCallback expectHeaderValueCallback) IExpectHeaders {
	return &expectHeader{
		cleanPath:     cleanPath,
		valueCallback: valueCallback,
	}
}

func (hdr *expectHeader) Contains(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Contains", values),
		Exec: func(hit *hitImpl) error {
			return minitest.Contains(hdr.valueCallback(hit), values...)
		},
	}
}

func (hdr *expectHeader) NotContains(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotContains", values),
		Exec: func(hit *hitImpl) error {
			return minitest.NotContains(hdr.valueCallback(hit), values...)
		},
	}
}

func (hdr *expectHeader) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit *hitImpl) error {
			return minitest.OneOf(hdr.valueCallback(hit), values...)
		},
	}
}

func (hdr *expectHeader) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit *hitImpl) error {
			return minitest.NotOneOf(hdr.valueCallback(hit), values...)
		},
	}
}

func (hdr *expectHeader) Empty() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit *hitImpl) error {
			return minitest.Empty(hdr.valueCallback(hit))
		},
	}
}

func (hdr *expectHeader) NotEmpty() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotEmpty", nil),
		Exec: func(hit *hitImpl) error {
			return minitest.NotEmpty(hdr.valueCallback(hit))
		},
	}
}

func (hdr *expectHeader) Len() IExpectInt {
	return newExpectInt(hdr.cleanPath.Push("Len", nil), func(hit Hit) int {
		return len(hdr.valueCallback(hit))
	})
}

func (hdr *expectHeader) Equal(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("Equal", values),
		Exec: func(hit *hitImpl) error {
			return minitest.Equal(hdr.valueCallback(hit), values)
		},
	}
}

func (hdr *expectHeader) NotEqual(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: hdr.cleanPath.Push("NotEqual", values),
		Exec: func(hit *hitImpl) error {
			return minitest.NotEqual(hdr.valueCallback(hit), values)
		},
	}
}

func (hdr *expectHeader) First() IExpectHeaderValue {
	return newExpectHeaderValue(hdr.cleanPath.Push("First", nil), func(hit Hit) *string {
		v := hdr.valueCallback(hit)
		size := len(v)
		if size == 0 {
			return nil
		}
		return &v[0]
	})
}

func (hdr *expectHeader) Last() IExpectHeaderValue {
	return newExpectHeaderValue(hdr.cleanPath.Push("Last", nil), func(hit Hit) *string {
		v := hdr.valueCallback(hit)
		size := len(v)
		if size == 0 {
			return nil
		}
		return &v[size-1]
	})
}

func (hdr *expectHeader) Nth(n int) IExpectHeaderValue {
	return newExpectHeaderValue(hdr.cleanPath.Push("Nth", []interface{}{n}), func(hit Hit) *string {
		if n < 0 {
			return nil
		}
		v := hdr.valueCallback(hit)
		size := len(v)
		if size == 0 {
			return nil
		}
		if n+1 > size {
			return nil
		}
		return &v[n]
	})
}
