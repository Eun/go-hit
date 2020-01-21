package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

type IExpectSpecificHeader interface {
	IStep
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
	header string
	expect IExpect
}

func newExpectSpecificHeader(expect IExpect, header string) IExpectSpecificHeader {
	return &expectSpecificHeader{
		expect: expect,
		header: header,
	}
}

func (hdr *expectSpecificHeader) when() StepTime {
	return hdr.expect.when()
}

func (hdr *expectSpecificHeader) exec(hit Hit) error {
	return hdr.expect.exec(hit)
}

// Contains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").Contains("application")
func (hdr *expectSpecificHeader) Contains(v string) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.Contains(hit.Response().Header.Get(hdr.header), v)
	})
}

// NotContains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").NotContains("application")
func (hdr *expectSpecificHeader) NotContains(v string) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.NotContains(hit.Response().Header.Get(hdr.header), v)
	})
}

// OneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").OneOf("application/json", "text/x-json")
func (hdr *expectSpecificHeader) OneOf(values ...interface{}) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.Contains(values, hit.Response().Header.Get(hdr.header))
	})
}

// NotOneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").NotOneOf("application/json", "text/x-json")
func (hdr *expectSpecificHeader) NotOneOf(values ...interface{}) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.NotContains(values, hit.Response().Header.Get(hdr.header))
	})
}

// Empty checks if the header value is empty
// Example:
//           Expect().Headers("Content-Type").Empty()
func (hdr *expectSpecificHeader) Empty() IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.Empty(hit.Response().Header.Get(hdr.header))
	})
}

// Len checks if the the length of the header value is equal to the specified size
// Example:
//           Expect().Header("Content-Type").Len(16)
func (hdr *expectSpecificHeader) Len(size int) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		minitest.Len(hit.Response().Header.Get(hdr.header), size)
	})
}

// Equal checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").Equal("application/json")
func (hdr *expectSpecificHeader) Equal(v interface{}) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
		minitest.NoError(err)
		minitest.Equal(v, compareData)
	})
}

// NotEqual checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").NotEqual("application/json")
func (hdr *expectSpecificHeader) NotEqual(v interface{}) IStep {
	return hdr.expect.Custom(func(hit Hit) {
		compareData, err := makeCompareable(hit.Response().Header.Get(hdr.header), v)
		minitest.NoError(err)
		minitest.NotEqual(v, compareData)
	})
}
