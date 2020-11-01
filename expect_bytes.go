package hit

import (
	"bytes"

	"github.com/Eun/go-hit/internal/minitest"
)

// IExpectBytes provides assertions for the byte slice type.
type IExpectBytes interface {
	// Equal expects the body to be equal to the specified value.
	//
	// Usage:
	//     Expect().Body().Bytes().Equal([]byte("Hello World"))
	Equal(value []byte) IStep

	// NotEqual expects the body to be not equal to the specified value
	//
	// Usage:
	//     Expect().Body().Bytes().NotEqual([]byte("Hello World"))
	NotEqual(value []byte) IStep

	// Contains expects the body to contain the specified value.
	//
	// Usage:
	//     Expect().Body().Bytes().Contains([]byte("Hello World"))
	Contains(value []byte) IStep

	// NotContains expects the body to not contain the specified value.
	//
	// Usage:
	//     Expect().Body().Bytes().NotContains([]byte("Hello World"))
	NotContains(value []byte) IStep

	// Len provides assertions to the body size.
	//
	// Usage:
	//     Expect().Body().Bytes().Len().Equal(10)
	Len() IExpectInt
}

type expectBytesValueCallback func(hit Hit) []byte
type expectBytes struct {
	cleanPath     callPath
	valueCallback expectBytesValueCallback
}

func newExpectBytes(cleanPath callPath, valueCallback expectBytesValueCallback) IExpectBytes {
	return &expectBytes{
		cleanPath:     cleanPath,
		valueCallback: valueCallback,
	}
}

func (v *expectBytes) Equal(value []byte) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Equal", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.Equal(v.valueCallback(hit), value)
		},
	}
}

func (v *expectBytes) NotEqual(value []byte) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotEqual", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return minitest.NotEqual(v.valueCallback(hit), value)
		},
	}
}

func (v *expectBytes) Contains(value []byte) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("Contains", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			buf := v.valueCallback(hit)
			if !bytes.Contains(buf, value) {
				return minitest.Errorf(`%s does not contain %s`, minitest.PrintValue(buf), minitest.PrintValue(value))
			}
			return nil
		},
	}
}

func (v *expectBytes) NotContains(value []byte) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     ExpectStep,
		CallPath: v.cleanPath.Push("NotContains", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			buf := v.valueCallback(hit)
			if bytes.Contains(buf, value) {
				return minitest.Errorf(`%s should not contain %s`, minitest.PrintValue(buf), minitest.PrintValue(value))
			}
			return nil
		},
	}
}

func (v *expectBytes) Len() IExpectInt {
	return newExpectInt(v.cleanPath.Push("Len", nil), func(hit Hit) int {
		return len(v.valueCallback(hit))
	})
}
