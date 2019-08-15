package hit_test

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectStatus_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status(http.StatusOK),
		)
	})

	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status(http.StatusNotFound),
			), PtrStr("Expected status code to be 404 but was 200 instead"),
		)
	})
}

func TestExpectStatusOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().OneOf(200, 300),
		)
	})

	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().OneOf(300, 400),
			),
			PtrStr(`[]int{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`} does not contain 200`),
		)
	})
}

func TestExpectStatus_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().GreaterThan(199),
		)
	})

	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().GreaterThan(299),
			),
			PtrStr("expected 200 to be greater than 299"),
		)
	})
}

func TestExpectStatus_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().GreaterOrEqualThan(200),
		)
	})
	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().GreaterOrEqualThan(299),
			),
			PtrStr("expected 200 to be greater or equal than 299"),
		)
	})
}

func TestExpectStatus_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().LessThan(201),
		)
	})

	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().LessThan(100),
			),
			PtrStr("expected 200 to be less than 100"),
		)
	})
}

func TestExpectStatus_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().LessOrEqualThan(200),
		)
	})

	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().LessOrEqualThan(100),
			),
			PtrStr("expected 200 to be less or equal than 100"),
		)
	})
}

func TestExpectStatus_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("ok", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Status().Between(100, 200),
		)
	})
	t.Run("failing", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Status().Between(300, 400),
			),
			PtrStr("expected 200 to be between 300 and 400"),
		)
	})
}
