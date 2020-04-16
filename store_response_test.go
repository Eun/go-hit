package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"

	"net/http"

	"fmt"

	"github.com/Eun/go-hit/httpbody"
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
		Send().Body("Hello World"),
		Store().Response().ContentLength().In(&contentLength),
	)

	require.Equal(t, 11, contentLength)
}

func TestStoreResponse_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		var transferEncoding []string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Response().TransferEncoding().In(&transferEncoding),
		)

		require.Equal(t, []string{"gzip", "chunked"}, transferEncoding)
	})

	t.Run("specific", func(t *testing.T) {
		var transferEncoding string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Response().TransferEncoding("0").In(&transferEncoding),
		)

		require.Equal(t, "gzip", transferEncoding)
	})

	t.Run("out of bounds", func(t *testing.T) {
		var transferEncoding string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Response().TransferEncoding("3").In(&transferEncoding),
		)

		require.Equal(t, "", transferEncoding)
	})
}

func TestStoreResponse_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		var str string

		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Store().Response().Body().In(&str),
		)

		require.Equal(t, `Hello World`, str)
	})

	t.Run("int", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Body("10"),
			Store().Response().Body().In(&n),
		)

		require.Equal(t, 10, n)
	})

	t.Run("httpbody", func(t *testing.T) {
		var body httpbody.HttpBody

		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Store().Response().Body().In(&body),
		)

		require.Equal(t, "Hello World", body.String())
	})
}

func TestStoreResponse_BodyJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("implicit", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Header("Content-Type", "application/json"),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Response().Body("Id").In(&n),
		)

		require.Equal(t, 12, n)
	})

	t.Run("string", func(t *testing.T) {
		var str string

		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Response().Body().JSON("Name").In(&str),
		)

		require.Equal(t, `Joe`, str)
	})

	t.Run("int", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Response().Body().JSON("Id").In(&n),
		)

		require.Equal(t, 12, n)
	})

	t.Run("map", func(t *testing.T) {
		m := map[string]interface{}{
			"Name": "", // force name as string,
			"Id":   0,  // force id as int
		}

		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Response().Body().JSON().In(&m),
		)

		require.Equal(t, map[string]interface{}{"Name": "Joe", "Id": 12}, m)
	})
}

func TestStoreResponse_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all headers", func(t *testing.T) {
		t.Run("http.Header", func(t *testing.T) {
			var headers http.Header
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Header("X-Header1", "Foo"),
				Send().Header("X-Header2", "Bar"),
				Store().Response().Header().In(&headers),
			)

			require.Equal(t, "Foo", headers.Get("X-Header1"))
			require.Equal(t, "Bar", headers.Get("X-Header2"))
		})

		t.Run("map", func(t *testing.T) {
			var headers map[string]interface{}
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Header("X-Header1", "Foo"),
				Send().Header("X-Header2", "Bar"),
				Store().Response().Header().In(&headers),
			)

			require.Equal(t, "Foo", headers["X-Header1"])
			require.Equal(t, "Bar", headers["X-Header2"])
		})
	})

	t.Run("specific header", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			var hdr string
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Header("X-Header1", "Foo"),
				Send().Header("X-Header2", "10"),
				Store().Response().Header("X-Header1").In(&hdr),
			)

			require.Equal(t, "Foo", hdr)
		})

		t.Run("int", func(t *testing.T) {
			var hdr int
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Header("X-Header1", "Foo"),
				Send().Header("X-Header2", "10"),
				Store().Response().Header("X-Header2").In(&hdr),
			)

			require.Equal(t, 10, hdr)
		})
	})
}

func TestStoreResponse_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all headers", func(t *testing.T) {
		t.Run("http.Trailer", func(t *testing.T) {
			var trailers http.Header
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer1", "Foo"),
				Send().Trailer("X-Trailer2", "Bar"),
				Store().Response().Trailer().In(&trailers),
			)

			require.Equal(t, "Foo", trailers.Get("X-Trailer1"))
			require.Equal(t, "Bar", trailers.Get("X-Trailer2"))
		})

		t.Run("map", func(t *testing.T) {
			var trailers map[string]interface{}
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer1", "Foo"),
				Send().Trailer("X-Trailer2", "Bar"),
				Store().Response().Trailer().In(&trailers),
			)

			require.Equal(t, "Foo", trailers["X-Trailer1"])
			require.Equal(t, "Bar", trailers["X-Trailer2"])
		})
	})

	t.Run("specific header", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			var trl string
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer1", "Foo"),
				Send().Trailer("X-Trailer2", "10"),
				Store().Response().Trailer("X-Trailer1").In(&trl),
			)

			require.Equal(t, "Foo", trl)
		})

		t.Run("int", func(t *testing.T) {
			var trl int
			Test(t,
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer1", "Foo"),
				Send().Trailer("X-Trailer2", "10"),
				Store().Response().Trailer("X-Trailer2").In(&trl),
			)

			require.Equal(t, 10, trl)
		})
	})
}

func TestStoreResponse_Uncompressed(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var uncompressed bool
	Test(t,
		Post(s.URL),
		Send().Body("Hello"),
		Store().Response().Uncompressed().In(&uncompressed),
	)
	require.False(t, uncompressed)
}
