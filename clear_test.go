package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClear_Send(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Send(),
			Expect().Body("Hello World"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), nil, nil, nil, nil, nil,
	)
}

func TestClear_Expect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body("Hello World"),
			Clear().Expect(),
			Expect().Body("Hello Nature"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
	)
}

func TestClear_CombineSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			CombineSteps(
				Post(s.URL),
				Send().Body("Hello"),
				Expect().Body("World"),
			),
			Clear().Expect(),
			Expect().Body("Nature"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "Nature"`), nil, nil, nil, nil, nil,
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
			PtrStr(`unable to find a step with Send()`),
		)
	})

	t.Run("Expect", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Clear().Expect(),
			),
			PtrStr(`unable to find a step with Expect()`),
		)
	})
}
