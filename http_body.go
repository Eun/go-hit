package hit

import (
	"fmt"
	"io"

	"io/ioutil"

	"bytes"

	"strconv"

	"strings"

	"github.com/Eun/go-doppelgangerreader"
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

type HTTPBody struct {
	hit     Hit
	factory doppelgangerreader.DoppelgangerFactory
}

func newHTTPBody(hit Hit, body io.Reader) *HTTPBody {
	b := &HTTPBody{
		hit: hit,
	}
	if body != nil {
		b.SetReader(body)
	}
	return b
}

// setters
// SetReader sets the body to the contents of the specified reader
func (body *HTTPBody) SetReader(r io.Reader) {
	if body.factory != nil {
		body.factory.Close()
		body.factory = nil
	}
	body.factory = doppelgangerreader.NewFactory(r)
}

// SetString sets the body to the specified byte slice
func (body *HTTPBody) SetBytes(b []byte) {
	body.SetReader(bytes.NewReader(b))
}

// SetString sets the body to the specified string
func (body *HTTPBody) SetString(s string) {
	body.SetBytes([]byte(s))
}

// SetStringf sets the body to the specified string format
func (body *HTTPBody) SetStringf(format string, a ...interface{}) {
	body.SetString(fmt.Sprintf(format, a...))
}

// SetInt sets the body to the specified int
func (body *HTTPBody) SetInt(i int) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt8 sets the body to the specified int8
func (body *HTTPBody) SetInt8(i int8) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt16 sets the body to the specified int16
func (body *HTTPBody) SetInt16(i int16) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt32 sets the body to the specified int64
func (body *HTTPBody) SetInt32(i int32) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt64 sets the body to the specified int64
func (body *HTTPBody) SetInt64(i int64) {
	body.SetString(strconv.FormatInt(i, 10))
}

// SetUint sets the body to the specified int
func (body *HTTPBody) SetUint(i uint) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint8 sets the body to the specified uint8
func (body *HTTPBody) SetUint8(i uint8) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint16 sets the body to the specified uint16
func (body *HTTPBody) SetUint16(i uint16) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint32 sets the body to the specified uint64
func (body *HTTPBody) SetUint32(i uint32) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint64 sets the body to the specified uint64
func (body *HTTPBody) SetUint64(i uint64) {
	body.SetString(strconv.FormatUint(i, 10))
}

// SetFloat32 sets the body to the specified float32
func (body *HTTPBody) SetFloat32(f float32) {
	body.SetString(strconv.FormatFloat(float64(f), 'f', 6, 32))
}

// SetFloat64 sets the body to the specified float64
func (body *HTTPBody) SetFloat64(f float64) {
	body.SetString(strconv.FormatFloat(f, 'f', 6, 64))
}

// SetBool sets the body to the specified bool
func (body *HTTPBody) SetBool(b bool) {
	body.SetString(strconv.FormatBool(b))
}

// Set sets the body to the specified value
func (body *HTTPBody) Set(data interface{}) {
	if body.setOnlyNativeTypes(data) {
		return
	}

	body.JSON().Set(data)
}

// getters
// Reader returns the body as an reader
func (body *HTTPBody) Reader() io.ReadCloser {
	if body.factory == nil {
		return nil
	}
	return body.factory.NewDoppelganger()
}

// Bytes returns the body as a byte slice
func (body *HTTPBody) Bytes() []byte {
	r := body.Reader()
	if r == nil {
		return nil
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		minitest.NoError(err, "failed to read body")
	}
	r.Close()
	return buf
}

// String returns the body as a string
func (body *HTTPBody) String() string {
	return string(body.Bytes())
}

// Int returns the body as an int
func (body *HTTPBody) Int() int {
	n, err := strconv.ParseInt(body.String(), 0, 0)
	minitest.NoError(err)
	return int(n)
}

// Int8 returns the body as an int8
func (body *HTTPBody) Int8() int8 {
	n, err := strconv.ParseInt(body.String(), 0, 8)
	minitest.NoError(err)
	return int8(n)
}

// Int16 returns the body as an int16
func (body *HTTPBody) Int16() int16 {
	n, err := strconv.ParseInt(body.String(), 0, 16)
	minitest.NoError(err)
	return int16(n)
}

// Int32 returns the body as an int64
func (body *HTTPBody) Int32() int32 {
	n, err := strconv.ParseInt(body.String(), 0, 32)
	minitest.NoError(err)
	return int32(n)
}

// Int64 returns the body as an int64
func (body *HTTPBody) Int64() int64 {
	n, err := strconv.ParseInt(body.String(), 0, 64)
	minitest.NoError(err)
	return n
}

// Uint returns the body as an uint
func (body *HTTPBody) Uint() uint {
	n, err := strconv.ParseUint(body.String(), 0, 0)
	minitest.NoError(err)
	return uint(n)
}

// Uint8 returns the body as an uint8
func (body *HTTPBody) Uint8() uint8 {
	n, err := strconv.ParseUint(body.String(), 0, 8)
	minitest.NoError(err)
	return uint8(n)
}

// Uint16 returns the body as an uint16
func (body *HTTPBody) Uint16() uint16 {
	n, err := strconv.ParseUint(body.String(), 0, 16)
	minitest.NoError(err)
	return uint16(n)
}

// Uint32 returns the body as an uint64
func (body *HTTPBody) Uint32() uint32 {
	n, err := strconv.ParseUint(body.String(), 0, 32)
	minitest.NoError(err)
	return uint32(n)
}

// Uint64 returns the body as an uint64
func (body *HTTPBody) Uint64() uint64 {
	n, err := strconv.ParseUint(body.String(), 0, 64)
	minitest.NoError(err)
	return n
}

// Float32 returns the body as an float32
func (body *HTTPBody) Float32() float32 {
	n, err := strconv.ParseFloat(body.String(), 32)
	minitest.NoError(err)
	return float32(n)
}

// Float64 returns the body as an float64
func (body *HTTPBody) Float64() float64 {
	n, err := strconv.ParseFloat(body.String(), 64)
	minitest.NoError(err)
	return n
}

// Bool returns the body as an bool
func (body *HTTPBody) Bool() bool {
	n, err := strconv.ParseBool(body.String())
	minitest.NoError(err)
	return n
}

// JSON returns the body as an jso
func (body *HTTPBody) JSON() *HTTPJson {
	return newHTTPJson(body)
}

// Length returns the body's length
func (body *HTTPBody) Length() int64 {
	if body.factory == nil {
		return 0
	}

	r := body.factory.NewDoppelganger()
	n, _ := io.Copy(internal.DevNullWriter(), r)
	_ = r.Close()
	return n
}

func (body *HTTPBody) setOnlyNativeTypes(a interface{}) bool {
	switch v := a.(type) {
	case string:
		body.SetString(v)
	case []byte:
		body.SetBytes(v)
	case io.Reader:
		body.SetReader(v)
	case int:
		body.SetInt(v)
	case int8:
		body.SetInt8(v)
	case int16:
		body.SetInt16(v)
	case int32:
		body.SetInt32(v)
	case int64:
		body.SetInt64(v)
	case uint:
		body.SetUint(v)
	case uint8:
		body.SetUint8(v)
	case uint16:
		body.SetUint16(v)
	case uint32:
		body.SetUint32(v)
	case uint64:
		body.SetUint64(v)
	case float32:
		body.SetFloat32(v)
	case float64:
		body.SetFloat64(v)
	case bool:
		body.SetBool(v)
	default:
		return false
	}
	return true
}

func (body *HTTPBody) equalOnlyNativeTypes(a interface{}, equal bool) bool {
	equalFunc := minitest.Equal
	if !equal {
		equalFunc = minitest.NotEqual
	}
	switch v := a.(type) {
	case string:
		equalFunc(v, body.String())
	case []byte:
		equalFunc(v, body.Bytes())
	case io.Reader:
		buf, err := ioutil.ReadAll(v)
		if err != nil {
			minitest.Errorf("unable to read data from reader: %s", err.Error())
		}
		equalFunc(buf, body.Bytes())
	case int:
		equalFunc(v, body.Int())
	case int8:
		equalFunc(v, body.Int8())
	case int16:
		equalFunc(v, body.Int16())
	case int32:
		equalFunc(v, body.Int32())
	case int64:
		equalFunc(v, body.Int64())
	case uint:
		equalFunc(v, body.Uint())
	case uint8:
		equalFunc(v, body.Uint8())
	case uint16:
		equalFunc(v, body.Uint16())
	case uint32:
		equalFunc(v, body.Uint32())
	case uint64:
		equalFunc(v, body.Uint64())
	case float32:
		equalFunc(v, body.Float32())
	case float64:
		equalFunc(v, body.Float64())
	case bool:
		equalFunc(v, body.Bool())
	default:
		// not handled
		return false
	}
	// handled and values are the same
	return true
}

func (body *HTTPBody) containsOnlyNativeTypes(a interface{}, equal bool) bool {
	switch v := a.(type) {
	case string:
		switch equal {
		case true:
			if !strings.Contains(body.String(), v) {
				minitest.Errorf(`"%s" does not contain "%s"`, body.String(), v)
			}
		default:
			if strings.Contains(body.String(), v) {
				minitest.Errorf(`"%s" does contain "%s"`, body.String(), v)
			}
		}

	case []byte:
		switch equal {
		case true:
			if !bytes.Contains(body.Bytes(), v) {
				minitest.Errorf(`"%v" does not contain "%v"`, body.Bytes(), v)
			}
		default:
			if bytes.Contains(body.Bytes(), v) {
				minitest.Errorf(`"%v" does contain "%v"`, body.Bytes(), v)
			}
		}

	case io.Reader:
		buf, err := ioutil.ReadAll(v)
		if err != nil {
			minitest.Errorf("unable to read data from reader: %s", err.Error())
		}
		return body.containsOnlyNativeTypes(buf, equal)
	default:
		// not handled
		return false
	}
	// handled and values are the same
	return true
}
