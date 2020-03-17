package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Custom(func(hit Hit) {
			require.Equal(t, "Hello World", hit.Response().Body().String())
		}),
	)
}

func TestExpect_Double(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Body().Equal(`Hello World`),
		Expect().Body().Equal(`Hello World`),
	)
}

func TestExpect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Response)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) {
						if hit.Response().Body().String() != "Hello Universe" {
							panic("Not equal")
						}
					}),
				),
				PtrStr("Not equal"),
			)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) {
						hit.MustDo(Expect("Hello Universe"))
					})),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
			)
		})
		t.Run("with correct parameter (using Hit) and error", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) error {
						return hit.Do(Expect("Hello Universe"))
					})),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
			)
		})

		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body("Hello World"),
				Expect(func() {
					calledFunc = true
				}),
			)
			require.True(t, calledFunc)
		})
	})

	t.Run("body", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect("Hello World"),
		)
	})
}

func TestExpect_DeepFunc(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect(func(h1 Hit) {
				h1.MustDo(Expect(func(h2 Hit) {
					h2.MustDo(Expect(func(h3 Hit) {
						calledFunc = true
						h3.MustDo(Expect().Body().Equal("Hello Universe"))
					}))
				}))
			}),
		),
		PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
	)
	require.True(t, calledFunc)
}

func TestExpectFinal(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	t.Run("Expect(value).Body()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Body())
		})
	})

	t.Run("Expect(value).Interface()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Interface(nil))
		})
	})

	t.Run("Expect(value).Header()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Header().Empty())
		})
	})

	t.Run("Expect(value).Status()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Status())
		})
	})

	t.Run("Expect(value).Custom()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Custom(nil))
		})
	})
}
