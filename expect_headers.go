package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectHeaders provides assertions on the http response headers
type IExpectHeaders interface {
	// Contains expects the headers to contain the specified header.
	//
	// Usage:
	//     Expect().Headers().Contains("Content-Type")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().Contains("Content-Type"),
	//     )
	Contains(headerName string) IStep

	// NotContains expects the headers to not contain the specified header.
	//
	// Usage:
	//     Expect().Headers().NotContains("Content-Type")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().NotContains("Content-Type"),
	//     )
	NotContains(headerName string) IStep

	// Empty expects the headers to be empty.
	//
	// Usage:
	//     Expect().Headers().Empty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().Empty(),
	//     )
	Empty() IStep

	// Len expects the header amount to be equal to the specified amount.
	//
	// Usage:
	//     Expect().Headers().Len(3)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().Len(3),
	//     )
	Len(size int) IStep
	// Equal expects the headers to be equal to the specified value.
	//
	// Usage:
	//     Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"}),
	//     )
	Equal(value interface{}) IStep

	// NotEqual expects the headers to be equal to the specified value.
	//
	// Usage:
	//     Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"}),
	//     )
	NotEqual(value interface{}) IStep

	// Get a specific header
	// Usage:
	//     Expect().Headers().Get("Content-Type").Equal("application/json")
	//     Expect().Headers().Get("Content-Type").Contains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Headers().Get("Content-Type").Contains("json"),
	//     )
	Get(headerName string) IExpectSpecificHeader
}

type expectHeaders struct {
	expect    IExpect
	cleanPath clearPath
}

func newExpectHeaders(expect IExpect, cleanPath clearPath) IExpectHeaders {
	return &expectHeaders{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (hdr *expectHeaders) Contains(headerName string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Contains", []interface{}{headerName}),
		Exec: func(hit Hit) error {
			minitest.Contains(hit.Response().Header, headerName)
			return nil
		},
	}
}

func (hdr *expectHeaders) NotContains(headerName string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{headerName}),
		Exec: func(hit Hit) error {
			minitest.NotContains(hit.Response().Header, headerName)
			return nil
		},
	}
}

func (hdr *expectHeaders) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			minitest.Empty(hit.Response().Header)
			return nil
		},
	}
}

func (hdr *expectHeaders) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			minitest.Len(hit.Response().Header, size)
			return nil
		},
	}
}

func (hdr *expectHeaders) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header, value)
			if err != nil {
				return err
			}
			minitest.Equal(value, compareData)
			return nil
		},
	}
}

func (hdr *expectHeaders) NotEqual(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header, value)
			if err != nil {
				return err
			}
			minitest.NotEqual(value, compareData)
			return nil
		},
	}
}

func (hdr *expectHeaders) Get(headerName string) IExpectSpecificHeader {
	return newExpectSpecificHeader(hdr.expect, hdr.cleanPath.Push("Get", []interface{}{headerName}), headerName)
}
