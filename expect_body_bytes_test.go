package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyBytes_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().Bytes().Equal([]byte("Hello World")),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().Bytes().Equal([]byte("Hello Universe")),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectBodyBytes_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().Bytes().NotEqual([]byte("Hello Universe")),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().Bytes().NotEqual([]byte("Hello World")),
		),
		PtrStr("should not be []uint8{"), nil, nil,
	)
}

func TestExpectBodyBytes_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().Bytes().Contains([]byte("World")),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().Bytes().Contains([]byte("Universe")),
		),
		PtrStr("[]uint8{"), nil, PtrStr("} does not contain []uint8{"), nil, nil,
	)
}

func TestExpectBodyBytes_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().Bytes().NotContains([]byte("Universe")),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().Bytes().NotContains([]byte("World")),
		),
		PtrStr("[]uint8{"), nil, PtrStr("} should not contain []uint8{"), nil, nil,
	)
}

func TestExpectBodyBytes_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().Bytes().Len().Equal(11),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().Bytes().Len().Equal(10),
		),
		PtrStr("not equal"), PtrStr("expected: 10"), PtrStr("actual: 11"), nil, nil, nil, nil,
	)
}
