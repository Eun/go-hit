package hit

import (
	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
)

type IExpectHeaders interface {
	IStep
	Contains(v string) IStep
	Empty() IStep
	Len(size int) IStep
	Equal(v interface{}) IStep
	Get(header string) IExpectSpecificHeader
}

type expectHeaders struct {
	expect IExpect
}

func newExpectHeaders(expect IExpect) IExpectHeaders {
	return &expectHeaders{expect}
}

func (exp *expectHeaders) when() StepTime {
	return exp.expect.when()
}

func (exp *expectHeaders) exec(hit Hit) {
	exp.expect.exec(hit)
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
func (hdr *expectHeaders) Get(header string) IExpectSpecificHeader {
	return newExpectSpecificHeader(hdr.expect, header)
}
