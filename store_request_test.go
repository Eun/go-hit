package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"

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
