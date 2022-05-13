package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`Hello World`),
			Expect().Body().Bytes().Equal([]byte("Hello World")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Equal(`Hello World`),
		)
	})

	t.Run("slice (JSON)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A", "B"]`),
			Expect().Body().JSON().Equal([]string{"A", "B"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("8"),
			Expect().Body().Int().Equal(8),
		)
	})
}

func TestExpectBody_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`Hello World`),
			Expect().Body().Bytes().NotEqual([]byte("Hello Universe")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().NotEqual(`Hello Universe`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("8"),
			Expect().Body().Int().NotEqual(6),
		)
	})

	t.Run("slice (JSON)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A", "B"]`),
			Expect().Body().JSON().NotEqual([]string{"A", "B", "C"}),
		)
	})
}

func TestExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Contains(`World`),
		)
	})

	t.Run("slice (JSON)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().Contains([]string{"Hello World"}),
			),
			PtrStr(`"Hello World" does not contain []string{`), PtrStr(`"Hello World",`), PtrStr("}"),
		)
	})
}

func TestExpectBody_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().NotContains(`Universe`),
		)
	})

	t.Run("slice (JSON)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().NotContains([]string{"Hello Universe"}),
		)
	})
}
