package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyInt_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().Equal(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().Equal(400),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyInt_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().NotEqual(400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().NotEqual(200),
		),
		PtrStr("should not be 200"),
	)
}

func TestExpectBodyIntOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().OneOf(200, 300),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().OneOf(300, 400),
		),
		PtrStr(`200 should be one of []interface {}{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`}`),
	)
}

func TestExpectBodyIntNotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().NotOneOf(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().NotOneOf(200, 300),
		),
		PtrStr(`200 should not contain 200`),
	)
}

func TestExpectBodyInt_GreaterThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().GreaterThan(199),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().GreaterThan(299),
		),
		PtrStr("expected 200 to be greater than 299"),
	)
}

func TestExpectBodyInt_GreaterOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().GreaterOrEqualThan(200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().GreaterOrEqualThan(299),
		),
		PtrStr("expected 200 to be greater or equal than 299"),
	)
}

func TestExpectBodyInt_LessThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().LessThan(201),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().LessThan(100),
		),
		PtrStr("expected 200 to be less than 100"),
	)
}

func TestExpectBodyInt_LessOrEqualThan(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().LessOrEqualThan(200),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().LessOrEqualThan(100),
		),
		PtrStr("expected 200 to be less or equal than 100"),
	)
}

func TestExpectBodyInt_Between(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().Between(100, 200),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().Between(300, 400),
		),
		PtrStr("expected 200 to be between 300 and 400"),
	)
}

func TestExpectBodyInt_NotBetween(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().Int(200),
		Expect().Body().Int().NotBetween(300, 400),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().NotBetween(200, 300),
		),
		PtrStr("expected 200 not to be between 200 and 300"),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().Int(200),
			Expect().Body().Int().NotBetween(100, 200),
		),
		PtrStr("expected 200 not to be between 100 and 200"),
	)
}
