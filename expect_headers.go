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
	hit       Hit
	cleanPath CleanPath
}

func newExpectHeaders(expect IExpect, hit Hit, cleanPath CleanPath) IExpectHeaders {
	return &expectHeaders{
		expect:    expect,
		hit:       hit,
		cleanPath: cleanPath,
	}
}

// Contains checks if the specified header is present
// Example:
//           Expect().Headers().Contains("Content-Type")
func (hdr *expectHeaders) Contains(v string) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: hdr.cleanPath.Push("Contains"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			minitest.Contains(hit.Response().Header, v)
		},
	})
}

// NotContains checks if the specified header is not present
// Example:
//           Expect().Headers().NotContains("Content-Type")
func (hdr *expectHeaders) NotContains(v string) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: append(hdr.cleanPath, "NotContains"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			minitest.NotContains(hit.Response().Header, v)
		},
	})
}

// Empty checks if the headers are empty
// Example:
//           Expect().Headers().Empty()
func (hdr *expectHeaders) Empty() IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: append(hdr.cleanPath, "Empty"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			minitest.Empty(hit.Response().Header)
		},
	})
}

// Len checks if the amount of headers are equal to the specified size
// Example:
//           Expect().Headers().Len(0)
func (hdr *expectHeaders) Len(size int) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: append(hdr.cleanPath, "Len"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			minitest.Len(hit.Response().Header, size)
		},
	})
}

// Equal checks if the headers are equal to the specified one
// Example:
//           Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"})
func (hdr *expectHeaders) Equal(v interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: append(hdr.cleanPath, "Equal"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			compareData, err := makeCompareable(hit.Response().Header, v)
			minitest.NoError(err)
			minitest.Equal(v, compareData)
		},
	})
}

// NotEqual checks if the headers are equal to the specified one
// Example:
//           Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"})
func (hdr *expectHeaders) NotEqual(v interface{}) IStep {
	return custom(Step{
		When:      ExpectStep,
		CleanPath: append(hdr.cleanPath, "NotEqual"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			compareData, err := makeCompareable(hit.Response().Header, v)
			minitest.NoError(err)
			minitest.NotEqual(v, compareData)
		},
	})
}

// Get a specific header
// Examples:
//           Expect().Headers().Get("Content-Type").Equal("application/json")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (hdr *expectHeaders) Get(header string) IExpectSpecificHeader {
	return newExpectSpecificHeader(hdr.expect, hdr.hit, hdr.cleanPath.Push("Get"), header)
}
