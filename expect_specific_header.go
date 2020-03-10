package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectSpecificHeader interface {
	Contains(v string) IStep
	NotContains(v string) IStep
	OneOf(values ...interface{}) IStep
	NotOneOf(values ...interface{}) IStep
	Empty() IStep
	Len(size int) IStep
	Equal(v interface{}) IStep
	NotEqual(v interface{}) IStep
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

// Contains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").Contains("application")
func (hdr *expectSpecificHeader) Contains(v string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Contains", []interface{}{v}),
		Exec: func(hit Hit) error {
			minitest.Contains(hit.Response().Header.Get(hdr.header), v)
			return nil
		},
	}
}

// NotContains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").NotContains("application")
func (hdr *expectSpecificHeader) NotContains(v string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{v}),
		Exec: func(hit Hit) error {
			minitest.NotContains(hit.Response().Header.Get(hdr.header), v)
			return nil
		},
	}
}

// OneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").OneOf("application/json", "text/x-json")
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

// NotOneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").NotOneOf("application/json", "text/x-json")
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

// Empty checks if the header value is empty
// Example:
//           Expect().Headers("Content-Type").Empty()
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

// Len checks if the the length of the header value is equal to the specified size
// Example:
//           Expect().Header("Content-Type").Len(16)
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

// Equal checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").Equal("application/json")
func (hdr *expectSpecificHeader) Equal(v interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Equal", []interface{}{v}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
			if err != nil {
				return err
			}
			minitest.Equal(v, compareData)
			return nil
		},
	}
}

// NotEqual checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").NotEqual("application/json")
func (hdr *expectSpecificHeader) NotEqual(v interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotEqual", []interface{}{v}),
		Exec: func(hit Hit) error {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
			if err != nil {
				return err
			}
			minitest.NotEqual(v, compareData)
			return nil
		},
	}
}
