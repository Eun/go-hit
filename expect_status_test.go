package hit_test

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectStatus_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status(http.StatusOK),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status(http.StatusNotFound),
		), PtrStr("Expected status code to be 404 but was 200 instead"),
	)
}

func TestExpectStatus_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().NotEqual(http.StatusNotFound),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().NotEqual(http.StatusOK),
		), PtrStr("Expected status code not to be 200"),
	)
}

func TestExpectStatusOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().OneOf(300, 400),
		),
		PtrStr(`[]int{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`} does not contain 200`),
	)
}

func TestExpectStatusNotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().NotOneOf(200, 300),
		),
		PtrStr(`[]int{`), PtrStr(`200,`), PtrStr(`300,`), PtrStr(`} should not contain 200`),
	)
}

func TestExpectStatus_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectStatus_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectStatus_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectStatus_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectStatus_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectStatus_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Status().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}

func TestExpectStatus_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Expect().Status(value).Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).Equal(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).NotEqual(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).OneOf()),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).NotOneOf()),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).GreaterThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).GreaterThan(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).LessThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).LessThan(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).GreaterOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).GreaterOrEqualThan(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).LessOrEqualThan()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).LessOrEqualThan(0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).Between()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).Between(0, 0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})

	t.Run("Expect().Status(value).NotBetween()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect().Status(0).NotBetween(0, 0)),
			PtrStr("only usable with Expect().Status() not with Expect().Status(value)"),
		)
	})
}

func TestExpectStatus_WithoutArgument(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Status(),
		),
		PtrStr("unable to run Expect().Status() without an argument or without a chain. Please use Expect().Status(something) or Expect().Status().Something"),
	)
}
