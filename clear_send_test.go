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
				Expect().Body("Hello World"),
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
				Expect().Body("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
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
					hit.MustDo(Send().Body("Hello World"))
				}),
				Clear().Send().Custom(),
				Send().Custom(func(hit Hit) {
					hit.MustDo(Send().Body("Hello Earth"))
				}),
				Expect().Body("Hello World"),
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
			Expect().Body("Hello World"),
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

func TestClearSend_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer", "Hello"),
				Clear().Send().Trailer(),
				Send().Trailer("X-Trailer", "World"),
				Expect().Trailer("X-Trailer").Equal("Hello"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello"`), PtrStr(`actual: "World"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello"),
				Send().Trailer("X-Trailer", "Hello"),
				Send().Trailer("X-Trailer", "World"),
				Clear().Send().Trailer("X-Trailer", "Hello"),
				Expect().Trailer("X-Trailer").Equal("Hello"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello"`), PtrStr(`actual: "World"`), nil, nil, nil, nil,
		)
	})
}

func TestClearSend_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Clear().Send().Body(),
		),
		PtrStr(`unable to find a step with Send().Body()`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Body("Hello World"),
			Clear().Send(),
		),
		PtrStr(`unable to find a step with Send()`), PtrStr(`got these steps:`), PtrStr(`Expect().Body("Hello World")`),
	)
}
