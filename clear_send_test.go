package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestClearSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Custom(func(hit Hit) {
					hit.MustDo(Send("Hello World"))
				}),
				Clear().Send().Custom(),
				Send().Custom(func(hit Hit) {
					hit.MustDo(Send("Hello Earth"))
				}),
				Expect("Hello World"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hello World"`), PtrStr(`actual: "Hello Earth"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific", func(t *testing.T) {
		ranCustomFunc := false
		fn := func(hit Hit) {
			hit.Request().Body().SetString("Hello Earth")
		}
		Test(t,
			Post(s.URL),
			Send().Custom(fn),
			Clear().Expect().Custom(fn),
			Send().Custom(func(hit Hit) {
				ranCustomFunc = true
				hit.Request().Body().SetString("Hello World")
			}),
			Expect("Hello World"),
		)
		require.True(t, ranCustomFunc)
	})
}
