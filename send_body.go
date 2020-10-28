package hit

import (
	"io"

	"github.com/Eun/go-hit/httpbody"
)

// ISendBody provides methods to set request body.
type ISendBody interface {
	// Bool sets the request body to the specified bool.
	//
	// Usage:
	//     Send().Body().Bool(true)
	Bool(value bool) IStep

	// Bytes sets the request body to the specified byte slice.
	//
	// Usage:
	//     Send().Body().Bytes([]byte{0x48, 0x65, 0x6C, 0x6C, 0x6F})
	Bytes(value []byte) IStep

	// Float32 sets the request body to the specified float32.
	//
	// Usage:
	//     Send().Body().Float32(3.2)
	Float32(value float32) IStep

	// Float64 sets the request body to the specified float64.
	//
	// Usage:
	//     Send().Body().Float64(3.2)
	Float64(value float64) IStep

	// FormValues sets the request body to the specified form values.
	//
	// Usage:
	//     Send().Body().FormValues("username").Add("admin")
	//     Send().Body().FormValues("password").Add("secret")
	FormValues(name string) ISendFormValues

	// Int sets the request body to the specified int.
	//
	// Usage:
	//     Send().Body().Int(3)
	Int(value int) IStep

	// Int16 sets the request body to the specified int16.
	//
	// Usage:
	//     Send().Body().Int16(3)
	Int16(value int16) IStep

	// Int32 sets the request body to the specified int32.
	//
	// Usage:
	//     Send().Body().Int32(3)
	Int32(value int32) IStep

	// Int64 sets the request body to the specified int64.
	//
	// Usage:
	//     Send().Body().Int64(3)
	Int64(value int64) IStep

	// Int8 sets the request body to the specified int8.
	//
	// Usage:
	//     Send().Body().Int8(3)
	Int8(value int8) IStep

	// JSON sets the request body to the json representation of the specified value.
	//
	// Usage:
	//     Send().Body().JSON(map[string]interface{}{"Name": "Joe"})
	JSON(value interface{}) IStep

	// Reader sets the request body to the specified reader.
	//
	// Usage:
	//     Send().Body().Reader(bytes.NewReader([]byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}))
	Reader(value io.Reader) IStep

	// String sets the request body to the specified string.
	//
	// Usage:
	//     Send().Body().String("Hello World")
	String(value string) IStep

	// Uint sets the request body to the specified uint.
	//
	// Usage:
	//     Send().Body().Uint(3)
	Uint(value uint) IStep

	// Uint16 sets the request body to the specified uint16.
	//
	// Usage:
	//     Send().Body().Uint16(3)
	Uint16(value uint16) IStep

	// Uint32 sets the request body to the specified uint32.
	//
	// Usage:
	//     Send().Body().Uint32(3)
	Uint32(value uint32) IStep

	// Uint64 sets the request body to the specified uint64.
	//
	// Usage:
	//     Send().Body().Uint64(3)
	Uint64(value uint64) IStep

	// Uint8 sets the request body to the specified uint8.
	//
	// Usage:
	//     Send().Body().Uint8(3)
	Uint8(value uint8) IStep

	// XML sets the request body to the XML representation of the specified value.
	//
	// Usage:
	//     Send().Body().XML([]string{"A", "B"})
	XML(value interface{}) IStep
}

type sendBody struct {
	cleanPath callPath
}

func newSendBody(clearPath callPath) ISendBody {
	return &sendBody{
		cleanPath: clearPath,
	}
}

func (body *sendBody) Bool(value bool) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Bool", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetBool(value)
			return nil
		},
	}
}

func (body *sendBody) Bytes(value []byte) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Bytes", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetBytes(value)
			return nil
		},
	}
}

func (body *sendBody) Float32(value float32) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Float32", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetFloat32(value)
			return nil
		},
	}
}

func (body *sendBody) Float64(value float64) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Float64", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetFloat64(value)
			return nil
		},
	}
}

func (body *sendBody) FormValues(name string) ISendFormValues {
	return newSendFormValues(body.cleanPath.Push("FormValues", []interface{}{name}), func(hit Hit) *httpbody.URLValues {
		return hit.Request().body.MustFormValues()
	}, name)
}

func (body *sendBody) Int(value int) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Int", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetInt(value)
			return nil
		},
	}
}

func (body *sendBody) Int16(value int16) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Int16", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetInt16(value)
			return nil
		},
	}
}

func (body *sendBody) Int32(value int32) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Int32", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetInt32(value)
			return nil
		},
	}
}

func (body *sendBody) Int64(value int64) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Int64", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetInt64(value)
			return nil
		},
	}
}

func (body *sendBody) Int8(value int8) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Int8", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetInt8(value)
			return nil
		},
	}
}

func (body *sendBody) JSON(value interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("JSON", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return hit.Request().Body().JSON().Set(value)
		},
	}
}

func (body *sendBody) Reader(value io.Reader) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Reader", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetReader(value)
			return nil
		},
	}
}

func (body *sendBody) String(value string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("String", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetString(value)
			return nil
		},
	}
}

func (body *sendBody) Uint(value uint) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Uint", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetUint(value)
			return nil
		},
	}
}

func (body *sendBody) Uint16(value uint16) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Uint16", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetUint16(value)
			return nil
		},
	}
}

func (body *sendBody) Uint32(value uint32) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Uint32", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetUint32(value)
			return nil
		},
	}
}

func (body *sendBody) Uint64(value uint64) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Uint64", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetUint64(value)
			return nil
		},
	}
}

func (body *sendBody) Uint8(value uint8) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("Uint8", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			hit.Request().Body().SetUint8(value)
			return nil
		},
	}
}

func (body *sendBody) XML(value interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: body.cleanPath.Push("XML", []interface{}{value}),
		Exec: func(hit *hitImpl) error {
			return hit.Request().Body().XML().Set(value)
		},
	}
}
