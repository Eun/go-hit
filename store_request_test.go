package hit_test

import (
	"fmt"
	"net"
	"testing"

	. "github.com/otto-eng/go-hit"

	"net/url"

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
		var v url.URL
		Test(t,
			Post("http://joe:secret@%s", s.Listener.Addr().String()),
			Store().Request().URL().In(&v),
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
		}, v)
	})

	t.Run("scheme", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://joe:secret@%s", s.Listener.Addr().String()),
			Store().Request().URL().Scheme().In(&v),
		)

		require.Equal(t, "http", v)
	})

	t.Run("opaque", func(t *testing.T) {
		var v string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Opaque().In(&v),
		)

		require.Equal(t, "", v)
	})

	t.Run("user", func(t *testing.T) {
		t.Run("all", func(t *testing.T) {
			var v url.Userinfo
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().In(&v),
			)

			require.Equal(t, url.UserPassword("joe", "secret"), &v)
		})

		t.Run("username", func(t *testing.T) {
			var v string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().Username().In(&v),
			)

			require.Equal(t, "joe", v)
		})
		t.Run("password", func(t *testing.T) {
			var v string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().Password().In(&v),
			)

			require.Equal(t, "secret", v)
		})

		t.Run("string", func(t *testing.T) {
			var v string
			Test(t,
				Post("http://joe:secret@%s", s.Listener.Addr().String()),
				Store().Request().URL().User().String().In(&v),
			)

			require.Equal(t, "joe:secret", v)
		})

		t.Run("nil password", func(t *testing.T) {
			var v string
			Test(t,
				Post(s.URL),
				Store().Request().URL().User().Password().In(&v),
			)

			require.Equal(t, "", v)
		})
	})

	t.Run("Host", func(t *testing.T) {
		var v string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Host().In(&v),
		)

		require.Equal(t, s.Listener.Addr().String(), v)
	})

	t.Run("Hostname", func(t *testing.T) {
		var v string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Hostname().In(&v),
		)

		host, _, err := net.SplitHostPort(s.Listener.Addr().String())
		require.NoError(t, err)
		require.Equal(t, host, v)
	})

	t.Run("Port", func(t *testing.T) {
		var v string
		Test(t,
			Post(s.URL),
			Store().Request().URL().Port().In(&v),
		)

		_, port, err := net.SplitHostPort(s.Listener.Addr().String())
		require.NoError(t, err)
		require.Equal(t, port, v)
	})

	t.Run("Path", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().Path().In(&v),
		)

		require.Equal(t, "/foo/bar", v)
	})

	t.Run("EscapedPath", func(t *testing.T) {
		var path string
		Test(t,
			Post("http://%s/foo/bar", s.Listener.Addr().String()),
			Store().Request().URL().EscapedPath().In(&path),
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

	t.Run("Query", func(t *testing.T) {
		t.Run("all", func(t *testing.T) {
			var v url.Values
			Test(t,
				Post("http://%s?foo=bar", s.Listener.Addr().String()),
				Store().Request().URL().Query().In(&v),
			)

			require.Equal(t, v, url.Values{
				"foo": []string{"bar"},
			})
		})
		t.Run("specific", func(t *testing.T) {
			var v string
			Test(t,
				Post("http://%s?foo=bar", s.Listener.Addr().String()),
				Store().Request().URL().Query("foo").In(&v),
			)

			require.Equal(t, v, "bar")
		})
		t.Run("specific invalid", func(t *testing.T) {
			var v string
			Test(t,
				Post("http://%s?foo=bar", s.Listener.Addr().String()),
				Store().Request().URL().Query("foo2").In(&v),
			)

			require.Equal(t, v, "")
		})
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

	t.Run("IsAbs", func(t *testing.T) {
		var v bool
		Test(t,
			Post("http://%s/foo/bar#hash", s.Listener.Addr().String()),
			Store().Request().URL().IsAbs().In(&v),
		)

		require.True(t, v)
	})

	t.Run("RequestURI", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://%s/foo/bar#hash", s.Listener.Addr().String()),
			Store().Request().URL().RequestURI().In(&v),
		)

		require.Equal(t, "/foo/bar", v)
	})

	t.Run("String", func(t *testing.T) {
		var v string
		Test(t,
			Post("http://%s/foo/bar#hash", s.Listener.Addr().String()),
			Store().Request().URL().String().In(&v),
		)

		require.Equal(t, fmt.Sprintf("http://%s/foo/bar#hash", s.Listener.Addr().String()), v)
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
		Send().Body().String("Hello World"),
		Store().Request().ContentLength().In(&contentLength),
	)

	require.Equal(t, 11, contentLength)
}

func TestStoreRequest_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var transferEncoding []string
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Custom(BeforeExpectStep, func(hit Hit) error {
			hit.Request().TransferEncoding = []string{"gzip", "chunked"}
			return nil
		}),
		Store().Request().TransferEncoding().In(&transferEncoding),
	)

	require.Equal(t, []string{"gzip", "chunked"}, transferEncoding)
}

func TestStoreRequest_Host(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var host string

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Send().Headers("Host").Add("example.com"),
		Store().Request().Host().In(&host),
	)

	require.Equal(t, "example.com", host)
}
