package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestClearExpectTrailer_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Trailer().Contains("X-Header1"),
				Clear().Expect().Trailer().Contains(),
				Expect().Trailer().Contains("X-Header2"),
			),
			PtrStr(`http.Header{} does not contain "X-Header2"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Trailer().Contains("X-Header1"),
				Expect().Trailer().Contains("X-Header2"),
				Clear().Expect().Trailer().Contains("X-Header1"),
			),
			PtrStr(`http.Header{} does not contain "X-Header2"`),
		)
	})
}

func TestClearExpectTrailer_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
						"X-Header2": []string{"World"},
					}
				}),
				Expect().Trailer().NotContains("X-Header1"),
				Clear().Expect().Trailer().NotContains(),
				Expect().Trailer().NotContains("X-Header2"),
			),
			nil, nil, nil, nil, nil, nil, nil, PtrStr(`} should not contain "X-Header2"`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
						"X-Header2": []string{"World"},
					}
				}),
				Expect().Trailer().NotContains("X-Header1"),
				Expect().Trailer().NotContains("X-Header2"),
				Clear().Expect().Trailer().NotContains("X-Header1"),
			),
			nil, nil, nil, nil, nil, nil, nil, PtrStr(`} should not contain "X-Header2"`))
	})
}

func TestClearExpectTrailer_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
				Clear().Expect().Trailer().OneOf(),
				Expect().Trailer().OneOf(map[string]interface{}{"X-Header2": "World"}),
			),
			nil, nil, PtrStr(`"X-Header2": "World",`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
				Expect().Trailer().OneOf(map[string]interface{}{"X-Header2": "World"}),
				Clear().Expect().Trailer().OneOf(map[string]interface{}{"X-Header1": "Hello"}),
			),
			nil, nil, PtrStr(`"X-Header2": "World",`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectTrailer_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Trailer = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Trailer().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
			Clear().Expect().Trailer().NotOneOf(),
			Expect().Trailer().NotOneOf(map[string]interface{}{"X-Header2": "World"}),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Trailer = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Trailer().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
			Expect().Trailer().NotOneOf(map[string]interface{}{"X-Header2": "World"}),
			Clear().Expect().Trailer().NotOneOf(map[string]interface{}{"X-Header1": "Hello"}),
		)
	})
}

func TestClearExpectTrailer_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Trailer().Empty(),
		Clear().Expect().Trailer().Empty(),
		Expect().Trailer().Empty(),
	)
}

func TestClearExpectTrailer_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().Len(2),
				Clear().Expect().Trailer().Len(),
				Expect().Trailer().Len(3),
			),
			nil, nil, nil, nil, PtrStr(`} should have 3 item(s), but has 1`),
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().Len(2),
				Expect().Trailer().Len(3),
				Clear().Expect().Trailer().Len(2),
			),
			nil, nil, nil, nil, PtrStr(`} should have 3 item(s), but has 1`),
		)
	})
}

func TestClearExpectTrailer_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Custom(BeforeExpectStep, func(hit Hit) {
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
				Clear().Expect().Trailer().Equal(),
				Expect().Trailer().Equal(map[string]string{
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
					hit.Response().Trailer = map[string][]string{
						"X-Header1": []string{"Hello"},
					}
				}),
				Expect().Trailer().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
				Expect().Trailer().Equal(map[string]string{
					"X-Header3": "Hello",
				}),
				Clear().Expect().Trailer().Equal(map[string]string{
					"X-Header2": "Hello",
				}),
			),
			PtrStr("Not equal"), nil, PtrStr(`"X-Header3": "Hello",`), nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})
}

func TestClearExpectTrailer_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Trailer = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Trailer().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
			Clear().Expect().Trailer().NotEqual(),
			Expect().Trailer().NotEqual(map[string]string{
				"X-Header2": "Hello",
			}),
		)
	})

	t.Run("specific", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.Response().Trailer = map[string][]string{
					"X-Header1": []string{"Hello"},
				}
			}),
			Expect().Trailer().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
			Expect().Trailer().NotEqual(map[string]string{
				"X-Header2": "Hello",
			}),
			Clear().Expect().Trailer().NotEqual(map[string]string{
				"X-Header1": "Hello",
			}),
		)
	})
}

func TestClearExpectTrailer_NotExistentStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Trailer("X-Header"),
		),
		PtrStr(`unable to find a step with Expect().Trailer("X-Header")`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Clear().Expect().Trailer(),
		),
		PtrStr(`unable to find a step with Expect().Trailer()`), PtrStr(`got these steps:`), PtrStr(`Send().Body("Hello World")`),
	)
}
