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
	cleanPath CleanPath
	header    string
}

func newExpectSpecificHeader(expect IExpect, cleanPath CleanPath, header string) IExpectSpecificHeader {
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
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("Contains", []interface{}{v}),
		Exec: func(hit Hit) {
			minitest.Contains(hit.Response().Header.Get(hdr.header), v)
		},
	})
}

// NotContains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").NotContains("application")
func (hdr *expectSpecificHeader) NotContains(v string) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("NotContains", []interface{}{v}),
		Exec: func(hit Hit) {
			minitest.NotContains(hit.Response().Header.Get(hdr.header), v)
		},
	})
}

// OneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").OneOf("application/json", "text/x-json")
func (hdr *expectSpecificHeader) OneOf(values ...interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) {
			minitest.Contains(values, hit.Response().Header.Get(hdr.header))
		},
	})
}

// NotOneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").NotOneOf("application/json", "text/x-json")
func (hdr *expectSpecificHeader) NotOneOf(values ...interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) {
			minitest.NotContains(values, hit.Response().Header.Get(hdr.header))
		},
	})
}

// Empty checks if the header value is empty
// Example:
//           Expect().Headers("Content-Type").Empty()
func (hdr *expectSpecificHeader) Empty() IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) {
			minitest.Empty(hit.Response().Header.Get(hdr.header))
		},
	})
}

// Len checks if the the length of the header value is equal to the specified size
// Example:
//           Expect().Header("Content-Type").Len(16)
func (hdr *expectSpecificHeader) Len(size int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) {
			minitest.Len(hit.Response().Header.Get(hdr.header), size)
		},
	})
}

// Equal checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").Equal("application/json")
func (hdr *expectSpecificHeader) Equal(v interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("Equal", []interface{}{v}),
		Exec: func(hit Hit) {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
			minitest.NoError(err)
			minitest.Equal(v, compareData)
		},
	})
}

// NotEqual checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").NotEqual("application/json")
func (hdr *expectSpecificHeader) NotEqual(v interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("NotEqual", []interface{}{v}),
		Exec: func(hit Hit) {
			compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
			minitest.NoError(err)
			minitest.NotEqual(v, compareData)
		},
	})
}
