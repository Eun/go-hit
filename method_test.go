package hit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/otto-eng/go-hit"
)

func TestMethod(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	t.Run("Connect", func(t *testing.T) {
		Test(t,
			Connect(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "CONNECT", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Delete", func(t *testing.T) {
		Test(t,
			Delete(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "DELETE", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Get", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "GET", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Head", func(t *testing.T) {
		Test(t,
			Head(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "HEAD", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Post", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "POST", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Options", func(t *testing.T) {
		Test(t,
			Options(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "OPTIONS", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Trace", func(t *testing.T) {
		Test(t,
			Trace(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "TRACE", hit.Request().Method)
				return nil
			}),
		)
	})

	t.Run("Put", func(t *testing.T) {
		Test(t,
			Put(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				require.Equal(t, "PUT", hit.Request().Method)
				return nil
			}),
		)
	})
}
