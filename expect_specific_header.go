package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectSpecificHeader provides assertions on one specific http response header
type IExpectSpecificHeader interface {
	// Contains expects the specified header to contain the specified value.
	//
	// Usage:
	//     Expect().Header("Content-Type").Contains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").Contains("json"),
	//     )
	Contains(value interface{}) IStep

	// NotContains expects the specified header to contain the specified value.
	//
	// Usage:
	//     Expect().Header("Content-Type").NotContains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").NotContains("json"),
	//     )
	NotContains(value interface{}) IStep

	// OneOf expects the specified header to contain one of the specified values.
	//
	// Usage:
	//     Expect().Header("Content-Type").OneOf("application/json", "text/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").OneOf("application/json", "text/json"),
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf expects the specified header to not contain one of the specified values.
	//
	// Usage:
	//     Expect().Header("Content-Type").NotOneOf("application/json", "text/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").NotOneOf("application/json", "text/json"),
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty expects the specified header to be empty.
	//
	// Usage:
	//     Expect().Header("Content-Type").Empty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").Empty(),
	//     )
	Empty() IStep

	// Len expects the specified header to be empty.
	//
	// Usage:
	//     Expect().Header("Content-Type").Len(16)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").Len(16),
	//     )
	Len(size int) IStep

	// Equal expects the specified header to be equal the specified value.
	//
	// Usage:
	//     Expect().Header("Content-Type").Equal("application/json")
	//     Expect().Header("Content-Length").Equal(11)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").Equal("application/json"),
	//     )
	Equal(value interface{}) IStep

	// NotEqual expects the specified header to be equal the specified value.
	//
	// Usage:
	//     Expect().Header("Content-Type").NotEqual("application/json")
	//     Expect().Header("Content-Length").NotEqual(11)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header("Content-Type").NotEqual("application/json"),
	//     )
	NotEqual(value interface{}) IStep
}
type expectSpecificHeader struct {
	expect    IExpect
	cleanPath clearPath
	header    string
}

func newExpectSpecificHeader(expect IExpect, cleanPath clearPath, header string) IExpectSpecificHeader {
	return &expectSpecificHeader{
		expect:    expect,
		cleanPath: cleanPath,
		header:    header,
	}
}

func (hdr *expectSpecificHeader) Contains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Contains", []interface{}{value}),
		Exec: func(hit Hit) error {
			minitest.Contains(hit.Response().Header.Get(hdr.header), value)
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) NotContains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit Hit) error {
			minitest.NotContains(hit.Response().Header.Get(hdr.header), value)
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			minitest.Contains(values, hit.Response().Header.Get(hdr.header))
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			minitest.NotContains(values, hit.Response().Header.Get(hdr.header))
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			minitest.Empty(hit.Response().Header.Get(hdr.header))
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			minitest.Len(hit.Response().Header.Get(hdr.header), size)
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), value)
			if err != nil {
				return err
			}
			minitest.Equal(value, compareData)
			return nil
		},
	}
}

func (hdr *expectSpecificHeader) NotEqual(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), value)
			if err != nil {
				return err
			}
			minitest.NotEqual(value, compareData)
			return nil
		},
	}
}
