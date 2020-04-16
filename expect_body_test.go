package hit_test

import (
	"testing"

	"errors"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBody_Interface(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Response)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect().Body(func(hit Hit) {
						if hit.Response().Body().String() != "Hello Universe" {
							panic("Not equal")
						}
					}),
				),
				PtrStr("Not equal"),
			)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect().Body(func(hit Hit) {
						hit.MustDo(Expect().Body("Hello Universe"))
					})),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
			)
		})
		t.Run("with correct parameter (using Hit) and error", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect().Body(func(hit Hit) error {
						return hit.Do(Expect().Body("Hello Universe"))
					})),
				PtrStr("Not equal"), PtrStr(`expected: "Hello Universe"`), PtrStr(`actual: "Hello World"`), nil, nil, nil, nil,
			)
		})
		t.Run("with correct parameter (using Hit) and error (return an error)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect().Body(func(hit Hit) error {
						return errors.New("whoops")
					})),
				PtrStr("whoops"),
			)
		})
		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body(func() {
					calledFunc = true
				}),
			)
			require.True(t, calledFunc)
		})
	})
}

func TestExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Expect().Body([]byte("Hello World")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A", "B"]`),
				Expect().Body([]string{"A", "B"}),
			),
			PtrStr("unsure howto compare text/plain; charset=utf-8 with []string"),
		)
	})

	t.Run("slice (Content-Type=json)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A", "B"]`),
			Expect().Custom(func(hit Hit) {
				hit.Response().Header.Set("Content-Type", "application/json")
			}),
			Expect().Body([]string{"A", "B"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("8"),
			Expect().Body(8),
		)
	})
}

func TestExpectBody_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Expect().Body().NotEqual([]byte("Hello Universe")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().NotEqual(`Hello Universe`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("8"),
			Expect().Body().NotEqual(6),
		)
	})

	t.Run("slice", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A", "B"]`),
				Expect().Body().NotEqual([]string{"A", "B", "C"}),
			),
			PtrStr("unsure howto compare text/plain; charset=utf-8 with []string"),
		)
	})

	t.Run("slice (Content-Type=json)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A", "B"]`),
			Expect().Custom(func(hit Hit) {
				hit.Response().Header.Set("Content-Type", "application/json")
			}),
			Expect().Body().NotEqual([]string{"A", "B", "C"}),
		)
	})
}

func TestExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains(`World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains([]string{"Hello World"}),
			),
			PtrStr("unsure howto compare text/plain; charset=utf-8 with []string"),
		)
	})

	t.Run("slice (Content-Type=json)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Custom(func(hit Hit) {
					hit.Response().Header.Set("Content-Type", "application/json")
				}),
				Expect().Body().Contains([]string{"Hello World"}),
			),
			PtrStr("invalid character 'H' looking for beginning of value"),
		)
	})
}

func TestExpectBody_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().NotContains(`Universe`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().NotContains([]string{"Hello Universe"}),
			),
			PtrStr("unsure howto compare text/plain; charset=utf-8 with []string"),
		)
	})

	t.Run("slice (Content-Type=json)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Custom(func(hit Hit) {
					hit.Response().Header.Set("Content-Type", "application/json")
				}),
				Expect().Body().NotContains([]string{"Hello Universe"}),
			),
			PtrStr("invalid character 'H' looking for beginning of value"),
		)
	})
}

func TestExpectBody_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Expect().Body(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").Interface(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").JSON(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).JSON().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").JSON().Equal("", "")),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").Equal(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").NotEqual(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).Contains()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").Contains(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})

	t.Run("Expect().Body(value).NotContains()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Body("Data").NotContains(nil)),
			PtrStr("only usable with Expect().Body() not with Expect().Body(value)"),
		)
	})
}

func TestExpectBody_WithoutArgument(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Body(),
		),
		PtrStr("unable to run Expect().Body() without an argument or without a chain. Please use Expect().Body(something) or Expect().Body().Something"),
	)
}
