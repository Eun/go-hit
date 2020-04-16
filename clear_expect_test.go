package hit_test

import (
	"testing"

	"net/http"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestClearExpect_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body("Hello Nature"),
				Clear().Expect().Body(),
				Expect().Body("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body("Hello Nature"),
				Expect().Body("Hello World"),
				Clear().Expect().Body("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpect_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Header("X-Header", "Hello Earth"),
				Expect().Header("X-Header").Equal("Hello Nature"),
				Clear().Expect().Header(),
				Expect().Header("X-Header").Equal("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Header("X-Header", "Hello Earth"),
				Expect().Header("X-Header").Equal("Hello Nature"),
				Expect().Header("X-Header").Equal("Hello World"),
				Clear().Expect().Header("X-Header").Equal("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpect_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Send().Trailer("X-Trailer", "Hello Earth"),
				Expect().Trailer("X-Trailer").Equal("Hello Nature"),
				Clear().Expect().Trailer(),
				Expect().Trailer("X-Trailer").Equal("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Send().Trailer("X-Trailer", "Hello Earth"),
				Expect().Trailer("X-Trailer").Equal("Hello Nature"),
				Expect().Trailer("X-Trailer").Equal("Hello World"),
				Clear().Expect().Trailer("X-Trailer").Equal("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpect_Status(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status(http.StatusOK),
				Clear().Expect().Status(),
				Expect().Status(http.StatusNotFound),
			),
			PtrStr("Expected status code to be 404 but was 200 instead"),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status(http.StatusOK),
				Expect().Status(http.StatusNotFound),
				Clear().Expect().Status(http.StatusOK),
			),
			PtrStr("Expected status code to be 404 but was 200 instead"),
		)
	})
}

func TestClearExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect().Body("Hello Universe"))
				}),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect().Body("Hello Earth"))
				}),
				Clear().Expect().Custom(),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect().Body("Hello Nature"))
				}),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ranCustomFunc := false
		fn := func(hit Hit) {
			require.Equal(t, "Hello Universe", hit.Response().Body().String())
		}
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Custom(fn),
			Expect().Custom(fn),
			Expect().Custom(func(hit Hit) {
				ranCustomFunc = true
			}),
			Clear().Expect().Custom(fn),
		)
		require.True(t, ranCustomFunc)
	})
}

func TestClearExpect_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Body("Hello Universe"),
		),
		PtrStr(`unable to find a step with Expect().Body("Hello Universe")`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect(),
		),
		PtrStr(`unable to find a step with Expect()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
