package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClearExpectBody_Interface(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Interface("Hello Earth"),
				Clear().Expect().Body().Interface(),
				Expect().Body().Interface("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Interface("Hello Earth"),
				Expect().Body().Interface("Hello Nature"),
				Clear().Expect().Body().Interface("Hello Earth"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON("Hello Earth"),
				Expect().Body().JSON("Hello World"),
				Clear().Expect().Body().JSON(),
				Expect().Body().JSON("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON("Hello Earth"),
				Expect().Body().JSON("Hello World"),
				Expect().Body().JSON("Hello Nature"),
				Clear().Expect().Body().JSON("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Equal("Hello Earth"),
				Clear().Expect().Body().Equal(),
				Expect().Body().Equal("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Equal("Hello World"),
				Expect().Body().Equal("Hello Nature"),
				Clear().Expect().Body().Equal("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectBody_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Body().NotEqual("Hello Earth"),
			Clear().Expect().Body().NotEqual(),
			Expect().Body().NotEqual("Hello World"),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Body().NotEqual("Hello Nature"),
			Expect().Body().NotEqual("Hello Earth"),
			Clear().Expect().Body().NotEqual("Hello Earth"),
		)
	})
}

func TestClearExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Contains("Hello Nature"),
				Clear().Expect().Body().Contains(),
				Expect().Body().Contains("Hello World"),
			),
			PtrStr(`"Hello Earth" does not contain "Hello World"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().Contains("Hello Nature"),
				Expect().Body().Contains("Hello World"),
				Clear().Expect().Body().Contains("Hello Nature"),
			),
			PtrStr(`"Hello Earth" does not contain "Hello World"`),
		)
	})
}

func TestClearExpectBody_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().NotContains("Hello Nature"),
				Clear().Expect().Body().NotContains(),
				Expect().Body().NotContains("Hello Earth"),
			),
			PtrStr(`"Hello Earth" does contain "Hello Earth"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello Earth"),
				Expect().Body().NotContains("Hello Nature"),
				Expect().Body().NotContains("Hello Earth"),
				Clear().Expect().Body().NotContains("Hello Nature"),
			),
			PtrStr(`"Hello Earth" does contain "Hello Earth"`),
		)
	})
}

func TestClearExpectBody_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Clear().Expect().Body(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").Interface()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").JSON()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).JSON().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").JSON().Equal()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").Equal()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").NotEqual()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).Contains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").Contains()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})

	t.Run("Clear().Expect().Body(value).NotContains()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Body("Data").NotContains()),
			PtrStr("only usable with Clear().Expect().Body() not with Clear().Expect().Body(value)"),
		)
	})
}

func TestClearExpectBody_NotExistentStep(t *testing.T) {
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
			Clear().Expect().Body(),
		),
		PtrStr(`unable to find a step with Expect().Body()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
