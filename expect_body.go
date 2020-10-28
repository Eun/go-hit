package hit

// IExpectBody provides assertions on the http response body.
type IExpectBody interface {
	// Bytes expects the body to be equal the specified byte slice.
	//
	// Usage:
	//     Expect().Body().Bytes().Equal([]byte("Hello World"))
	//     Expect().Body().Bytes().Contains([]byte("H"))
	Bytes() IExpectBytes

	// FormValues expects the body to be equal to the specified FormValues
	//
	// Usage:
	//     Expect().Body().FormValues("username").Equal("joe")
	FormValues(name string) IExpectFormValues

	// Float32 expects the body to be equal the specified int.
	//
	// Usage:
	//     Expect().Body().Float32().Equal(0.0)
	//     Expect().Body().Float32().GreaterThan(5.0)
	Float32() IExpectFloat32

	// Float64 expects the body to be equal the specified int.
	//
	// Usage:
	//     Expect().Body().Float64().Equal(0.0)
	//     Expect().Body().Float64().GreaterThan(5.0)
	Float64() IExpectFloat64

	// Int expects the body to be equal the specified int.
	//
	// Usage:
	//     Expect().Body().Int().Equal(0)
	//     Expect().Body().Int().GreaterThan(5)
	Int() IExpectInt

	// Int8 expects the body to be equal the specified int8.
	//
	// Usage:
	//     Expect().Body().Int8().Equal(0)
	//     Expect().Body().Int8().GreaterThan(5)
	Int8() IExpectInt8

	// Int16 expects the body to be equal the specified int16.
	//
	// Usage:
	//     Expect().Body().Int16().Equal(0)
	//     Expect().Body().Int16().GreaterThan(5)
	Int16() IExpectInt16

	// Int32 expects the body to be equal the specified int32.
	//
	// Usage:
	//     Expect().Body().Int32().Equal(0)
	//     Expect().Body().Int32().GreaterThan(5)
	Int32() IExpectInt32

	// Int64 expects the body to be equal the specified int64.
	//
	// Usage:
	//     Expect().Body().Int64().Equal(0)
	//     Expect().Body().Int64().GreaterThan(5)
	Int64() IExpectInt64

	// JSON provides assertions for the body in the JSON format
	//
	// Usage:
	//     Expect().Body().JSON().Equal([]int{1, 2, 3})
	//     Expect().Body().JSON().Contains(1)
	//     Expect().Body().JSON().JQ(".Name").Equal("Joe")
	JSON() IExpectBodyJSON

	// String expects the body to be equal the specified string.
	//
	// Usage:
	//     Expect().Body().String().Equal("Hello World")
	//     Expect().Body().String().Contains("Hello")
	String() IExpectString

	// Uint expects the body to be equal the specified uint.
	//
	// Usage:
	//     Expect().Body().Uint().Equal(0)
	//     Expect().Body().Uint().GreaterThan(5)
	Uint() IExpectUint

	// Uint8 expects the body to be equal the specified uint8.
	//
	// Usage:
	//     Expect().Body().Uint8().Equal(0)
	//     Expect().Body().Uint8().GreaterThan(5)
	Uint8() IExpectUint8

	// Uint16 expects the body to be equal the specified uint16.
	//
	// Usage:
	//     Expect().Body().Uint16().Equal(0)
	//     Expect().Body().Uint16().GreaterThan(5)
	Uint16() IExpectUint16

	// Uint32 expects the body to be equal the specified uint32.
	//
	// Usage:
	//     Expect().Body().Uint32().Equal(0)
	//     Expect().Body().Uint32().GreaterThan(5)
	Uint32() IExpectUint32

	// Uint64 expects the body to be equal the specified uint64.
	//
	// Usage:
	//     Expect().Body().Uint64().Equal(0)
	//     Expect().Body().Uint64().GreaterThan(5)
	Uint64() IExpectUint64
}

type expectBody struct {
	expect    IExpect
	cleanPath callPath
}

func newExpectBody(expect IExpect, cleanPath callPath) IExpectBody {
	return &expectBody{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (body *expectBody) Bytes() IExpectBytes {
	return newExpectBytes(body.cleanPath.Push("Bytes", nil), func(hit Hit) []byte {
		return hit.Response().Body().MustBytes()
	})
}

func (body *expectBody) FormValues(name string) IExpectFormValues {
	return newExpectFormValues(body.cleanPath.Push("FormValues", []interface{}{name}), func(hit Hit) []string {
		return hit.Response().Body().MustFormValues().Values(name)
	})
}

func (body *expectBody) Float32() IExpectFloat32 {
	return newExpectFloat32(body.cleanPath.Push("Float32", nil), func(hit Hit) float32 {
		return hit.Response().Body().MustFloat32()
	})
}

func (body *expectBody) Float64() IExpectFloat64 {
	return newExpectFloat64(body.cleanPath.Push("Float64", nil), func(hit Hit) float64 {
		return hit.Response().Body().MustFloat64()
	})
}

func (body *expectBody) Int() IExpectInt {
	return newExpectInt(body.cleanPath.Push("Int", nil), func(hit Hit) int {
		return hit.Response().Body().MustInt()
	})
}

func (body *expectBody) Int8() IExpectInt8 {
	return newExpectInt8(body.cleanPath.Push("Int8", nil), func(hit Hit) int8 {
		return hit.Response().Body().MustInt8()
	})
}

func (body *expectBody) Int16() IExpectInt16 {
	return newExpectInt16(body.cleanPath.Push("Int16", nil), func(hit Hit) int16 {
		return hit.Response().Body().MustInt16()
	})
}

func (body *expectBody) Int32() IExpectInt32 {
	return newExpectInt32(body.cleanPath.Push("Int32", nil), func(hit Hit) int32 {
		return hit.Response().Body().MustInt32()
	})
}

func (body *expectBody) Int64() IExpectInt64 {
	return newExpectInt64(body.cleanPath.Push("Int64", nil), func(hit Hit) int64 {
		return hit.Response().Body().MustInt64()
	})
}

func (body *expectBody) JSON() IExpectBodyJSON {
	return newExpectBodyJSON(body, body.cleanPath.Push("JSON", nil))
}

func (body *expectBody) String() IExpectString {
	return newExpectString(body.cleanPath.Push("String", nil), func(hit Hit) string {
		return hit.Response().Body().MustString()
	})
}

func (body *expectBody) Uint() IExpectUint {
	return newExpectUint(body.cleanPath.Push("Uint", nil), func(hit Hit) uint {
		return hit.Response().Body().MustUint()
	})
}

func (body *expectBody) Uint8() IExpectUint8 {
	return newExpectUint8(body.cleanPath.Push("Uint8", nil), func(hit Hit) uint8 {
		return hit.Response().Body().MustUint8()
	})
}

func (body *expectBody) Uint16() IExpectUint16 {
	return newExpectUint16(body.cleanPath.Push("Uint16", nil), func(hit Hit) uint16 {
		return hit.Response().Body().MustUint16()
	})
}

func (body *expectBody) Uint32() IExpectUint32 {
	return newExpectUint32(body.cleanPath.Push("Uint32", nil), func(hit Hit) uint32 {
		return hit.Response().Body().MustUint32()
	})
}

func (body *expectBody) Uint64() IExpectUint64 {
	return newExpectUint64(body.cleanPath.Push("Uint64", nil), func(hit Hit) uint64 {
		return hit.Response().Body().MustUint64()
	})
}
