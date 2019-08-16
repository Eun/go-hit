package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Expect().Body([]byte("Hello World")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body([]interface{}{"A", "B"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("8"),
			Expect().Body(8),
		)
	})
}

func TestExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains(`World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		// slice goes to json
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains([]string{"Hello World"}),
			),
			PtrStr("invalid character 'H' looking for beginning of value"),
		)
	})
}
