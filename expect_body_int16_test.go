package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyInt16_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyInt16_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().NotEqual(200),
		),
		PtrStr("should not be 200"),
	)
}

func TestExpectBodyInt16OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().OneOf(300, 400),
		),
		PtrStr(`200 should be one of []interface {}{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`}`),
	)
}

func TestExpectBodyInt16NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().NotOneOf(200, 300),
		),
		PtrStr(`200 should not contain 200`),
	)
}

func TestExpectBodyInt16_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyInt16_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyInt16_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyInt16_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyInt16_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyInt16_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int16(200),
		Expect().Body().Int16().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int16(200),
			Expect().Body().Int16().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
