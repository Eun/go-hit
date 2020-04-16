package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("inline", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) {
				calledFunc = true
				hit.Request().Body().SetString("Hello World")
			}),
			Expect().Body().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
	t.Run("MustDo", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) {
				calledFunc = true
				hit.MustDo(Send().Body("Hello World"))
			}),
			Expect().Body().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
}

func TestSendHeader(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Header("X-Header", "World"),
		Expect().Header("X-Header").Equal("World"),
	)
}

func TestSendTrailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Trailer("X-Trailer", "Hello"),
		Send().Body("Hello"),
		Expect().Trailer("X-Trailer").Equal("Hello"),
	)
}

func TestSendHeader_DoubleSet(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Header("X-Header", "World"),
		Send().Header("X-Header", "Universe"),
		Expect().Header("X-Header").Equal("Universe"),
	)
}
