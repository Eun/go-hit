package hit_test

import (
	"testing"

	"net/url"

	. "github.com/Eun/go-hit"
)

func mockFormValues() IStep {
	return Custom(BeforeExpectStep, func(hit Hit) {
		s := url.Values{
			"X-String":  {"Foo"},
			"X-Strings": {"Hello", "World"},
			"X-Int":     {"3"},
			"X-Ints":    {"3", "4"},
			"X-Mixed":   {"3.0", "4", "Foo"},
		}.Encode()
		hit.Request().Body().SetString(s)
		hit.Response().Body().SetString(s)
	})
}

func TestExpectFormValue_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-String").Len().Equal(1),
		Expect().Body().FormValues("X-Strings").Len().Equal(2),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").Len().Equal(1),
		),
		PtrStr("not equal"), PtrStr("expected: 1"), PtrStr("actual: 2"), nil, nil, nil, nil,
	)
}

func TestExpectFormValue_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-Unknown").Empty(),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-String").Empty(),
		),
		PtrStr(`[]string{`), PtrStr(`"Foo",`), PtrStr(`} should be empty, but has 1 item(s)`),
	)
}

func TestExpectFormValue_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-String").Equal("Foo"),
		Expect().Body().FormValues("X-Strings").Equal("Hello", "World"),
		Expect().Body().FormValues("X-Int").Equal(3),
		Expect().Body().FormValues("X-Ints").Equal(3, 4),
		Expect().Body().FormValues("X-Mixed").Equal(3.0, 4, "Foo"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-String").Equal("Bye"),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Bye",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-String").Equal("Foo", "Bar"),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`"Bar",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").Equal("Hello", "Earth"),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"Earth",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`), nil, nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").Equal("World", "Hello"),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`"World",`), PtrStr(`"Hello",`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`), nil, nil, nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Int").Equal(6),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`}`), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Int").Equal(3, 6),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`3,`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`}`), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Ints").Equal(3, 6),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`3,`), PtrStr(`6,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`4,`), PtrStr(`}`), nil, nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Ints").Equal(4, 3),
		),
		PtrStr("not equal"), PtrStr(`expected: []interface {}{`), PtrStr(`4,`), PtrStr(`3,`), PtrStr(`}`), PtrStr(`actual: []interface {}{`), PtrStr(`3,`), PtrStr(`4,`), PtrStr(`}`), nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectFormValue_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-String").NotEqual("Bye"),
		Expect().Body().FormValues("X-Strings").NotEqual("Hello"),
		Expect().Body().FormValues("X-Strings").NotEqual("Hello", "Universe"),
		Expect().Body().FormValues("X-Int").NotEqual("1"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-String").NotEqual("Foo"),
		),
		PtrStr(`should not be []interface {}{`), PtrStr(`"Foo",`), PtrStr(`}`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").NotEqual("Hello", "World"),
		),
		PtrStr(`should not be []interface {}{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`}`),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Int").NotEqual("3"),
		),
		PtrStr(`should not be []interface {}{`), PtrStr(`"3",`), PtrStr(`}`),
	)
}

func TestExpectFormValue_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-Strings").Contains("Hello"),
		Expect().Body().FormValues("X-Strings").Contains("Hello", "World"),
		Expect().Body().FormValues("X-Strings").Contains("World", "Hello"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").Contains("Bye"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} does not contain "Bye"`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").Contains("Hello", "Universe"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} does not contain "Universe"`),
	)
}

func TestExpectFormValue_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-Strings").NotContains("Bye"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").NotContains("Hello"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "Hello"`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").NotContains("Bye", "World"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "World"`),
	)
}

func TestExpectFormValue_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-Strings").OneOf("Hello", "Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").OneOf("Universe"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should be one of []interface {}{`), PtrStr(`"Universe",`), PtrStr(`}`),
	)
}

func TestExpectFormValue_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		mockFormValues(),
		Expect().Body().FormValues("X-Strings").NotOneOf("Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			mockFormValues(),
			Expect().Body().FormValues("X-Strings").NotOneOf("Hello", "World"),
		),
		PtrStr(`[]string{`), PtrStr(`"Hello",`), PtrStr(`"World",`), PtrStr(`} should not contain "Hello"`),
	)
}
