package httpbody

import (
	"fmt"
	"io"
	"strings"

	"io/ioutil"

	"bytes"

	"strconv"

	"net/http"

	"github.com/Eun/go-doppelgangerreader"
	"github.com/Eun/go-hit/internal/minitest"
	"github.com/Eun/go-hit/internal/misc"
)

type HttpBody struct {
	factory doppelgangerreader.DoppelgangerFactory
	headers http.Header
}

func NewHttpBody(body io.Reader, headers http.Header) *HttpBody {
	b := &HttpBody{
		headers: headers,
	}
	if body != nil {
		b.SetReader(body)
	}
	return b
}

// setters
// SetReader sets the body to the contents of the specified reader
func (body *HttpBody) SetReader(r io.Reader) {
	if body.factory != nil {
		body.factory.Close()
		body.factory = nil
	}
	body.factory = doppelgangerreader.NewFactory(r)
}

// SetString sets the body to the specified byte slice
func (body *HttpBody) SetBytes(b []byte) {
	body.SetReader(bytes.NewReader(b))
}

// SetString sets the body to the specified string
func (body *HttpBody) SetString(s string) {
	body.SetBytes([]byte(s))
}

// SetStringf sets the body to the specified string format
func (body *HttpBody) SetStringf(format string, a ...interface{}) {
	body.SetString(fmt.Sprintf(format, a...))
}

// SetInt sets the body to the specified int
func (body *HttpBody) SetInt(i int) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt8 sets the body to the specified int8
func (body *HttpBody) SetInt8(i int8) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt16 sets the body to the specified int16
func (body *HttpBody) SetInt16(i int16) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt32 sets the body to the specified int64
func (body *HttpBody) SetInt32(i int32) {
	body.SetString(strconv.FormatInt(int64(i), 10))
}

// SetInt64 sets the body to the specified int64
func (body *HttpBody) SetInt64(i int64) {
	body.SetString(strconv.FormatInt(i, 10))
}

// SetUint sets the body to the specified int
func (body *HttpBody) SetUint(i uint) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint8 sets the body to the specified uint8
func (body *HttpBody) SetUint8(i uint8) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint16 sets the body to the specified uint16
func (body *HttpBody) SetUint16(i uint16) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint32 sets the body to the specified uint64
func (body *HttpBody) SetUint32(i uint32) {
	body.SetString(strconv.FormatUint(uint64(i), 10))
}

// SetUint64 sets the body to the specified uint64
func (body *HttpBody) SetUint64(i uint64) {
	body.SetString(strconv.FormatUint(i, 10))
}

// SetFloat32 sets the body to the specified float32
func (body *HttpBody) SetFloat32(f float32) {
	body.SetString(strconv.FormatFloat(float64(f), 'f', 6, 32))
}

// SetFloat64 sets the body to the specified float64
func (body *HttpBody) SetFloat64(f float64) {
	body.SetString(strconv.FormatFloat(f, 'f', 6, 64))
}

// SetBool sets the body to the specified bool
func (body *HttpBody) SetBool(b bool) {
	body.SetString(strconv.FormatBool(b))
}

// getters
// Reader returns the body as an reader
func (body *HttpBody) Reader() io.ReadCloser {
	if body.factory == nil {
		return nil
	}
	return body.factory.NewDoppelganger()
}

// Bytes returns the body as a byte slice
func (body *HttpBody) Bytes() []byte {
	r := body.Reader()
	if r == nil {
		return nil
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		minitest.Panic.NoError(err, "failed to read body")
	}
	r.Close()
	return buf
}

// String returns the body as a string
func (body *HttpBody) String() string {
	return string(body.Bytes())
}

// Int returns the body as an int
func (body *HttpBody) Int() int {
	n, err := strconv.ParseInt(body.String(), 0, 0)
	minitest.Panic.NoError(err)
	return int(n)
}

// Int8 returns the body as an int8
func (body *HttpBody) Int8() int8 {
	n, err := strconv.ParseInt(body.String(), 0, 8)
	minitest.Panic.NoError(err)
	return int8(n)
}

// Int16 returns the body as an int16
func (body *HttpBody) Int16() int16 {
	n, err := strconv.ParseInt(body.String(), 0, 16)
	minitest.Panic.NoError(err)
	return int16(n)
}

// Int32 returns the body as an int64
func (body *HttpBody) Int32() int32 {
	n, err := strconv.ParseInt(body.String(), 0, 32)
	minitest.Panic.NoError(err)
	return int32(n)
}

// Int64 returns the body as an int64
func (body *HttpBody) Int64() int64 {
	n, err := strconv.ParseInt(body.String(), 0, 64)
	minitest.Panic.NoError(err)
	return n
}

// Uint returns the body as an uint
func (body *HttpBody) Uint() uint {
	n, err := strconv.ParseUint(body.String(), 0, 0)
	minitest.Panic.NoError(err)
	return uint(n)
}

// Uint8 returns the body as an uint8
func (body *HttpBody) Uint8() uint8 {
	n, err := strconv.ParseUint(body.String(), 0, 8)
	minitest.Panic.NoError(err)
	return uint8(n)
}

// Uint16 returns the body as an uint16
func (body *HttpBody) Uint16() uint16 {
	n, err := strconv.ParseUint(body.String(), 0, 16)
	minitest.Panic.NoError(err)
	return uint16(n)
}

// Uint32 returns the body as an uint64
func (body *HttpBody) Uint32() uint32 {
	n, err := strconv.ParseUint(body.String(), 0, 32)
	minitest.Panic.NoError(err)
	return uint32(n)
}

// Uint64 returns the body as an uint64
func (body *HttpBody) Uint64() uint64 {
	n, err := strconv.ParseUint(body.String(), 0, 64)
	minitest.Panic.NoError(err)
	return n
}

// Float32 returns the body as an float32
func (body *HttpBody) Float32() float32 {
	n, err := strconv.ParseFloat(body.String(), 32)
	minitest.Panic.NoError(err)
	return float32(n)
}

// Float64 returns the body as an float64
func (body *HttpBody) Float64() float64 {
	n, err := strconv.ParseFloat(body.String(), 64)
	minitest.Panic.NoError(err)
	return n
}

// Bool returns the body as an bool
func (body *HttpBody) Bool() bool {
	n, err := strconv.ParseBool(body.String())
	minitest.Panic.NoError(err)
	return n
}

// JSON treats the body as JSON encoded data
func (body *HttpBody) JSON() *HttpBodyJson {
	return newHttpBodyJson(body)
}

// XML treats the body as XML encoded data
func (body *HttpBody) XML() *HttpBodyXml {
	return newHttpBodyXml(body)
}

// GetBestFittingObject tries its best to return an appropriate object for the body
func (body *HttpBody) GetBestFittingObject() interface{} {
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

// Length returns the body's length
func (body *HttpBody) Length() int64 {
	if body.factory == nil {
		return 0
	}

	r := body.factory.NewDoppelganger()
	n, _ := io.Copy(misc.DevNullWriter(), r)
	_ = r.Close()
	return n
}

func (body *HttpBody) ContainsOnlyNativeTypes(a interface{}, equal bool) (handled bool, err error) {
	switch v := a.(type) {
	case string:
		switch equal {
		case true:
			if !strings.Contains(body.String(), v) {
				return true, minitest.Error.Errorf(`"%s" does not contain "%s"`, body.String(), v)
			}
			return true, nil
		default:
			if strings.Contains(body.String(), v) {
				return true, minitest.Error.Errorf(`"%s" does contain "%s"`, body.String(), v)
			}
			return true, nil
		}

	case []byte:
		switch equal {
		case true:
			if !bytes.Contains(body.Bytes(), v) {
				return true, minitest.Error.Errorf(`"%v" does not contain "%v"`, body.Bytes(), v)
			}
			return true, nil
		default:
			if bytes.Contains(body.Bytes(), v) {
				return true, minitest.Error.Errorf(`"%v" does contain "%v"`, body.Bytes(), v)
			}
			return true, nil
		}

	case io.Reader:
		buf, err := ioutil.ReadAll(v)
		if err != nil {
			return true, minitest.Error.Errorf("unable to read data from reader: %s", err.Error())
		}
		return body.ContainsOnlyNativeTypes(buf, equal)
	default:
		// not handled
		return false, nil
	}
}

func (body *HttpBody) EqualOnlyNativeTypes(a interface{}, equal bool) (handled bool, err error) {
	equalFunc := minitest.Error.Equal
	if !equal {
		equalFunc = minitest.Error.NotEqual
	}
	switch v := a.(type) {
	case string:
		return true, equalFunc(v, body.String())
	case []byte:
		return true, equalFunc(v, body.Bytes())
	case io.Reader:
		buf, err := ioutil.ReadAll(v)
		if err != nil {
			return true, minitest.Error.Errorf("unable to read data from reader: %s", err.Error())
		}
		return true, equalFunc(buf, body.Bytes())
	case int:
		return true, equalFunc(v, body.Int())
	case int8:
		return true, equalFunc(v, body.Int8())
	case int16:
		return true, equalFunc(v, body.Int16())
	case int32:
		return true, equalFunc(v, body.Int32())
	case int64:
		return true, equalFunc(v, body.Int64())
	case uint:
		return true, equalFunc(v, body.Uint())
	case uint8:
		return true, equalFunc(v, body.Uint8())
	case uint16:
		return true, equalFunc(v, body.Uint16())
	case uint32:
		return true, equalFunc(v, body.Uint32())
	case uint64:
		return true, equalFunc(v, body.Uint64())
	case float32:
		return true, equalFunc(v, body.Float32())
	case float64:
		return true, equalFunc(v, body.Float64())
	case bool:
		return true, equalFunc(v, body.Bool())
	default:
		return false, nil
	}
}
