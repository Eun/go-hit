package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestClearExpect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect("Hello Universe"),
				Expect("Hello Earth"),
				Clear().Expect(),
				Expect("Hello Nature"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect("Hello Universe"),
				Expect("Hello Earth"),
				Clear().Expect("Hello Universe"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Earth"`), nil, nil, nil, nil, nil,
		)
	})
}

func TestClearExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect("Hello Universe"))
				}),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect("Hello Earth"))
				}),
				Clear().Expect().Custom(),
				Expect().Custom(func(hit Hit) {
					hit.MustDo(Expect("Hello Nature"))
				}),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello Nature"`), nil, nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ranCustomFunc := false
		fn := func(hit Hit) {
			require.Equal(t, "Hello Universe", hit.Response().Body().String())
		}
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Custom(fn),
			Expect().Custom(fn),
			Expect().Custom(func(hit Hit) {
				ranCustomFunc = true
			}),
			Clear().Expect().Custom(fn),
		)
		require.True(t, ranCustomFunc)
	})
}
