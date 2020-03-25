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

func TestClearExpect_Interface(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Interface("Hello Earth"),
				Expect().Interface("Hello Nature"),
				Clear().Expect().Interface(),
				Expect().Interface("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Interface("Hello Earth"),
				Expect().Interface("Hello Nature"),
				Expect().Interface("Hello World"),
				Clear().Expect().Interface("Hello Nature"),
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
					hit.MustDo(Expect("Hello Universe"))
				}),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect("Hello Earth"))
				}),
				Clear().Expect().Custom(),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect("Hello Nature"))
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

func TestClearExpect_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	t.Run("Clear().Expect(value).Body()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Body()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Body().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Body().Equal()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Interface()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Custom()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Custom()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().Contains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().Contains()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().NotContains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().NotContains()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().OneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().NotOneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().Empty()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().Empty()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().Len()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().Len()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().Equal()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Header().NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Header().NotEqual()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().Contains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().Contains()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().NotContains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().NotContains()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().OneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().NotOneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().Empty()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().Empty()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().Len()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().Len()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().Equal()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Trailer().NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Trailer().NotEqual()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().Equal()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().NotEqual()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().OneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().NotOneOf()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().GreaterThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().GreaterThan()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().LessThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().LessThan()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().GreaterOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().GreaterOrEqualThan()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().LessOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().LessOrEqualThan()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().Between()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().Between()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})

	t.Run("Clear().Expect(value).Status().NotBetween()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect("Data").Status().NotBetween()),
			PtrStr("only usable with Clear().Expect() not with Clear().Expect(value)"),
		)
	})
}

func TestClearExpect_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect("Hello Universe"),
		),
		PtrStr(`unable to find a step with Expect("Hello Universe")`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
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
