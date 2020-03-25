package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectSpecificTrailer_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer("X-Header").Len(5),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").Len(0),
		),
		PtrStr(`"Hello" should have 0 item(s), but has 5`),
	)
}

func TestExpectSpecificTrailer_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Trailer("X-Header").Empty(),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").Empty(),
		),
		PtrStr(`"Hello" should be empty, but has 5 item(s)`),
	)
}

func TestExpectSpecificTrailer_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-String", "Hello"),
		Send().Trailer("X-Int", 3),
		Expect().Trailer("X-String").Equal("Hello"),
		Expect().Trailer("X-Int").Equal(3),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-String", "Hello"),
			Send().Trailer("X-Int", 3),
			Expect().Trailer("X-String").Equal("Bye"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "Bye"`), nil, nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-String", "Hello"),
			Send().Trailer("X-Int", 3),
			Expect().Trailer("X-Int").Equal(1),
		),
		PtrStr("Not equal"), PtrStr("expected: 1"), nil, nil, nil, nil, nil,
	)
}

func TestExpectSpecificTrailer_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-String", "Hello"),
		Send().Trailer("X-Int", 3),
		Expect().Trailer("X-String").NotEqual("Bye"),
		Expect().Trailer("X-Int").NotEqual(1),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-String", "Hello"),
			Send().Trailer("X-Int", 3),
			Expect().Trailer("X-String").NotEqual("Hello"),
		),
		PtrStr(`should not be "Hello"`),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-String", "Hello"),
			Send().Trailer("X-Int", 3),
			Expect().Trailer("X-Int").NotEqual(3),
		),
		PtrStr(`should not be 3`),
	)
}

func TestExpectSpecificTrailer_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer("X-Header").Contains("Hello"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").Contains("Bye"),
		),
		PtrStr(`"Hello" does not contain "Bye"`),
	)
}

func TestExpectSpecificTrailer_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer("X-Header").NotContains("Bye"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").NotContains("H"),
		),
		PtrStr(`"Hello" should not contain "H"`),
	)
}

func TestExpectSpecificTrailer_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer("X-Header").OneOf("Hello", "World"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").OneOf("Universe"),
		),
		PtrStr("[]interface {}{"), PtrStr(`"Universe",`), PtrStr(`} does not contain "Hello"`),
	)
}

func TestExpectSpecificTrailer_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer("X-Header").NotOneOf("Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer("X-Header").NotOneOf("Hello", "World"),
		),
		nil, nil, nil, PtrStr(`} should not contain "Hello"`),
	)
}
