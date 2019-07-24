package hit_test

import (
	"testing"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestMethod(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		e := hit.Connect(t, "http://127.0.0.1")
		require.Equal(t, "CONNECT", e.Request().Method)
	})

	t.Run("Delete", func(t *testing.T) {
		e := hit.Delete(t, "http://127.0.0.1")
		require.Equal(t, "DELETE", e.Request().Method)
	})

	t.Run("Get", func(t *testing.T) {
		e := hit.Get(t, "http://127.0.0.1")
		require.Equal(t, "GET", e.Request().Method)
	})

	t.Run("Head", func(t *testing.T) {
		e := hit.Head(t, "http://127.0.0.1")
		require.Equal(t, "HEAD", e.Request().Method)
	})

	t.Run("Post", func(t *testing.T) {
		e := hit.Post(t, "http://127.0.0.1")
		require.Equal(t, "POST", e.Request().Method)
	})

	t.Run("Options", func(t *testing.T) {
		e := hit.Options(t, "http://127.0.0.1")
		require.Equal(t, "OPTIONS", e.Request().Method)
	})

	t.Run("Trace", func(t *testing.T) {
		e := hit.Trace(t, "http://127.0.0.1")
		require.Equal(t, "TRACE", e.Request().Method)
	})

	t.Run("Put", func(t *testing.T) {
		e := hit.Put(t, "http://127.0.0.1")
		require.Equal(t, "PUT", e.Request().Method)
	})
}
