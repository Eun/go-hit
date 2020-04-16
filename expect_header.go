package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
	"github.com/Eun/go-hit/internal/misc"
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
	//     Expect().Header().Empty()
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
	name, ok := misc.GetLastStringArgument(headerName)
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
			return minitest.Error.Contains(hit.Response().Header, headerName)
		},
	}
}

func (hdr *expectHeaders) NotContains(headerName interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{headerName}),
		Exec: func(hit Hit) error {
			return minitest.Error.NotContains(hit.Response().Header, headerName)
		},
	}
}

func (hdr *expectHeaders) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			return minitest.Error.Empty(hit.Response().Header)
		},
	}
}

func (hdr *expectHeaders) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			return minitest.Error.Len(hit.Response().Header, size)
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
			return minitest.Error.Equal(value, compareData)
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
			return minitest.Error.NotEqual(value, compareData)
		},
	}
}

func (hdr *expectHeaders) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			// convert to map[string]string because its easier to compare
			var v []map[string]string
			if err := converter.Convert(values, &v); err != nil {
				return err
			}
			var hdr map[string]string
			if err := converter.Convert(hit.Response().Header, &hdr); err != nil {
				return err
			}
			return minitest.Error.Contains(v, hdr)
		},
	}
}

func (hdr *expectHeaders) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			// convert to map[string]string because its easier to compare
			var v []map[string]string
			if err := converter.Convert(values, &v); err != nil {
				return err
			}
			var hdr map[string]string
			if err := converter.Convert(hit.Response().Header, &hdr); err != nil {
				return err
			}
			return minitest.Error.NotContains(v, hdr)
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
			return minitest.Error.Contains(hit.Response().Header.Get(hdr.header), value)
		},
	}
}

func (hdr *expectSpecificHeader) NotContains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit Hit) error {
			return minitest.Error.NotContains(hit.Response().Header.Get(hdr.header), value)
		},
	}
}

func (hdr *expectSpecificHeader) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			return minitest.Error.Contains(values, hit.Response().Header.Get(hdr.header))
		},
	}
}

func (hdr *expectSpecificHeader) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			return minitest.Error.NotContains(values, hit.Response().Header.Get(hdr.header))
		},
	}
}

func (hdr *expectSpecificHeader) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			return minitest.Error.Empty(hit.Response().Header.Get(hdr.header))
		},
	}
}

func (hdr *expectSpecificHeader) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			return minitest.Error.Len(hit.Response().Header.Get(hdr.header), size)
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
			return minitest.Error.Equal(value, compareData)
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
			return minitest.Error.NotEqual(value, compareData)
		},
	}
}
