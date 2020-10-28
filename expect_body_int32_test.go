package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyInt32_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyInt32_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().NotEqual(200),
		),
		PtrStr("should not be 200"),
	)
}

func TestExpectBodyInt32OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().OneOf(300, 400),
		),
		PtrStr(`200 should be one of []interface {}{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`}`),
	)
}

func TestExpectBodyInt32NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().NotOneOf(200, 300),
		),
		PtrStr(`200 should not contain 200`),
	)
}

func TestExpectBodyInt32_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyInt32_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyInt32_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyInt32_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyInt32_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyInt32_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int32(200),
		Expect().Body().Int32().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int32(200),
			Expect().Body().Int32().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
