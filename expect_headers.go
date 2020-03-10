package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectHeaders interface {
	Contains(v string) IStep
	NotContains(v string) IStep
	Empty() IStep
	Len(size int) IStep
	Equal(v interface{}) IStep
	NotEqual(v interface{}) IStep
	Get(header string) IExpectSpecificHeader
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

// Contains checks if the specified header is present
// Example:
//           Expect().Headers().Contains("Content-Type")
func (hdr *expectHeaders) Contains(v string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Contains", []interface{}{v}),
		Exec: func(hit Hit) error {
			minitest.Contains(hit.Response().Header, v)
			return nil
		},
	}
}

// NotContains checks if the specified header is not present
// Example:
//           Expect().Headers().NotContains("Content-Type")
func (hdr *expectHeaders) NotContains(v string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{v}),
		Exec: func(hit Hit) error {
			minitest.NotContains(hit.Response().Header, v)
			return nil
		},
	}
}

// Empty checks if the headers are empty
// Example:
//           Expect().Headers().Empty()
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

// Len checks if the amount of headers are equal to the specified size
// Example:
//           Expect().Headers().Len(0)
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

// Equal checks if the headers are equal to the specified one
// Example:
//           Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"})
func (hdr *expectHeaders) Equal(v interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Equal", []interface{}{v}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header, v)
			if err != nil {
				return err
			}
			minitest.Equal(v, compareData)
			return nil
		},
	}
}

// NotEqual checks if the headers are equal to the specified one
// Example:
//           Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"})
func (hdr *expectHeaders) NotEqual(v interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotEqual", []interface{}{v}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header, v)
			if err != nil {
				return err
			}
			minitest.NotEqual(v, compareData)
			return nil
		},
	}
}

// Get a specific header
// Examples:
//           Expect().Headers().Get("Content-Type").Equal("application/json")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (hdr *expectHeaders) Get(header string) IExpectSpecificHeader {
	return newExpectSpecificHeader(hdr.expect, hdr.cleanPath.Push("Get", []interface{}{header}), header)
}
