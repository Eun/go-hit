package hit

import (
	"github.com/Eun/go-hit/httpbody"
	"github.com/Eun/go-hit/internal/converter"
)

// IStoreBody defines the functions that can be used to store data from the http request/response body.
type IStoreBody interface {
	// Bool treats the body contents as Bool data and stores it
	//
	// Usage:
	//     var body bool
	//     Store().Response().Body().Bool().In(&body)
	Bool() IStoreStep

	// Bytes treats the body contents as byte data and stores it
	//
	// Usage:
	//     var body []byte
	//     Store().Response().Body().Bytes().In(&body)
	Bytes() IStoreStep

	// Float32 treats the body contents as Float32 data and stores it
	//
	// Usage:
	//     var body float32
	//     Store().Response().Body().Float32().In(&body)
	Float32() IStoreStep

	// Float64 treats the body contents as Float64 data and stores it
	//
	// Usage:
	//     var body float64
	//     Store().Response().Body().Float64().In(&body)
	Float64() IStoreStep

	// FormValues treats the body contents as FormValues data and stores it
	//
	// Usage:
	//     var values url.Values
	//     Store().Response().Body().FormValues().In(&values)
	//     var username string
	//     Store().Response().Body().FormValues("username").In(&username)
	FormValues(name ...string) IStoreStep

	// Int treats the body contents as Int data and stores it
	//
	// Usage:
	//     var body int
	//     Store().Response().Body().Int().In(&body)
	Int() IStoreStep

	// Int8 treats the body contents as Int8 data and stores it
	//
	// Usage:
	//     var body int8
	//     Store().Response().Body().Int8().In(&body)
	Int8() IStoreStep

	// Int16 treats the body contents as Int16 data and stores it
	//
	// Usage:
	//     var body int16
	//     Store().Response().Body().Int16().In(&body)
	Int16() IStoreStep

	// Int32 treats the body contents as Int32 data and stores it
	//
	// Usage:
	//     var body int32
	//     Store().Response().Body().Int32().In(&body)
	Int32() IStoreStep

	// Int64 treats the body contents as Int64 data and stores it
	//
	// Usage:
	//     var body int64
	//     Store().Response().Body().Int64().In(&body)
	Int64() IStoreStep

	// JSON treats the body as JSON data and stores it
	//
	// Example:
	//     // given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     var body map[string]interface{}
	//     var name string
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Store().Response().Body().JSON().In(&body),             // store the whole body as a map
	//         Store().Response().Body().JSON().JQ(".Name").In(&name), // store "Joe" in name
	//     )
	JSON() IStoreBodyJSON

	// Reader treats the body contents as Reader and stores it
	//
	// Usage:
	//     var body io.Reader
	//     Store().Response().Body().Reader().In(&body)
	Reader() IStoreStep

	// String treats the body contents as String data and stores it
	//
	// Usage:
	//     var body string
	//     Store().Response().Body().String().In(&body)
	String() IStoreStep

	// Uint treats the body contents as Uint data and stores it
	//
	// Usage:
	//     var body uint
	//     Store().Response().Body().Uint().In(&body)
	Uint() IStoreStep

	// Uint8 treats the body contents as Uint8 data and stores it
	//
	// Usage:
	//     var body uint8
	//     Store().Response().Body().Uint8().In(&body)
	Uint8() IStoreStep

	// Uint16 treats the body contents as Uint16 data and stores it
	//
	// Usage:
	//     var body uint16
	//     Store().Response().Body().Uint16().In(&body)
	Uint16() IStoreStep

	// Uint32 treats the body contents as Uint32 data and stores it
	//
	// Usage:
	//     var body uint32
	//     Store().Response().Body().Uint32().In(&body)
	Uint32() IStoreStep

	// Uint64 treats the body contents as Uint64 data and stores it
	//
	// Usage:
	//     var body uint64
	//     Store().Response().Body().Uint64().In(&body)
	Uint64() IStoreStep
}

type storeBodyMode int

const (
	storeBodyRequest storeBodyMode = iota
	storeBodyResponse
)

type storeBody struct {
	mode storeBodyMode
}

func newStoreBody(mode storeBodyMode) IStoreBody {
	return &storeBody{
		mode: mode,
	}
}

func (s *storeBody) body(hit Hit) *httpbody.HTTPBody {
	if s.mode == storeBodyRequest {
		return hit.Request().Body()
	}
	return hit.Response().Body()
}

func (s *storeBody) Bool() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustBool(), v)
	})
}

// Bytes returns the body as a byte slice.
func (s *storeBody) Bytes() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustBytes(), v)
	})
}

func (s *storeBody) Float32() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustFloat32(), v)
	})
}

func (s *storeBody) Float64() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustFloat64(), v)
	})
}

func (s *storeBody) FormValues(name ...string) IStoreStep {
	if n, ok := getLastStringArgument(name); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			return storeStringSlice(s.body(hit).MustFormValues().Values(n), v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustFormValues(), v)
	})
}

func (s *storeBody) Int() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustInt(), v)
	})
}

func (s *storeBody) Int8() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustInt8(), v)
	})
}

func (s *storeBody) Int16() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustInt16(), v)
	})
}

func (s *storeBody) Int32() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustInt32(), v)
	})
}

func (s *storeBody) Int64() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustInt64(), v)
	})
}

func (s *storeBody) JSON() IStoreBodyJSON {
	return newStoreBodyJSON(s.mode)
}

func (s *storeBody) Reader() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).Reader(), v)
	})
}

// String returns the body as a string.
func (s *storeBody) String() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustString(), v)
	})
}

func (s *storeBody) Uint() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustUint(), v)
	})
}

func (s *storeBody) Uint8() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustUint8(), v)
	})
}

func (s *storeBody) Uint16() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustUint16(), v)
	})
}

func (s *storeBody) Uint32() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustUint32(), v)
	})
}

func (s *storeBody) Uint64() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(s.body(hit).MustUint64(), v)
	})
}
