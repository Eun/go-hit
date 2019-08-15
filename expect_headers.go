package hit

import (
	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
)

type expectHeaders struct {
	expect *defaultExpect
}

func newExpectHeaders(expect *defaultExpect) *expectHeaders {
	return &expectHeaders{
		expect: expect,
	}
}

// Contains checks if the specified header is present
// Example:
//           Expect().Headers().Contains("Content-Type")
func (hdr *expectHeaders) Contains(v string) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Contains(hit.Response().Header, v)
	})
}

// Empty checks if the headers are empty
// Example:
//           Expect().Headers().Empty()
func (hdr *expectHeaders) Empty() IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Empty(hit.Response().Header)
	})
}

// Len checks if the amount of headers are equal to the specified size
// Example:
//           Expect().Headers().Len(0)
func (hdr *expectHeaders) Len(size int) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Len(hit.Response().Header, size)
	})
}

// Equal checks if the headers are equal to the specified one
// Example:
//           Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"})
func (hdr *expectHeaders) Equal(v interface{}) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		compareData, err := converter.Convert(hit.Response().Header, v, convert.Options.ConvertEmbeddedStructToParentType())
		et.NoError(err)
		et.Equal(v, compareData)
	})
}

// Get a specific header
// Examples:
//           Expect().Headers().Get("Content-Type").Equal("application/json")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (hdr *expectHeaders) Get(header string) *expectSpecificHeader {
	return newExpectSpecificHeader(hdr.expect, header)
}
