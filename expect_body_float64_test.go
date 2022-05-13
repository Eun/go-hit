package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyFloat64_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyFloat64_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().NotEqual(200),
		),
		PtrStr("should not be 200.000000"),
	)
}

func TestExpectBodyFloat64OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().OneOf(300, 400),
		),
		PtrStr(`200.000000 should be one of []interface {}{`), PtrStr(`300.000000,`), PtrStr(`400.000000,`), PtrStr(`}`),
	)
}

func TestExpectBodyFloat64NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().NotOneOf(200, 300),
		),
		PtrStr(`200.000000 should not contain 200.000000`),
	)
}

func TestExpectBodyFloat64_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().GreaterThan(299),
		),
		PtrStr("expected 200.000000 to be greater than 299.000000"),
	)
}

func TestExpectBodyFloat64_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200.000000 to be greater or equal than 299.000000"),
	)
}

func TestExpectBodyFloat64_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().LessThan(100),
		),
		PtrStr("expected 200.000000 to be less than 100.000000"),
	)
}

func TestExpectBodyFloat64_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().LessOrEqualThan(100),
		),
		PtrStr("expected 200.000000 to be less or equal than 100.000000"),
	)
}

func TestExpectBodyFloat64_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().Between(300, 400),
		),
		PtrStr("expected 200.000000 to be between 300.000000 and 400.000000"),
	)
}

func TestExpectBodyFloat64_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float64(200),
		Expect().Body().Float64().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().NotBetween(200, 300),
		),
		PtrStr("expected 200.000000 not to be between 200.000000 and 300.000000"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float64(200),
			Expect().Body().Float64().NotBetween(100, 200),
		),
		PtrStr("expected 200.000000 not to be between 100.000000 and 200.000000"),
	)
}
