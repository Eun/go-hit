package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyUint_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyUint_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().NotEqual(200),
		),
		PtrStr("should not be 0xc8"),
	)
}

func TestExpectBodyUintOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().OneOf(300, 400),
		),
		PtrStr(`0xc8 should be one of []interface {}{`), PtrStr(`0x12c,`), PtrStr(`0x190,`), PtrStr(`}`),
	)
}

func TestExpectBodyUintNotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().NotOneOf(200, 300),
		),
		PtrStr(`0xc8 should not contain 0xc8`),
	)
}

func TestExpectBodyUint_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyUint_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyUint_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyUint_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyUint_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyUint_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint(200),
		Expect().Body().Uint().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint(200),
			Expect().Body().Uint().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
