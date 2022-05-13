package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyString_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().Equal("Hello World"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Equal("Hello Universe"),
		),
		PtrStr("not equal"), PtrStr(`expected: "Hello Universe"`), PtrStr(`actual: "Hello World"`), nil, nil, nil, nil,
	)
}

func TestExpectBodyString_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().NotEqual("Hello Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().NotEqual("Hello World"),
		),
		PtrStr(`should not be "Hello World"`),
	)
}

func TestExpectBodyString_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().Contains("World"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Contains("Universe"),
		),
		PtrStr(`"Hello World" does not contain "Universe"`),
	)
}

func TestExpectBodyString_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().NotContains("Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().NotContains("World"),
		),
		PtrStr(`"Hello World" should not contain "World"`),
	)
}

func TestExpectBodyString_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().Len().Equal(11),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Len().Equal(10),
		),
		PtrStr("not equal"), PtrStr("expected: 10"), PtrStr("actual: 11"), nil, nil, nil, nil,
	)
}

func TestExpectBodyString_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().OneOf("Hello World", "Hello Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().OneOf("Hello Universe", "Hello Earth"),
		),
		nil, nil, nil, nil,
	)
}

func TestExpectBodyString_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().NotOneOf("Hello Universe", "Hello Earth"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().NotOneOf("Hello World", "Hello Universe"),
		),
		PtrStr(`"Hello World" should not contain "Hello World"`),
	)
}
