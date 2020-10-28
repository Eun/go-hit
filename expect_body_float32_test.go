package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyFloat32_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyFloat32_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().NotEqual(200),
		),
		PtrStr("should not be 200.000000"),
	)
}

func TestExpectBodyFloat32OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().OneOf(300, 400),
		),
		PtrStr(`200.000000 should be one of []interface {}{`), PtrStr(`300.000000,`), PtrStr(`400.000000,`), PtrStr(`}`),
	)
}

func TestExpectBodyFloat32NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().NotOneOf(200, 300),
		),
		PtrStr(`200.000000 should not contain 200.000000`),
	)
}

func TestExpectBodyFloat32_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().GreaterThan(299),
		),
		PtrStr("expected 200.000000 to be greater than 299.000000"),
	)
}

func TestExpectBodyFloat32_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200.000000 to be greater or equal than 299.000000"),
	)
}

func TestExpectBodyFloat32_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().LessThan(100),
		),
		PtrStr("expected 200.000000 to be less than 100.000000"),
	)
}

func TestExpectBodyFloat32_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().LessOrEqualThan(100),
		),
		PtrStr("expected 200.000000 to be less or equal than 100.000000"),
	)
}

func TestExpectBodyFloat32_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().Between(300, 400),
		),
		PtrStr("expected 200.000000 to be between 300.000000 and 400.000000"),
	)
}

func TestExpectBodyFloat32_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Float32(200),
		Expect().Body().Float32().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().NotBetween(200, 300),
		),
		PtrStr("expected 200.000000 not to be between 200.000000 and 300.000000"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Float32(200),
			Expect().Body().Float32().NotBetween(100, 200),
		),
		PtrStr("expected 200.000000 not to be between 100.000000 and 200.000000"),
	)
}
