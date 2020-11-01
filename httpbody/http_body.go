// Package httpbody contains a http body representation with a reusable body that can be consumed multiple times.
// It also provides setters and getters for special body formats such as JSON or XML.
package httpbody

import (
	"fmt"
	"io"
	"strings"

	"io/ioutil"

	"bytes"

	"strconv"

	"net/http"

	"net/url"

	"github.com/Eun/go-convert"
	"github.com/Eun/go-doppelgangerreader"
)

// HTTPBody provides setters and getters for the http body.
type HTTPBody struct {
	factory doppelgangerreader.DoppelgangerFactory
	headers http.Header
}

// NewHTTPBody creates a new HTTPBody from the specified reader and headers.
func NewHTTPBody(body io.Reader, headers http.Header) *HTTPBody {
	b := &HTTPBody{
		headers: headers,
	}
	if body != nil {
		b.SetReader(body)
	}
	return b
}

// setters

// SetReader sets the body to the contents of the specified reader.
func (body *HTTPBody) SetReader(r io.Reader) {
	if body.factory != nil {
		_ = body.factory.Close()
		body.factory = nil
	}
	body.factory = doppelgangerreader.NewFactory(r)
}

// SetBytes sets the body to the specified byte slice.
func (body *HTTPBody) SetBytes(b []byte) {
	body.SetReader(bytes.NewReader(b))
}

// SetString sets the body to the specified string.
func (body *HTTPBody) SetString(s string) {
	body.SetBytes([]byte(s))
}

// SetStringf sets the body to the specified string format.
func (body *HTTPBody) SetStringf(format string, a ...interface{}) {
	body.SetString(fmt.Sprintf(format, a...))
}

// SetInt sets the body to the specified int.
func (body *HTTPBody) SetInt(i int) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt8 sets the body to the specified int8.
func (body *HTTPBody) SetInt8(i int8) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt16 sets the body to the specified int16.
func (body *HTTPBody) SetInt16(i int16) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt32 sets the body to the specified int64.
func (body *HTTPBody) SetInt32(i int32) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt64 sets the body to the specified int64.
func (body *HTTPBody) SetInt64(i int64) {
	body.SetString(strconv.FormatInt(i, 10))
}

// SetUint sets the body to the specified int.
func (body *HTTPBody) SetUint(i uint) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint8 sets the body to the specified uint8.
func (body *HTTPBody) SetUint8(i uint8) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint16 sets the body to the specified uint16.
func (body *HTTPBody) SetUint16(i uint16) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint32 sets the body to the specified uint64.
func (body *HTTPBody) SetUint32(i uint32) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint64 sets the body to the specified uint64.
func (body *HTTPBody) SetUint64(i uint64) {
	body.SetString(strconv.FormatUint(i, 10))
}

// SetFloat32 sets the body to the specified float32.
func (body *HTTPBody) SetFloat32(f float32) {
	body.SetString(strconv.FormatFloat(float64(f), 'f', 6, 32))
}

// SetFloat64 sets the body to the specified float64.
func (body *HTTPBody) SetFloat64(f float64) {
	body.SetString(strconv.FormatFloat(f, 'f', 6, 64))
}

// SetFormValues sets the body to the specified form values.
func (body *HTTPBody) SetFormValues(v url.Values) {
	body.SetString(v.Encode())
}

// SetBool sets the body to the specified bool.
func (body *HTTPBody) SetBool(b bool) {
	body.SetString(strconv.FormatBool(b))
}

// getters

// Reader returns the body as an reader.
func (body *HTTPBody) Reader() io.ReadCloser {
	if body.factory == nil {
		return nil
	}
	return body.factory.NewDoppelganger()
}

// Bytes returns the body as a byte slice.
func (body *HTTPBody) Bytes() ([]byte, error) {
	r := body.Reader()
	if r == nil {
		return nil, nil
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return buf, r.Close()
}

// MustBytes returns the body as a byte slice, it panics on failure.
func (body *HTTPBody) MustBytes() []byte {
	buf, err := body.Bytes()
	if err != nil {
		panic(err)
	}
	return buf
}

// String returns the body as a string.
func (body *HTTPBody) String() (string, error) {
	buf, err := body.Bytes()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// MustString returns the body as a string, it panics on failure.
func (body *HTTPBody) MustString() string {
	buf, err := body.String()
	if err != nil {
		panic(err)
	}
	return buf
}

// Int returns the body as an int.
func (body *HTTPBody) Int() (int, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// MustInt returns the body as an int, it panics on failure.
func (body *HTTPBody) MustInt() int {
	v, err := body.Int()
	if err != nil {
		panic(err)
	}
	return v
}

// Int8 returns the body as an int8.
func (body *HTTPBody) Int8() (int8, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		panic(err)
	}
	return int8(n), nil
}

// MustInt8 returns the body as an int8, it panics on failure.
func (body *HTTPBody) MustInt8() int8 {
	v, err := body.Int8()
	if err != nil {
		panic(err)
	}
	return v
}

// Int16 returns the body as an int16.
func (body *HTTPBody) Int16() (int16, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}

// MustInt16 returns the body as an int16, it panics on failure.
func (body *HTTPBody) MustInt16() int16 {
	v, err := body.Int16()
	if err != nil {
		panic(err)
	}
	return v
}

// Int32 returns the body as an int64.
func (body *HTTPBody) Int32() (int32, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

// MustInt32 returns the body as an int64, it panics on failure.
func (body *HTTPBody) MustInt32() int32 {
	v, err := body.Int32()
	if err != nil {
		panic(err)
	}
	return v
}

// Int64 returns the body as an int64.
func (body *HTTPBody) Int64() (int64, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// MustInt64 returns the body as an int64, it panics on failure.
func (body *HTTPBody) MustInt64() int64 {
	v, err := body.Int64()
	if err != nil {
		panic(err)
	}
	return v
}

// Uint returns the body as an uint.
func (body *HTTPBody) Uint() (uint, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(s, 0, 0)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}

// MustUint returns the body as an uint, it panics on failure.
func (body *HTTPBody) MustUint() uint {
	v, err := body.Uint()
	if err != nil {
		panic(err)
	}
	return v
}

// Uint8 returns the body as an uint8.
func (body *HTTPBody) Uint8() (uint8, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}

// MustUint8 returns the body as an uint8, it panics on failure.
func (body *HTTPBody) MustUint8() uint8 {
	v, err := body.Uint8()
	if err != nil {
		panic(err)
	}
	return (v)
}

// Uint16 returns the body as an uint16.
func (body *HTTPBody) Uint16() (uint16, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

// MustUint16 returns the body as an uint16, it panics on failure.
func (body *HTTPBody) MustUint16() uint16 {
	v, err := body.Uint16()
	if err != nil {
		panic(err)
	}
	return (v)
}

// Uint32 returns the body as an uint32.
func (body *HTTPBody) Uint32() (uint32, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}

// MustUint32 returns the body as an uint32, it panics on failure.
func (body *HTTPBody) MustUint32() uint32 {
	v, err := body.Uint32()
	if err != nil {
		panic(err)
	}
	return (v)
}

// Uint64 returns the body as an uint64.
func (body *HTTPBody) Uint64() (uint64, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// MustUint64 returns the body as an uint64, it panics on failure.
func (body *HTTPBody) MustUint64() uint64 {
	v, err := body.Uint64()
	if err != nil {
		panic(err)
	}
	return v
}

// Float32 returns the body as an float32.
func (body *HTTPBody) Float32() (float32, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(n), nil
}

// MustFloat32 returns the body as an float32, it panics on failure.
func (body *HTTPBody) MustFloat32() float32 {
	v, err := body.Float32()
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns the body as an float64.
func (body *HTTPBody) Float64() (float64, error) {
	s, err := body.String()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

// MustFloat64 returns the body as an float64, it panics on failure.
func (body *HTTPBody) MustFloat64() float64 {
	v, err := body.Float64()
	if err != nil {
		panic(err)
	}
	return v
}

// FormValues returns the body as url.Values.
func (body *HTTPBody) FormValues() (*URLValues, error) {
	return ParseURLValues(body)
}

// MustFormValues returns the body as url.FormValues.
func (body *HTTPBody) MustFormValues() *URLValues {
	v, err := body.FormValues()
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns the body as an bool.
func (body *HTTPBody) Bool() (bool, error) {
	s, err := body.String()
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}

// MustBool returns the body as an bool, it panics on failure.
func (body *HTTPBody) MustBool() bool {
	v, err := body.Bool()
	if err != nil {
		panic(err)
	}
	return v
}

// JSON treats the body as JSON encoded data.
func (body *HTTPBody) JSON() *HTTPBodyJSON {
	return newHTTPBodyJSON(body.Reader, body.SetBytes)
}

// XML treats the body as XML encoded data.
func (body *HTTPBody) XML() *HTTPBodyXML {
	return newHTTPBodyXML(body)
}

// GetBestFittingObject tries its best to return an appropriate object for the body.
func (body *HTTPBody) GetBestFittingObject() interface{} {
	var decoder *decoder
	switch strings.ToLower(body.headers.Get("Content-Type")) {
	case "application/json", "text/json":
		decoder = newDecoder(jsonDecoder())
	// case "application/xml", "text/xml":
	// 	decoder = newDecoder(xmlDecoder())
	default:
		decoder = newDecoder()
	}

	if result := decoder.Run(body); result != nil {
		return result
	}
	return nil
}

// Length returns the body's length.
func (body *HTTPBody) Length() (int64, error) {
	if body.factory == nil {
		return 0, nil
	}

	r := body.factory.NewDoppelganger()
	n, err := io.Copy(ioutil.Discard, r)
	if err != nil {
		return 0, err
	}
	return n, r.Close()
}

// MustLength returns the body's length, it panics on failure.
func (body *HTTPBody) MustLength() int64 {
	n, err := body.Length()
	if err != nil {
		panic(err)
	}
	return n
}

// ConvertRecipes contains recipes for go-convert.
func (body *HTTPBody) ConvertRecipes() []convert.Recipe {
	return convert.MustMakeRecipes(func(_ convert.Converter, in HTTPBody, out *HTTPBody) error {
		// just copy
		*out = in
		return nil
	})
}
