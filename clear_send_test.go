package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestClearSend_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Clear().Send().Body(),
				Send().Body("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Send().Body("Hello Earth"),
				Clear().Send().Body("Hello World"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearSend_Interface(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Interface("Hello World"),
				Clear().Send().Interface(),
				Send().Interface("Hello Earth"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Interface("Hello World"),
				Send().Interface("Hello Earth"),
				Clear().Send().Interface("Hello World"),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearSend_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().JSON("Hello World"),
				Clear().Send().JSON(),
				Send().JSON("Hello Earth"),
				Expect(`"Hello World"`),
			),
			PtrStr("Not equal"), PtrStr(`expected: "\"Hello World\""`), PtrStr(`actual: "\"Hello Earth\""`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().JSON("Hello World"),
				Send().JSON("Hello Earth"),
				Clear().Send().JSON("Hello World"),
				Expect(`"Hello World"`),
			),
			PtrStr("Not equal"), PtrStr(`expected: "\"Hello World\""`), PtrStr(`actual: "\"Hello Earth\""`), nil, nil, nil, nil,
		)
	})
}

func TestClearSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Custom(func(hit Hit) {
					hit.MustDo(Send("Hello World"))
				}),
				Clear().Send().Custom(),
				Send().Custom(func(hit Hit) {
					hit.MustDo(Send("Hello Earth"))
				}),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ranCustomFunc := false
		fn := func(hit Hit) {
			hit.Request().Body().SetString("Hello Earth")
		}
		Test(t,
			Post(s.URL),
			Send().Custom(fn),
			Clear().Send().Custom(fn),
			Send().Custom(func(hit Hit) {
				ranCustomFunc = true
				hit.Request().Body().SetString("Hello World")
			}),
			Expect("Hello World"),
		)
		require.True(t, ranCustomFunc)
	})
}

func TestClearSend_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Header("X-Header", "Hello"),
				Clear().Send().Header(),
				Send().Header("X-Header", "World"),
				Expect().Header("X-Header").Equal("Hello"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello"`), PtrStr(`actual: "World"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Header("X-Header", "Hello"),
				Send().Header("X-Header", "World"),
				Clear().Send().Header("X-Header", "Hello"),
				Expect().Header("X-Header").Equal("Hello"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello"`), PtrStr(`actual: "World"`), nil, nil, nil, nil,
		)
	})
}

func TestClearSend_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Clear().Send(value).Body()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").Body()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})

	t.Run("Clear().Send(value).Body().JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").Body().JSON()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})

	t.Run("Clear().Send(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").Interface()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})

	t.Run("Clear().Send(value).JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").JSON()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})

	t.Run("Clear().Send(value).Header()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").Header()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})

	t.Run("Clear().Send(value).Custom()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Send("Data").Custom()),
			PtrStr("only usable with Clear().Send() not with Clear().Send(value)"),
		)
	})
}
