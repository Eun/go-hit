package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectTrailers_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().Equal(map[string]string{"X-Header": "Hello"}),
		Expect().Trailer().Equal(map[string]interface{}{"X-Header": "Hello"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().Equal(map[string]interface{}{"X-Header": "World"}),
		),
		PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectTrailers_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().NotEqual(map[string]string{"X-Header": "World"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().NotEqual(map[string]string{"X-Header": "Hello"}),
		),
		PtrStr("should not be map[string]string{"), nil, nil,
	)
}

func TestExpectTrailers_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().Contains("X-Header"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().Contains("X-Header2"),
		),
		PtrStr("http.Header{"), nil, nil, nil, PtrStr(`} does not contain "X-Header2"`),
	)
}

func TestExpectTrailers_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().NotContains("X-Header2"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().NotContains("X-Header"),
		),
		PtrStr("http.Header{"), nil, nil, nil, PtrStr(`} should not contain "X-Header"`),
	)
}

func TestExpectTrailers_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().OneOf(map[string]string{"X-Header": "Hello"}),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().OneOf(map[string]string{"X-Header": "World"}),
		),
		nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectTrailers_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Send().Trailer("X-Header", "Hello"),
		Expect().Trailer().NotOneOf(map[string]string{"X-Header": "World"}),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Send().Trailer("X-Header", "Hello"),
			Expect().Trailer().NotOneOf(map[string]string{"X-Header": "Hello"}),
		),
		nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectTrailers_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Trailer().Empty(),
		Expect().Trailer().Len(0),
	)
}
