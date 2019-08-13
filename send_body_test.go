package hit_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSendBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})
}

func TestSendBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body([]byte("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("reader", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(bytes.NewBufferString("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(8),
			Expect().Body().Equal(`8`),
		)
	})
}

func TestSendBody_ModifyPreviousBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
			Send().Body(func(hit Hit) {
				hit.Request().Body().SetString(fmt.Sprintf("%s World", hit.Request().Body().String()))
			}),
			Expect().Body().Equal(`Hello World`),
		)
	})
}

func TestSendBody_EmptyBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			Send().Body(func(hit Hit) {
				require.Empty(t, hit.Request().Body().String())
			}),
		)
	})
}
