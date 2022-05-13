package hit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/otto-eng/go-hit"
)

func TestSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("inline", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) error {
				calledFunc = true
				hit.Request().Body().SetString("Hello World")
				return nil
			}),
			Expect().Body().String().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
	t.Run("MustDo", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) error {
				calledFunc = true
				hit.MustDo(Send().Body().String("Hello World"))
				return nil
			}),
			Expect().Body().String().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
}

func TestSendHeaders(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Headers("X-Header").Add("World"),
		Expect().Headers("X-Header").Equal("World"),
	)
}

func TestSendTrailers(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Trailers("X-Trailer").Add("Hello"),
		Send().Body().String("Hello"),
		Expect().Trailers("X-Trailer").Equal("Hello"),
	)
}
