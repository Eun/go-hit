package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClearSend(t *testing.T) {
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

func TestClearExpect(t *testing.T) {
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
