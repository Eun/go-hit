package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

// for convenience we test headers and trailers here

func mockHeadersAndTrailers() IStep {
	return Custom(BeforeExpectStep, func(hit Hit) {
		m := map[string][]string{
			"X-String":  {"Foo"},
			"X-Strings": {"Hello", "World"},
			"X-Int":     {"3"},
			"X-Ints":    {"3", "4"},
			"X-Mixed":   {"3.0", "4", "Foo"},
		}
		hit.Response().Header, hit.Response().Trailer = m, m
		hit.Request().Header, hit.Request().Trailer = m, m
	})
}

func runHeaderTrailerExpectTest(f func(expect func(string) IExpectHeaders)) {
	for _, v := range []func(string) IExpectHeaders{Expect().Headers, Expect().Trailers} {
		f(v)
	}
}

func TestExpectHeaderTrailer_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").Len().Equal(1),
			expect("X-Strings").Len().Equal(2),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").Len().Equal(1),
			),
			PtrStr("not equal"), PtrStr("expected: 1"), PtrStr("actual: 2"), nil, nil, nil, nil,
		)
	})
}

func TestExpectHeaderTrailer_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
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
			PtrStr(`[]string{`), PtrStr(`"Foo",`), PtrStr(`} should be empty, but has 1 item(s)`),
		)
	})
}

func TestExpectHeaderTrailer_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").Equal("Foo"),
			expect("X-Strings").Equal("Hello", "World"),
			expect("X-Int").Equal(3),
			expect("X-Ints").Equal(3, 4),
			expect("X-Mixed").Equal(3.0, 4, "Foo"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Equal("Bye"),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Bye",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`), nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").Equal("Foo", "Bar"),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`"Bar",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`), nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").Equal("Hello", "Earth"),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"Earth",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`), nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").Equal("World", "Hello"),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"World",`), PtrStr(`"Hello",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`), nil, nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").Equal(6),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`}`), nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").Equal(3, 6),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`3,`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`}`), nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Ints").Equal(3, 6),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`3,`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`4,`), PtrStr(`}`), nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Ints").Equal(4, 3),
			),
			PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`4,`), PtrStr(`3,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`4,`), PtrStr(`}`), nil, nil, nil, nil, nil, nil,
		)
	})
}

func TestExpectHeaderTrailer_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-String").NotEqual("Bye"),
			expect("X-Strings").NotEqual("Hello"),
			expect("X-Strings").NotEqual("Hello", "Universe"),
			expect("X-Int").NotEqual("1"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-String").NotEqual("Foo"),
			),
			PtrStr(`should not be []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").NotEqual("Hello", "World"),
			),
			PtrStr(`should not be []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Int").NotEqual("3"),
			),
			PtrStr(`should not be []interface {}{`), PtrStr(`"3",`), PtrStr(`}`),
		)
	})
}

func TestExpectHeaderTrailer_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-Strings").Contains("Hello"),
			expect("X-Strings").Contains("Hello", "World"),
			expect("X-Strings").Contains("World", "Hello"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").Contains("Bye"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} does not contain "Bye"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").Contains("Hello", "Universe"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} does not contain "Universe"`),
		)
	})
}

func TestExpectHeaderTrailer_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-Strings").NotContains("Bye"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").NotContains("Hello"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "Hello"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").NotContains("Bye", "World"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "World"`),
		)
	})
}

func TestExpectHeaderTrailer_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-Strings").OneOf("Hello", "Universe"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").OneOf("Universe"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should be one of []interface {}{`), PtrStr(`"Universe",`), PtrStr(`}`),
		)
	})
}

func TestExpectHeaderTrailer_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerExpectTest(func(expect func(string) IExpectHeaders) {
		Test(t,
			Post(s.URL),
			mockHeadersAndTrailers(),
			expect("X-Strings").NotOneOf("Universe"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				mockHeadersAndTrailers(),
				expect("X-Strings").NotOneOf("Hello", "World"),
			),
			PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "Hello"`),
		)
	})
}
