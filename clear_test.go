package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClear_Send(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send("Hello World"),
				Clear().Send(),
				Send("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send("Hello World"),
				Clear().Send("Hello World"),
				Send("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), nil, nil, nil, nil, nil,
		)
	})
}

func TestClear_Expect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect("Hello Universe"),
				Expect("Hello Earth"),
				Clear().Expect(),
				Expect("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect("Hello Universe"),
				Expect("Hello Earth"),
				Clear().Expect("Hello Universe"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Earth"`), nil, nil, nil, nil, nil,
		)
	})
}

func TestClear_CombineSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			CombineSteps(
				Post(s.URL),
				Send("Hello"),
				Expect("World"),
			),
			Clear().Expect(),
			Expect("Nature"),
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
