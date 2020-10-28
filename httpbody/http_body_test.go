package httpbody

import (
	"bytes"
	"context"
	"io"
	"testing"

	"io/ioutil"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func testServer(response string, fn func(body *HTTPBody)) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		io.WriteString(writer, response)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.URL, nil)
	if err != nil {
		panic(err)
	}
	res, err := s.Client().Do(req)
	if err != nil {
		panic(err)
	}
	fn(NewHTTPBody(res.Body, res.Header))
	if res.Body != nil {
		res.Body.Close()
	}
}

func TestHTTPBody_SetGet(t *testing.T) {
	t.Run("Reader", func(t *testing.T) {
		testServer("Some Old Data", func(body *HTTPBody) {
			buf, err := ioutil.ReadAll(body.Reader())
			require.NoError(t, err)
			require.Equal(t, []byte("Some Old Data"), buf)
			body.SetReader(bytes.NewReader([]byte("Hello World")))

			buf, err = ioutil.ReadAll(body.Reader())
			require.NoError(t, err)
			require.Equal(t, []byte("Hello World"), buf)
		})
	})

	t.Run("Bytes", func(t *testing.T) {
		testServer("Some Old Data", func(body *HTTPBody) {
			require.Equal(t, []byte("Some Old Data"), body.MustBytes())
			body.SetBytes([]byte("Hello World"))
			require.Equal(t, []byte("Hello World"), body.MustBytes())
		})
	})

	t.Run("String", func(t *testing.T) {
		testServer("Some Old Data", func(body *HTTPBody) {
			require.Equal(t, "Some Old Data", body.MustString())
			body.SetString("Hello World")
			require.Equal(t, "Hello World", body.MustString())
		})
	})

	t.Run("Stringf", func(t *testing.T) {
		testServer("Some Old Data", func(body *HTTPBody) {
			require.Equal(t, "Some Old Data", body.MustString())
			body.SetStringf("Hello %s", "World")
			require.Equal(t, "Hello World", body.MustString())
		})
	})

	t.Run("Int", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, 5, body.MustInt())
			body.SetInt(10)
			require.Equal(t, 10, body.MustInt())
		})
	})

	t.Run("Int8", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, int8(5), body.MustInt8())
			body.SetInt8(10)
			require.Equal(t, int8(10), body.MustInt8())
		})
	})

	t.Run("Int16", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, int16(5), body.MustInt16())
			body.SetInt16(10)
			require.Equal(t, int16(10), body.MustInt16())
		})
	})

	t.Run("Int32", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, int32(5), body.MustInt32())
			body.SetInt32(10)
			require.Equal(t, int32(10), body.MustInt32())
		})
	})

	t.Run("Int64", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, int64(5), body.MustInt64())
			body.SetInt64(10)
			require.Equal(t, int64(10), body.MustInt64())
		})
	})

	t.Run("Uint", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, uint(5), body.MustUint())
			body.SetUint(10)
			require.Equal(t, uint(10), body.MustUint())
		})
	})

	t.Run("Uint8", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, uint8(5), body.MustUint8())
			body.SetUint8(10)
			require.Equal(t, uint8(10), body.MustUint8())
		})
	})

	t.Run("Uint16", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, uint16(5), body.MustUint16())
			body.SetUint16(10)
			require.Equal(t, uint16(10), body.MustUint16())
		})
	})

	t.Run("Uint32", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, uint32(5), body.MustUint32())
			body.SetUint32(10)
			require.Equal(t, uint32(10), body.MustUint32())
		})
	})

	t.Run("Uint64", func(t *testing.T) {
		testServer("5", func(body *HTTPBody) {
			require.Equal(t, uint64(5), body.MustUint64())
			body.SetUint64(10)
			require.Equal(t, uint64(10), body.MustUint64())
		})
	})

	t.Run("Float32", func(t *testing.T) {
		testServer("5.0", func(body *HTTPBody) {
			require.Equal(t, float32(5.0), body.MustFloat32())
			body.SetFloat32(10.0)
			require.Equal(t, float32(10.0), body.MustFloat32())
		})
	})

	t.Run("Float64", func(t *testing.T) {
		testServer("5.0", func(body *HTTPBody) {
			require.Equal(t, float64(5.0), body.MustFloat64())
			body.SetFloat64(10.0)
			require.Equal(t, float64(10.0), body.MustFloat64())
		})
	})

	t.Run("Bool", func(t *testing.T) {
		testServer("false", func(body *HTTPBody) {
			require.Equal(t, false, body.MustBool())
			body.SetBool(true)
			require.Equal(t, true, body.MustBool())
		})
	})
}

func TestHttpBody_Length(t *testing.T) {
	tests := []struct {
		response string
		want     int64
	}{
		{
			"",
			0,
		},
		{
			"Hello World",
			11,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			testServer(tt.response, func(body *HTTPBody) {
				got := body.MustLength()
				if got != tt.want {
					t.Errorf("Length() got = %v, want %v", got, tt.want)
				}
			})
		})
	}
}
