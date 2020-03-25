package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectHeaders_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().Equal(map[string]string{"X-Header": "Hello"}),
		Expect().Header().Equal(map[string]interface{}{"X-Header": "Hello"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().Equal(map[string]string{"X-Header": "World"}),
		),
		PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectHeaders_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().NotEqual(map[string]string{"X-Header": "World"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().NotEqual(map[string]string{"X-Header": "Hello"}),
		),
		PtrStr("should not be map[string]string{"), nil, nil,
	)
}

func TestExpectHeaders_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().Contains("X-Header"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().Contains("X-Header2"),
		),
		PtrStr("http.Header{"), nil, nil, nil, PtrStr(`} does not contain "X-Header2"`),
	)
}

func TestExpectHeaders_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().NotContains("X-Header2"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().NotContains("X-Header"),
		),
		PtrStr("http.Header{"), nil, nil, nil, PtrStr(`} should not contain "X-Header"`),
	)
}

func TestExpectHeaders_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().OneOf(map[string]string{"X-Header": "Hello"}),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().OneOf(map[string]string{"X-Header": "World"}),
		),
		nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectHeaders_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().NotOneOf(map[string]string{"X-Header": "World"}),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header": []string{"Hello"},
				}
			}),
			Expect().Header().NotOneOf(map[string]string{"X-Header": "Hello"}),
		),
		nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectHeaders_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{}
		}),
		Expect().Header().Empty(),
		Expect().Header().Len(0),
	)
}
