package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyUint16_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyUint16_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().NotEqual(200),
		),
		PtrStr("should not be 0x00c8"),
	)
}

func TestExpectBodyUint16OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().OneOf(300, 400),
		),
		PtrStr(`0x00c8 should be one of []interface {}{`), PtrStr(`0x012c,`), PtrStr(`0x0190,`), PtrStr(`}`),
	)
}

func TestExpectBodyUint16NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().NotOneOf(200, 300),
		),
		PtrStr(`0x00c8 should not contain 0x00c8`),
	)
}

func TestExpectBodyUint16_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyUint16_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyUint16_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyUint16_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyUint16_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyUint16_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint16(200),
		Expect().Body().Uint16().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint16(200),
			Expect().Body().Uint16().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
