package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyUint64_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyUint64_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().NotEqual(200),
		),
		PtrStr("should not be 0x00000000000000c8"),
	)
}

func TestExpectBodyUint64OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().OneOf(300, 400),
		),
		PtrStr(`0x00000000000000c8 should be one of []interface {}{`), PtrStr(`0x000000000000012c,`), PtrStr(`0x0000000000000190,`), PtrStr(`}`),
	)
}

func TestExpectBodyUint64NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().NotOneOf(200, 300),
		),
		PtrStr(`0x00000000000000c8 should not contain 0x00000000000000c8`),
	)
}

func TestExpectBodyUint64_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyUint64_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyUint64_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyUint64_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyUint64_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyUint64_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Uint64(200),
		Expect().Body().Uint64().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Uint64(200),
			Expect().Body().Uint64().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
