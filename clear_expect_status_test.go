package hit_test

import (
	"testing"

	"net/http"

	. "github.com/Eun/go-hit"
)

func TestClearExpectStatus_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().Equal(http.StatusOK),
				Clear().Expect().Status().Equal(),
				Expect().Status().Equal(http.StatusNotFound),
			),
			PtrStr("Expected status code to be 404 but was 200 instead"),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().Equal(http.StatusOK),
				Expect().Status().Equal(http.StatusNotFound),
				Clear().Expect().Status().Equal(http.StatusOK),
			),
			PtrStr("Expected status code to be 404 but was 200 instead"),
		)
	})
}

func TestClearExpectStatus_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().NotEqual(http.StatusOK),
			Clear().Expect().Status().NotEqual(),
			Expect().Status().NotEqual(http.StatusNotFound),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().NotEqual(http.StatusOK),
			Expect().Status().NotEqual(http.StatusNotFound),
			Clear().Expect().Status().NotEqual(http.StatusOK),
		)
	})
}

func TestClearExpectStatus_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().OneOf(http.StatusNotFound),
			Clear().Expect().Status().OneOf(),
			Expect().Status().OneOf(http.StatusOK),
		)
	})

	t.Run("specific only first parameter ", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().OneOf(http.StatusNotFound, http.StatusNoContent),
			Expect().Status().OneOf(http.StatusOK),
			Clear().Expect().Status().OneOf(http.StatusNotFound),
		)
	})
	t.Run("specific (all)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().OneOf(http.StatusNotFound, http.StatusNoContent),
			Expect().Status().OneOf(http.StatusOK),
			Clear().Expect().Status().OneOf(http.StatusNotFound, http.StatusNoContent),
		)
	})
}

func TestClearExpectStatus_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().NotOneOf(http.StatusOK),
			Clear().Expect().Status().NotOneOf(),
			Expect().Status().NotOneOf(http.StatusNotFound),
		)
	})

	t.Run("specific only first parameter ", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().NotOneOf(http.StatusOK, http.StatusNoContent),
			Expect().Status().NotOneOf(http.StatusNotFound),
			Clear().Expect().Status().NotOneOf(http.StatusOK),
		)
	})
	t.Run("specific (all)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().NotOneOf(http.StatusOK, http.StatusNoContent),
			Expect().Status().NotOneOf(http.StatusNotFound),
			Clear().Expect().Status().NotOneOf(http.StatusOK, http.StatusNoContent),
		)
	})
}

func TestClearExpectStatus_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().GreaterThan(http.StatusOK),
			Clear().Expect().Status().GreaterThan(),
			Expect().Status().GreaterThan(0),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().GreaterThan(http.StatusOK),
			Expect().Status().GreaterThan(0),
			Clear().Expect().Status().GreaterThan(http.StatusOK),
		)
	})
}

func TestClearExpectStatus_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().LessThan(http.StatusOK),
			Clear().Expect().Status().LessThan(),
			Expect().Status().LessThan(500),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().LessThan(http.StatusOK),
			Expect().Status().LessThan(500),
			Clear().Expect().Status().LessThan(http.StatusOK),
		)
	})
}

func TestClearExpectStatus_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().GreaterOrEqualThan(500),
			Clear().Expect().Status().GreaterOrEqualThan(),
			Expect().Status().GreaterOrEqualThan(0),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().GreaterOrEqualThan(500),
			Expect().Status().GreaterOrEqualThan(0),
			Clear().Expect().Status().GreaterOrEqualThan(500),
		)
	})
}

func TestClearExpectStatus_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().LessOrEqualThan(1),
			Clear().Expect().Status().LessOrEqualThan(),
			Expect().Status().LessOrEqualThan(500),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello Earth"),
			Expect().Status().LessOrEqualThan(1),
			Expect().Status().LessOrEqualThan(500),
			Clear().Expect().Status().LessOrEqualThan(1),
		)
	})
}

func TestClearExpectStatus_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().Between(300, 400),
				Clear().Expect().Status().Between(),
				Expect().Status().Between(400, 500),
			),
			PtrStr("expected 200 to be between 400 and 500"),
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().Between(300, 400),
				Expect().Status().Between(400, 500),
				Clear().Expect().Status().Between(300),
			),
			PtrStr("expected 200 to be between 400 and 500"),
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().Between(300, 400),
				Expect().Status().Between(400, 500),
				Clear().Expect().Status().Between(300, 400),
			),
			PtrStr("expected 200 to be between 400 and 500"),
		)
	})
}

func TestClearExpectStatus_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().NotBetween(200, 201),
				Clear().Expect().Status().NotBetween(),
				Expect().Status().NotBetween(200, 300),
			),
			PtrStr("expected 200 not to be between 200 and 300"),
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().NotBetween(200, 201),
				Expect().Status().NotBetween(199, 300),
				Clear().Expect().Status().NotBetween(200),
			),
			PtrStr("expected 200 not to be between 199 and 300"),
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Status().NotBetween(200, 201),
				Expect().Status().NotBetween(199, 300),
				Clear().Expect().Status().NotBetween(200, 201),
			),
			PtrStr("expected 200 not to be between 199 and 300"),
		)
	})
}

func TestClearExpectStatus_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Clear().Expect().Status(value).Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).Equal()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).NotEqual()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).OneOf()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).NotOneOf()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).GreaterThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).GreaterThan()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).LessThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).LessThan()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).GreaterOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).GreaterOrEqualThan()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).LessOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).LessOrEqualThan()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).Between()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).Between()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})

	t.Run("Clear().Expect().Status(value).NotBetween()", func(t *testing.T) {
		ExpectError(t,
			Do(Clear().Expect().Status(0).NotBetween()),
			PtrStr("only usable with Clear().Expect().Status() not with Clear().Expect().Status(value)"),
		)
	})
}

func TestClearExpectStatus_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Status(200),
		),
		PtrStr(`unable to find a step with Expect().Status(200)`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Status(),
		),
		PtrStr(`unable to find a step with Expect().Status()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
