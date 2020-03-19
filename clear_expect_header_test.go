package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClearExpectHeader_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Header().Contains("X-Header1"),
				Clear().Expect().Header().Contains(),
				Expect().Header().Contains("X-Header2"),
			),
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, PtrStr(`} does not contain "X-Header2"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Header().Contains("X-Header1"),
				Expect().Header().Contains("X-Header2"),
				Clear().Expect().Header().Contains("X-Header1"),
			),
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, PtrStr(`} does not contain "X-Header2"`))
	})
}

func TestClearExpectHeader_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
						"X-Header2": []string{"World"},
					}
				}),
				Expect().Header().NotContains("X-Header1"),
				Clear().Expect().Header().NotContains(),
				Expect().Header().NotContains("X-Header2"),
			),
			nil, nil, nil, nil, nil, nil, nil, PtrStr(`} should not contain "X-Header2"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
						"X-Header2": []string{"World"},
					}
				}),
				Expect().Header().NotContains("X-Header1"),
				Expect().Header().NotContains("X-Header2"),
				Clear().Expect().Header().NotContains("X-Header1"),
			),
			nil, nil, nil, nil, nil, nil, nil, PtrStr(`} should not contain "X-Header2"`))
	})
}

func TestClearExpectHeader_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
				Clear().Expect().Header().OneOf(),
				Expect().Header().OneOf(map[string]interface{}{"X-Header2": "World"}),
			),
			nil, nil, PtrStr(`"X-Header2": "World",`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
				Expect().Header().OneOf(map[string]interface{}{"X-Header2": "World"}),
				Clear().Expect().Header().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
			),
			nil, nil, PtrStr(`"X-Header2": "World",`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectHeader_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Header().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
			Clear().Expect().Header().NotOneOf(),
			Expect().Header().NotOneOf(map[string]interface{}{"X-Header2": "World"}),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Header().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
			Expect().Header().NotOneOf(map[string]interface{}{"X-Header2": "World"}),
			Clear().Expect().Header().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
		)
	})
}

func TestClearExpectHeader_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{}
		}),
		Expect().Header().Empty(),
		Clear().Expect().Header().Empty(),
		Expect().Header().Empty(),
	)
}

func TestClearExpectHeader_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().Len(2),
				Clear().Expect().Header().Len(),
				Expect().Header().Len(3),
			),
			nil, nil, nil, nil, PtrStr(`} should have 3 item(s), but has 1`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().Len(2),
				Expect().Header().Len(3),
				Clear().Expect().Header().Len(2),
			),
			nil, nil, nil, nil, PtrStr(`} should have 3 item(s), but has 1`),
		)
	})
}

func TestClearExpectHeader_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
				Clear().Expect().Header().Equal(),
				Expect().Header().Equal(map[string]string{
					"X-Header3": "Hello",
				}),
			),
			PtrStr("Not equal"), nil, PtrStr(`"X-Header3": "Hello",`), nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Header = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Header().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
				Expect().Header().Equal(map[string]string{
					"X-Header3": "Hello",
				}),
				Clear().Expect().Header().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
			),
			PtrStr("Not equal"), nil, PtrStr(`"X-Header3": "Hello",`), nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})
}

func TestClearExpectHeader_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Header().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
			Clear().Expect().Header().NotEqual(),
			Expect().Header().NotEqual(map[string]string{
				"X-Header2": "Hello",
			}),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Header = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Header().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
			Expect().Header().NotEqual(map[string]string{
				"X-Header2": "Hello",
			}),
			Clear().Expect().Header().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
		)
	})
}

func TestClearExpectHeader_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Header("X-Header"),
		),
		PtrStr(`unable to find a step with Expect().Header("X-Header")`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Header(),
		),
		PtrStr(`unable to find a step with Expect().Header()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
