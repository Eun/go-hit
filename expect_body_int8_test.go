package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyInt8_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().Equal(20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().Equal(40),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyInt8_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().NotEqual(40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().NotEqual(20),
		),
		PtrStr("should not be 20"),
	)
}

func TestExpectBodyInt8OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().OneOf(20, 30),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().OneOf(30, 40),
		),
		PtrStr(`20 should be one of []interface {}{`), PtrStr(`30,`), PtrStr(`40,`), PtrStr(`}`),
	)
}

func TestExpectBodyInt8NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().NotOneOf(30, 40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().NotOneOf(20, 30),
		),
		PtrStr(`20 should not contain 20`),
	)
}

func TestExpectBodyInt8_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().GreaterThan(19),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().GreaterThan(21),
		),
		PtrStr("expected 20 to be greater than 21"),
	)
}

func TestExpectBodyInt8_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().GreaterOrEqualThan(20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().GreaterOrEqualThan(21),
		),
		PtrStr("expected 20 to be greater or equal than 21"),
	)
}

func TestExpectBodyInt8_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().LessThan(21),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().LessThan(10),
		),
		PtrStr("expected 20 to be less than 10"),
	)
}

func TestExpectBodyInt8_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().LessOrEqualThan(20),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().LessOrEqualThan(10),
		),
		PtrStr("expected 20 to be less or equal than 10"),
	)
}

func TestExpectBodyInt8_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().Between(10, 20),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().Between(30, 40),
		),
		PtrStr("expected 20 to be between 30 and 40"),
	)
}

func TestExpectBodyInt8_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int8(20),
		Expect().Body().Int8().NotBetween(30, 40),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().NotBetween(20, 30),
		),
		PtrStr("expected 20 not to be between 20 and 30"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int8(20),
			Expect().Body().Int8().NotBetween(10, 20),
		),
		PtrStr("expected 20 not to be between 10 and 20"),
	)
}
