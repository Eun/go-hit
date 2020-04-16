package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"

	"net/url"

	"net/http"

	"github.com/Eun/go-hit/httpbody"
	"github.com/stretchr/testify/require"
)

func TestStoreRequest_Method(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var method string
	Test(t,
		Post(s.URL),
		Store().Request().Method().In(&method),
	)

	require.Equal(t, `POST`, method)
}

func TestStoreRequest_URL(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		var u url.URL
		Test(t,
			Post("http://joe:secret@%s", s.Listener.Addr().String()),
			Store().Request().URL().In(&u),
		)

		require.Equal(t, url.URL{
			Scheme:     "http",
			Opaque:     "",
			User:       url.UserPassword("joe", "secret"),
			Host:       s.Listener.Addr().String(),
			Path:       "",
			RawPath:    "",
			ForceQuery: false,
			RawQuery:   "",
			Fragment:   "",
		}, u)
	})

	t.Run("scheme", func(t *testing.T) {
		var scheme string
		Test(t,
			Post("http://joe:secret@%s", s.Listener.Addr().String()),
			Store().Request().URL().Scheme().In(&scheme),
		)

		require.Equal(t, "http", scheme)
	})

	t.Run("opaque", func(t *testing.T) {
		var opaque string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Opaque().In(&opaque),
		)

		require.Equal(t, "", opaque)
	})

	t.Run("user", func(t *testing.T) {
		t.Run("all", func(t *testing.T) {
			var u url.Userinfo
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().In(&u),
			)

			require.Equal(t, url.UserPassword("joe", "secret"), &u)
		})

		t.Run("username", func(t *testing.T) {
			var u string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().Username().In(&u),
			)

			require.Equal(t, "joe", u)
		})
		t.Run("password", func(t *testing.T) {
			var p string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().Password().In(&p),
			)

			require.Equal(t, "secret", p)
		})

		t.Run("string", func(t *testing.T) {
			var p string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().String().In(&p),
			)

			require.Equal(t, "joe:secret", p)
		})

		t.Run("nil password", func(t *testing.T) {
			var u string
			Test(t,
				Post(s.URL),
				Store().Request().URL().User().Password().In(&u),
			)

			require.Equal(t, "", u)
		})
	})

	t.Run("host", func(t *testing.T) {
		var host string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Host().In(&host),
		)

		require.Equal(t, s.Listener.Addr().String(), host)
	})

	t.Run("path", func(t *testing.T) {
		var path string
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().Path().In(&path),
		)

		require.Equal(t, "/foo/bar", path)
	})

	t.Run("RawPath", func(t *testing.T) {
		var rawPath string
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().RawPath().In(&rawPath),
		)

		require.Equal(t, "", rawPath)
	})

	t.Run("ForceQuery", func(t *testing.T) {
		var v bool
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().ForceQuery().In(&v),
		)

		require.False(t, v)
	})

	t.Run("RawQuery", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().RawQuery().In(&v),
		)

		require.Equal(t, "", v)
	})

	t.Run("Fragment", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://%s/foo/bar#hash", s.Listener.Addr().String()),
			Store().Request().URL().Fragment().In(&v),
		)

		require.Equal(t, "hash", v)
	})
}

func TestStoreRequest_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var proto string
	Test(t,
		Post(s.URL),
		Store().Request().Proto().In(&proto),
	)

	require.Equal(t, `HTTP/1.1`, proto)
}

func TestStoreRequest_ProtoMajor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var protoMajor int
	Test(t,
		Post(s.URL),
		Store().Request().ProtoMajor().In(&protoMajor),
	)

	require.Equal(t, 1, protoMajor)
}

func TestStoreRequest_ProtoMinor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var protoMinor int
	Test(t,
		Post(s.URL),
		Store().Request().ProtoMinor().In(&protoMinor),
	)

	require.Equal(t, 1, protoMinor)
}

func TestStoreRequest_ContentLength(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var contentLength int
	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Store().Request().ContentLength().In(&contentLength),
	)

	require.Equal(t, 11, contentLength)
}

func TestStoreRequest_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		var transferEncoding []string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Request().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Request().TransferEncoding().In(&transferEncoding),
		)

		require.Equal(t, []string{"gzip", "chunked"}, transferEncoding)
	})

	t.Run("specific", func(t *testing.T) {
		var transferEncoding string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Request().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Request().TransferEncoding("0").In(&transferEncoding),
		)

		require.Equal(t, "gzip", transferEncoding)
	})

	t.Run("out of bounds", func(t *testing.T) {
		var transferEncoding string
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Request().TransferEncoding = []string{"gzip", "chunked"}
			}),
			Store().Request().TransferEncoding("3").In(&transferEncoding),
		)

		require.Equal(t, "", transferEncoding)
	})
}

func TestStoreRequest_Host(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var host string

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Header("Host", "example.com"),
		Store().Request().Host().In(&host),
	)

	require.Equal(t, "example.com", host)
}

func TestStoreRequest_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		var str string

		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Store().Request().Body().In(&str),
		)

		require.Equal(t, `Hello World`, str)
	})

	t.Run("int", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Body("10"),
			Store().Request().Body().In(&n),
		)

		require.Equal(t, 10, n)
	})

	t.Run("httpbody", func(t *testing.T) {
		var body httpbody.HttpBody

		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Store().Request().Body().In(&body),
		)

		require.Equal(t, "Hello World", body.String())
	})
}

func TestStoreRequest_BodyJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("implicit", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Header("Content-Type", "application/json"),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Request().Body("Id").In(&n),
		)

		require.Equal(t, 12, n)
	})

	t.Run("string", func(t *testing.T) {
		var str string

		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Request().Body().JSON("Name").In(&str),
		)

		require.Equal(t, `Joe`, str)
	})

	t.Run("int", func(t *testing.T) {
		var n int

		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
			Store().Request().Body().JSON("Id").In(&n),
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
			Store().Request().Body().JSON().In(&m),
		)

		require.Equal(t, map[string]interface{}{"Name": "Joe", "Id": 12}, m)
	})
}

func TestStoreRequest_Header(t *testing.T) {
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
				Store().Request().Header().In(&headers),
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
				Store().Request().Header().In(&headers),
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
				Store().Request().Header("X-Header1").In(&hdr),
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
				Store().Request().Header("X-Header2").In(&hdr),
			)

			require.Equal(t, 10, hdr)
		})
	})
}

func TestStoreRequest_Trailer(t *testing.T) {
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
				Store().Request().Trailer().In(&trailers),
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
				Store().Request().Trailer().In(&trailers),
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
				Store().Request().Trailer("X-Trailer1").In(&trl),
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
				Store().Request().Trailer("X-Trailer2").In(&trl),
			)

			require.Equal(t, 10, trl)
		})
	})
}
