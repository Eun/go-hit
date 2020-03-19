package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClearSendBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON("Hello World"),
				Clear().Send().Body().JSON(),
				Send().Body().JSON("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "\"Hello Earth\""`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON("Hello World"),
				Send().Body().JSON("Hello Earth"),
				Clear().Send().Body().JSON("Hello World"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "\"Hello Earth\""`), nil, nil, nil, nil,
		)
	})
}

func TestClearSendBody_Interface(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().Interface("Hello World"),
				Clear().Send().Body().Interface(),
				Send().Body().Interface("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().Interface("Hello World"),
				Send().Body().Interface("Hello Earth"),
				Clear().Send().Body().Interface("Hello World"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearSendBody_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Clear().Send().Body(value).JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send().Body("Data").JSON()),
			PtrStr("only usable with Clear().Send().Body() not with Clear().Send().Body(value)"),
		)
	})
	t.Run("Clear().Send().Body(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send().Body("Data").Interface()),
			PtrStr("only usable with Clear().Send().Body() not with Clear().Send().Body(value)"),
		)
	})
}

func TestClearSendBody_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Send().Body().JSON("Hello World"),
		),
		PtrStr(`unable to find a step with Send().Body().JSON("Hello World")`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Send().Body().JSON(),
		),
		PtrStr(`unable to find a step with Send().Body().JSON()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
