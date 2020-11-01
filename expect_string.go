package hit

import (
	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectString provides assertions for the string type.
type IExpectString interface {
	// Equal expects the string to be equal to the specified value.
	Equal(value string) IStep

	// NotEqual expects the string to be not equal to the specified value.
	NotEqual(value string) IStep

	// Contains expects the string to contain the specified value.
	Contains(value string) IStep

	// NotContains expects the string to not contain the specified value.
	NotContains(value string) IStep

	// Len provides assertions to the string size.
	Len() IExpectInt

	// OneOf expects the string to be equal to one of the specified values.
	OneOf(values ...string) IStep

	// NotOneOf expects the string to be not equal to one of the specified values.
	NotOneOf(values ...string) IStep
}

type expectExpectStringValueCallback func(hit Hit) string
type expectString struct {
	cleanPath     callPath
	valueCallback expectExpectStringValueCallback
}

func newExpectString(cleanPath callPath, valueCallback expectExpectStringValueCallback) IExpectString {
	return &expectString{
		cleanPath:     cleanPath,
		valueCallback: valueCallback,
	}
}

func (v *expectString) Equal(value string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.Equal(v.valueCallback(hit), value)
		},
	}
}

func (v *expectString) NotEqual(value string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.NotEqual(v.valueCallback(hit), value)
		},
	}
}

func (v *expectString) Contains(value string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Contains", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.Contains(v.valueCallback(hit), value)
		},
	}
}

func (v *expectString) NotContains(value string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.NotContains(v.valueCallback(hit), value)
		},
	}
}

func (v *expectString) Len() IExpectInt {
	return newExpectInt(v.cleanPath.Push("Len", nil), func(hit Hit) int {
		return len(v.valueCallback(hit))
	})
}

func (v *expectString) OneOf(values ...string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("OneOf", stringSliceToInterfaceSlice(values)),
		Exec: func(hit *hitImpl) error {
			return minitest.OneOf(v.valueCallback(hit), stringSliceToInterfaceSlice(values)...)
		}}
}
func (v *expectString) NotOneOf(values ...string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotOneOf", stringSliceToInterfaceSlice(values)),
		Exec: func(hit *hitImpl) error {
			return minitest.NotOneOf(v.valueCallback(hit), stringSliceToInterfaceSlice(values)...)
		}}
}
