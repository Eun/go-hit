package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func runHeaderTrailerValueExpectTest(f func(expect func(string) IExpectHeaderValue)) {
	f(func(s string) IExpectHeaderValue {
		return Expect().Headers(s).First()
	})
	f(func(s string) IExpectHeaderValue {
		return Expect().Headers(s).Last()
	})
	f(func(s string) IExpectHeaderValue {
		return Expect().Headers(s).Nth(0)
	})

	f(func(s string) IExpectHeaderValue {
		return Expect().Trailers(s).First()
	})
	f(func(s string) IExpectHeaderValue {
		return Expect().Trailers(s).Last()
	})
	f(func(s string) IExpectHeaderValue {
		return Expect().Trailers(s).Nth(0)
	})
}

func TestExpectHeaderTrailerValue_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").Len().Equal(3),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Len().Equal(1),
			),
			PtrStr("not equal"), PtrStr("expected: 1"), PtrStr("actual: 3"), nil, nil, nil, nil,
		)
	})
}

func TestExpectHeaderTrailerValue_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-Unknown").Empty(),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Empty(),
			),
			PtrStr(`"Foo" should be empty, but has 3 item(s)`),
		)
	})
}

func TestExpectHeaderTrailerValue_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").Equal("Foo"),
			expect("X-Int").Equal(3),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Equal("Bye"),
			),
			PtrStr("not equal"), PtrStr(`expected: "Bye"`), PtrStr(`actual: "Foo"`), nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").Equal(6),
			),
			PtrStr("not equal"), PtrStr(`expected: 6`), PtrStr(`actual: 3`), nil, nil, nil, nil,
		)
	})
}

func TestExpectHeaderTrailerValue_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").NotEqual("Bye"),
			expect("X-Int").NotEqual("1"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").NotEqual("Foo"),
			),
			PtrStr(`should not be "Foo"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").NotEqual("3"),
			),
			PtrStr(`should not be "3"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").NotEqual(1, 2, 3),
			),
			PtrStr(`should not be 3`),
		)
	})
}

func TestExpectHeaderTrailerValue_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").Contains("Foo"),
			expect("X-String").Contains("F", "o", "o"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Contains("Bye"),
			),
			PtrStr(`"Foo" does not contain "Bye"`),
		)
	})
}

func TestExpectHeaderTrailerValue_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").NotContains("Bye"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").NotContains("Foo"),
			),
			PtrStr(`"Foo" should not contain "Foo"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").NotContains("Bar", "Foo"),
			),
			PtrStr(`"Foo" should not contain "Foo"`),
		)
	})
}

func TestExpectHeaderTrailerValue_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").OneOf("Foo", "Bar"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").OneOf("tar", "taz"),
			),
			PtrStr(`"Foo" should be one of []interface {}{`), PtrStr(`"tar",`), PtrStr(`"taz",`), PtrStr(`}`),
		)
	})
}

func TestExpectHeaderTrailerValue_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerValueExpectTest(func(expect func(string) IExpectHeaderValue) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").NotOneOf("Bar"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").NotOneOf("Bar", "Foo"),
			),
			PtrStr(`"Foo" should not contain "Foo"`),
		)
	})
}

func TestExpectHeaderTrailerValue_OutOfBounds(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockHeadersAndTrailers(),
		Expect().Headers("X-String").Nth(-1).Equal(nil),
		Expect().Headers("X-String").Nth(100).Equal(nil),
	)
}
