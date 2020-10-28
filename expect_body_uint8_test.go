package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyUint8_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().Equal(20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().Equal(40),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyUint8_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().NotEqual(40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().NotEqual(20),
		),
		PtrStr("should not be 0x14"),
	)
}

func TestExpectBodyUint8OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().OneOf(20, 30),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().OneOf(30, 40),
		),
		PtrStr(`0x14 should be one of []interface {}{`), PtrStr(`0x1e,`), PtrStr(`0x28,`), PtrStr(`}`),
	)
}

func TestExpectBodyUint8NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().NotOneOf(30, 40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().NotOneOf(20, 30),
		),
		PtrStr(`0x14 should not contain 0x14`),
	)
}

func TestExpectBodyUint8_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().GreaterThan(19),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().GreaterThan(29),
		),
		PtrStr("expected 20 to be greater than 29"),
	)
}

func TestExpectBodyUint8_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().GreaterOrEqualThan(20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().GreaterOrEqualThan(29),
		),
		PtrStr("expected 20 to be greater or equal than 29"),
	)
}

func TestExpectBodyUint8_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().LessThan(21),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().LessThan(10),
		),
		PtrStr("expected 20 to be less than 10"),
	)
}

func TestExpectBodyUint8_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().LessOrEqualThan(20),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().LessOrEqualThan(10),
		),
		PtrStr("expected 20 to be less or equal than 10"),
	)
}

func TestExpectBodyUint8_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().Between(10, 20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().Between(30, 40),
		),
		PtrStr("expected 20 to be between 30 and 40"),
	)
}

func TestExpectBodyUint8_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint8(20),
		Expect().Body().Uint8().NotBetween(30, 40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().NotBetween(20, 30),
		),
		PtrStr("expected 20 not to be between 20 and 30"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint8(20),
			Expect().Body().Uint8().NotBetween(10, 20),
		),
		PtrStr("expected 20 not to be between 10 and 20"),
	)
}
