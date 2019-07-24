package hit_test

import (
	"testing"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`Hello World`).
			Expect().Body([]byte("Hello World")).
			Do()
	})

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("Hello World").
			Expect().Body(`Hello World`).
			Do()
	})

	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`["A","B"]`).
			Expect().Body([]interface{}{"A", "B"}).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("8").
			Expect().Body(8).
			Do()
	})
}

func TestExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("Hello World").
			Expect().Body().Contains(`World`).
			Do()
	})

	t.Run("slice", func(t *testing.T) {
		// slice goes to json
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("invalid character 'H' looking for beginning of value")), s.URL).
				Send().Body("Hello World").
				Expect().Body().Contains([]string{"Hello World"}).
				Do()
		})
	})
}
