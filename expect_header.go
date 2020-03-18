package hit

import (
	"net/http"

	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
	"golang.org/x/xerrors"
)

// IExpectHeader provides assertions on the http response header(s)
type IExpectHeader interface {
	// Contains expects the specified header to contain the specified value.
	//
	// Usage:
	//     Expect().Header().Contains("Content-Type")
	//     Expect().Header("Content-Type").Contains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header().Contains("Content-Type"),
	//         Expect().Header("Content-Type").Contains("json"),
	//     )
	Contains(value interface{}) IStep

	// NotContains expects the specified header to contain the specified value.
	//
	// Usage:
	//     Expect().Header().NotContains("Content-Type")
	//     Expect().Header("Content-Type").NotContains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header().NotContains("Set-Cookie"),
	//         Expect().Header("Content-Type").NotContains("json"),
	//     )
	NotContains(value interface{}) IStep

	// OneOf expects the specified header to contain one of the specified values.
	//
	// Usage:
	//     Expect().Header().OneOf(map[string]string{"Content-Type": "application/json"}, map[string]string{"Content-Type": "text/json"})
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
	//     Expect().Header().NotOneOf(map[string]string{"Content-Type": "application/json"}, map[string]string{"Content-Type": "text/json"})
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
	//     Expect().Headers().Empty()
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
	//     Expect().Header().Len(1)
	//     Expect().Header("Content-Type").Len(16)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Header().Len(1),
	//         Expect().Header("Content-Type").Len(16),
	//     )
	Len(size int) IStep

	// Equal expects the specified header to be equal the specified value.
	//
	// Usage:
	//     Expect().Header().Equal(map[string]string{"Content-Type": "application/json"})
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
	//     Expect().Header().NotEqual(map[string]string{"Content-Type": "application/json"})
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

func newExpectHeader(expect IExpect, cleanPath clearPath, headerName ...string) IExpectHeader {
	name, ok := internal.GetLastStringArgument(headerName)
	if ok {
		return newExpectSpecificHeader(expect, cleanPath, name)
	}
	return newExpectHeaders(expect, cleanPath)
}

type expectHeaders struct {
	expect    IExpect
	cleanPath clearPath
}

func newExpectHeaders(expect IExpect, cleanPath clearPath) IExpectHeader {
	return &expectHeaders{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (hdr *expectHeaders) Contains(headerName interface{}) IStep {
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

func (hdr *expectHeaders) NotContains(headerName interface{}) IStep {
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

func (hdr *expectHeaders) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			var h []http.Header
			if err := converter.Convert(values, &h); err != nil {
				return err
			}
			minitest.Contains(h, hit.Response().Header)
			return nil
		},
	}
}

func (hdr *expectHeaders) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			minitest.NotContains(values, hit.Response().Header)
			return nil
		},
	}
}

type expectSpecificHeader struct {
	expect    IExpect
	cleanPath clearPath
	header    string
}

func newExpectSpecificHeader(expect IExpect, cleanPath clearPath, header string) IExpectHeader {
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

type finalExpectHeader struct {
	IStep
	message string
}

func (hdr *finalExpectHeader) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(hdr.message)
		},
	}
}

func (hdr *finalExpectHeader) Contains(interface{}) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) NotContains(interface{}) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) OneOf(...interface{}) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) NotOneOf(...interface{}) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) Empty() IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) Len(int) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) Equal(interface{}) IStep {
	return hdr.fail()
}

func (hdr *finalExpectHeader) NotEqual(interface{}) IStep {
	return hdr.fail()
}
