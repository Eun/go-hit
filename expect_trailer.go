package hit

import (
	"io"

	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
	"golang.org/x/xerrors"
)

// IExpectTrailer provides assertions on the http response trailer(s)
type IExpectTrailer interface {
	// Contains expects the specified trailer to contain the specified value.
	//
	// Usage:
	//     Expect().Trailer().Contains("Content-Type")
	//     Expect().Trailer("Content-Type").Contains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer().Contains("Content-Type"),
	//         Expect().Trailer("Content-Type").Contains("json"),
	//     )
	Contains(value interface{}) IStep

	// NotContains expects the specified trailer to contain the specified value.
	//
	// Usage:
	//     Expect().Trailer().NotContains("Content-Type")
	//     Expect().Trailer("Content-Type").NotContains("json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer().NotContains("Set-Cookie"),
	//         Expect().Trailer("Content-Type").NotContains("json"),
	//     )
	NotContains(value interface{}) IStep

	// OneOf expects the specified trailer to contain one of the specified values.
	//
	// Usage:
	//     Expect().Trailer().OneOf(map[string]string{"Content-Type": "application/json"}, map[string]string{"Content-Type": "text/json"})
	//     Expect().Trailer("Content-Type").OneOf("application/json", "text/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer("Content-Type").OneOf("application/json", "text/json"),
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf expects the specified trailer to not contain one of the specified values.
	//
	// Usage:
	//     Expect().Trailer().NotOneOf(map[string]string{"Content-Type": "application/json"}, map[string]string{"Content-Type": "text/json"})
	//     Expect().Trailer("Content-Type").NotOneOf("application/json", "text/json")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer("Content-Type").NotOneOf("application/json", "text/json"),
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty expects the specified trailer to be empty.
	//
	// Usage:
	//     Expect().Headers().Empty()
	//     Expect().Trailer("Content-Type").Empty()
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer("Content-Type").Empty(),
	//     )
	Empty() IStep

	// Len expects the specified trailer to be empty.
	//
	// Usage:
	//     Expect().Trailer().Len(1)
	//     Expect().Trailer("Content-Type").Len(16)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer().Len(1),
	//         Expect().Trailer("Content-Type").Len(16),
	//     )
	Len(size int) IStep

	// Equal expects the specified trailer to be equal the specified value.
	//
	// Usage:
	//     Expect().Trailer().Equal(map[string]string{"Content-Type": "application/json"})
	//     Expect().Trailer("Content-Type").Equal("application/json")
	//     Expect().Trailer("Content-Length").Equal(11)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer("Content-Type").Equal("application/json"),
	//     )
	Equal(value interface{}) IStep

	// NotEqual expects the specified trailer to be equal the specified value.
	//
	// Usage:
	//     Expect().Trailer().NotEqual(map[string]string{"Content-Type": "application/json"})
	//     Expect().Trailer("Content-Type").NotEqual("application/json")
	//     Expect().Trailer("Content-Length").NotEqual(11)
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Trailer("Content-Type").NotEqual("application/json"),
	//     )
	NotEqual(value interface{}) IStep
}

func newExpectTrailer(expect IExpect, cleanPath clearPath, trailerName ...string) IExpectTrailer {
	name, ok := internal.GetLastStringArgument(trailerName)
	if ok {
		return newExpectSpecificTrailer(expect, cleanPath, name)
	}
	return newExpectTrailers(expect, cleanPath)
}

type expectTrailers struct {
	expect    IExpect
	cleanPath clearPath
}

func newExpectTrailers(expect IExpect, cleanPath clearPath) IExpectTrailer {
	return &expectTrailers{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (hdr *expectTrailers) Contains(headerName interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Contains", []interface{}{headerName}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Contains(hit.Response().Trailer, headerName)
			return nil
		},
	}
}

func (hdr *expectTrailers) NotContains(headerName interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotContains", []interface{}{headerName}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.NotContains(hit.Response().Trailer, headerName)
			return nil
		},
	}
}

func (hdr *expectTrailers) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Empty(hit.Response().Trailer)
			return nil
		},
	}
}

func (hdr *expectTrailers) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Len(hit.Response().Trailer, size)
			return nil
		},
	}
}

func (hdr *expectTrailers) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			compareData, err := makeCompareable(hit.Response().Trailer, value)
			if err != nil {
				return err
			}
			minitest.Equal(value, compareData)
			return nil
		},
	}
}

func (hdr *expectTrailers) NotEqual(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			compareData, err := makeCompareable(hit.Response().Trailer, value)
			if err != nil {
				return err
			}
			minitest.NotEqual(value, compareData)
			return nil
		},
	}
}

func (hdr *expectTrailers) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			// convert to map[string]string because its easier to compare
			var v []map[string]string
			if err := converter.Convert(values, &v); err != nil {
				return err
			}
			var hdr map[string]string
			if err := converter.Convert(hit.Response().Trailer, &hdr); err != nil {
				return err
			}
			minitest.Contains(v, hdr)
			return nil
		},
	}
}

func (hdr *expectTrailers) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: hdr.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			// convert to map[string]string because its easier to compare
			var v []map[string]string
			if err := converter.Convert(values, &v); err != nil {
				return err
			}
			var hdr map[string]string
			if err := converter.Convert(hit.Response().Trailer, &hdr); err != nil {
				return err
			}
			minitest.NotContains(v, hdr)
			return nil
		},
	}
}

type expectSpecificTrailer struct {
	expect    IExpect
	cleanPath clearPath
	header    string
}

func newExpectSpecificTrailer(expect IExpect, cleanPath clearPath, header string) IExpectTrailer {
	return &expectSpecificTrailer{
		expect:    expect,
		cleanPath: cleanPath,
		header:    header,
	}
}

func (trl *expectSpecificTrailer) Contains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("Contains", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Contains(hit.Response().Trailer.Get(trl.header), value)
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) NotContains(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.NotContains(hit.Response().Trailer.Get(trl.header), value)
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) OneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("OneOf", values),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Contains(values, hit.Response().Trailer.Get(trl.header))
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) NotOneOf(values ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("NotOneOf", values),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.NotContains(values, hit.Response().Trailer.Get(trl.header))
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) Empty() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("Empty", nil),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Empty(hit.Response().Trailer.Get(trl.header))
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) Len(size int) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("Len", []interface{}{size}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			minitest.Len(hit.Response().Trailer.Get(trl.header), size)
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) Equal(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			compareData, err := makeCompareable(hit.Response().Trailer.Get(trl.header), value)
			if err != nil {
				return err
			}
			minitest.Equal(value, compareData)
			return nil
		},
	}
}

func (trl *expectSpecificTrailer) NotEqual(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      ExpectStep,
		ClearPath: trl.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit Hit) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(internal.DevNullWriter(), hit.Response().Body().Reader())
			compareData, err := makeCompareable(hit.Response().Trailer.Get(trl.header), value)
			if err != nil {
				return err
			}
			minitest.NotEqual(value, compareData)
			return nil
		},
	}
}

type finalExpectTrailer struct {
	IStep
	message string
}

func (trl *finalExpectTrailer) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(trl.message)
		},
	}
}

func (trl *finalExpectTrailer) Contains(interface{}) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) NotContains(interface{}) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) OneOf(...interface{}) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) NotOneOf(...interface{}) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) Empty() IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) Len(int) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) Equal(interface{}) IStep {
	return trl.fail()
}

func (trl *finalExpectTrailer) NotEqual(interface{}) IStep {
	return trl.fail()
}
