package hit_test

import (
	"bytes"
	"testing"

	"fmt"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSendBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body().JSON([]string{"A", "B"}).
			Expect().Body().Equal(`["A","B"]`).
			Do()
	})
}

func TestSendBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body([]byte("Hello World")).
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("Hello World").
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("reader", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(bytes.NewBufferString("Hello World")).
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body([]string{"A", "B"}).
			Expect().Body().Equal(`["A","B"]`).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(8).
			Expect().Body().Equal(`8`).
			Do()
	})
}

func TestSendBody_ModifyPreviousBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("Hello").
			Send().Body(func(hit hit.Hit) {
			hit.Request().Body().SetString(fmt.Sprintf("%s World", hit.Request().Body().String()))
		}).
			Expect().Body().Equal(`Hello World`).
			Do()
	})
}

func TestSendBody_EmptyBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Get(t, s.URL).
			Send().Body(func(hit hit.Hit) {
			require.Empty(t, hit.Request().Body().String())
		}).
			Do()
	})
}

func TestSendBody_AfterDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL).
				Do().
				Send().Body("Hello")
		})
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL).
				Do().
				Send().Body().JSON("Hello")
		})
	})
}
