package hit

import (
	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
)

type IExpectSpecificHeader interface {
	IStep
	Contains(v string) IStep
	OneOf(values ...interface{}) IStep
	Empty() IStep
	Len(size int) IStep
	Equal(v interface{}) IStep
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

func (exp *expectSpecificHeader) when() StepTime {
	return exp.expect.when()
}

func (exp *expectSpecificHeader) exec(hit Hit) {
	exp.expect.exec(hit)
}

// Contains checks if the header value contains the specified value
// Example:
//           Expect().Header("Content-Type").Contains("application")
func (hdr *expectSpecificHeader) Contains(v string) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Contains(hit.Response().Header.Get(hdr.header), v)
	})
}

// OneOf checks if the header value is one of the specified values
// Example:
//           Expect().Header("Content-Type").OneOf("application/json", "text/x-json")
func (hdr *expectSpecificHeader) OneOf(values ...interface{}) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Contains(values, hit.Response().Header.Get(hdr.header))
	})
}

// Empty checks if the header value is empty
// Example:
//           Expect().Headers("Content-Type").Empty()
func (hdr *expectSpecificHeader) Empty() IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Empty(hit.Response().Header.Get(hdr.header))
	})
}

// Len checks if the the length of the header value is equal to the specified size
// Example:
//           Expect().Header("Content-Type").Len(16)
func (hdr *expectSpecificHeader) Len(size int) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Len(hit.Response().Header.Get(hdr.header), size)
	})
}

// Equal checks if the header value is equal to the specified value
// Example:
//           Expect().Headers("Content-Type").Equal("application/json")
func (hdr *expectSpecificHeader) Equal(v interface{}) IStep {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		compareData, err := converter.Convert(hit.Response().Header.Get(hdr.header), v, convert.Options.ConvertEmbeddedStructToParentType())
		et.NoError(err)
		et.Equal(v, compareData)
	})
}
