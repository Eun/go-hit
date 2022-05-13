package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestClear_Send(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Clear().Send(),
			Expect().Body().String().Equal("Hello World"),
		),
		PtrStr("not equal"), PtrStr(`expected: "Hello World"`), nil, nil, nil, nil, nil,
	)
}

func TestClear_Expect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Equal("Hello World"),
			Clear().Expect(),
			Expect().Body().String().Equal("Hello Nature"),
		),
		PtrStr("not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
	)
}

func TestClear_CombineSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			CombineSteps(
				Post(s.URL),
				Send().Body().String("Hello"),
				Expect().Body().String().Equal("World"),
			),
			Clear().Expect(),
			Expect().Body().String().Equal("Nature"),
		),
		PtrStr("not equal"), PtrStr(`expected: "Nature"`), nil, nil, nil, nil, nil,
	)
}

func TestClear_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Send", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Clear().Send(),
			),
			PtrStr(`unable to find a step with Send()`), PtrStr("got these steps:"), PtrStr("Post()"),
		)
	})

	t.Run("Expect", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Clear().Expect(),
			),
			PtrStr(`unable to find a step with Expect()`), PtrStr("got these steps:"), PtrStr("Post()"),
		)
	})
}
