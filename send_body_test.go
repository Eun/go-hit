package hit_test

import (
	"bytes"
	"fmt"
	"testing"

	"errors"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSendBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body([]byte("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("reader", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(bytes.NewBufferString("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body([]string{"A", "B"}),
				Expect().Body().Equal(`["A","B"]`),
			),
			PtrStr(`unable to set http body to []string{"A", "B"}, either specify a Content-Type or use a stringable value`),
		)
	})

	t.Run("slice (Content-Type=json)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Header("Content-Type", "application/json"),
			Send().Body([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})

	t.Run("slice (Content-Type=xml)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Header("Content-Type", "application/xml"),
			Send().Body([]string{"A", "B"}),
			Expect().Body().Equal(`<string>A</string><string>B</string>`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(8),
			Expect().Body().Equal(`8`),
		)
	})

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Request)", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body(func(e Hit) {
					calledFunc = true
					e.Request().Body().SetString("Hello World")
				}),
				Expect().Body().Equal(`Hello World`),
			)
			require.True(t, calledFunc)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body(func(e Hit) {
					calledFunc = true
					e.MustDo(Send().Body(`Hello World`))
				}),
				Expect().Body().Equal(`Hello World`),
			)
			require.True(t, calledFunc)
		})
		t.Run("with correct parameter (using Hit) and error", func(t *testing.T) {
			calledFunc := false
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(func(e Hit) error {
						calledFunc = true
						return e.Do(Send().Body(`Hello Earth`))
					}),
					Expect().Body().Equal(`Hello World`),
				),
				PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
			)
			require.True(t, calledFunc)
		})
		t.Run("with correct parameter (using Hit) and error (return an error)", func(t *testing.T) {
			calledFunc := false
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(func(e Hit) error {
						calledFunc = true
						return errors.New("whoops")
					}),
				),
				PtrStr("whoops"),
			)
			require.True(t, calledFunc)
		})
		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body(func() {
					calledFunc = true
				}),
				Expect().Body().Equal(``),
			)
			require.True(t, calledFunc)
		})
	})
}

func TestSendBody_Interface(t *testing.T) {
	// Interface is implicit tested trough TestSendBody
}
func TestSendBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON("Hello World"),
			Expect().Body().Equal(`"Hello World"`),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}),
			Expect().Body().Equal(`{"A":"1","B":"2"}`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(8),
			Expect().Body().Equal(`8`),
		)
	})

	t.Run("struct", func(t *testing.T) {
		var user = struct {
			Name string
			ID   int
		}{
			"Joe",
			10,
		}

		Test(t,
			Post(s.URL),
			Send().Body().JSON(user),
			Expect().Body().Equal(`{"Name":"Joe","ID":10}`),
		)
	})
}

func TestSendBody_XML(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().XML([]string{"A", "B"}),
		Expect().Body().Equal(`<string>A</string><string>B</string>`),
	)
}

func TestSendBody_ModifyPreviousBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello"),
		Send().Body(func(hit Hit) {
			hit.Request().Body().SetString(fmt.Sprintf("%s World", hit.Request().Body().String()))
		}),
		Expect().Body().Equal(`Hello World`),
	)
}

func TestSendBody_EmptyBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Get(s.URL),
		Send().Body(func(hit Hit) {
			require.Empty(t, hit.Request().Body().String())
		}),
	)
}

func TestSendBody_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Send().Body(value).JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Send().Body("Data").JSON(nil)),
			PtrStr("only usable with Send().Body() not with Send().Body(value)"),
		)
	})

	t.Run("Send().Body(value).XML()", func(t *testing.T) {
		ExpectError(t,
			Do(Send().Body("Data").XML(nil)),
			PtrStr("only usable with Send().Body() not with Send().Body(value)"),
		)
	})

	t.Run("Send().Body(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Send().Body("Data").Interface(nil)),
			PtrStr("only usable with Send().Body() not with Send().Body(value)"),
		)
	})
}

func TestSendBody_WithoutArgument(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body(),
		),
		PtrStr("unable to run Send().Body() without an argument or without a chain. Please use Send().Body(something) or Send().Body().Something"),
	)
}
