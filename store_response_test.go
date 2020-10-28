package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"

	"net/http"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestStoreResponse_Status(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var v string
	Test(t,
		Post(s.URL),
		Store().Response().Status().In(&v),
	)

	require.Equal(t, fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK)), v)
}

func TestStoreResponse_StatusCode(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var v int
	Test(t,
		Post(s.URL),
		Store().Response().StatusCode().In(&v),
	)

	require.Equal(t, 200, v)
}

func TestStoreResponse_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var proto string
	Test(t,
		Post(s.URL),
		Store().Response().Proto().In(&proto),
	)

	require.Equal(t, `HTTP/1.1`, proto)
}

func TestStoreResponse_ProtoMajor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var protoMajor int
	Test(t,
		Post(s.URL),
		Store().Response().ProtoMajor().In(&protoMajor),
	)

	require.Equal(t, 1, protoMajor)
}

func TestStoreResponse_ProtoMinor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var protoMinor int
	Test(t,
		Post(s.URL),
		Store().Response().ProtoMinor().In(&protoMinor),
	)

	require.Equal(t, 1, protoMinor)
}

func TestStoreResponse_ContentLength(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var contentLength int
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Store().Response().ContentLength().In(&contentLength),
	)

	require.Equal(t, 11, contentLength)
}

func TestStoreResponse_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var transferEncoding []string
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().TransferEncoding = []string{"gzip", "chunked"}
		}),
		Store().Response().TransferEncoding().In(&transferEncoding),
	)

	require.Equal(t, []string{"gzip", "chunked"}, transferEncoding)
}

func TestStoreResponse_Uncompressed(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var uncompressed bool
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello"),
		Store().Response().Uncompressed().In(&uncompressed),
	)
	require.False(t, uncompressed)
}
