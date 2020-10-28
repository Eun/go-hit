package hit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Custom(func(hit Hit) {
			require.Equal(t, "Hello World", hit.Response().Body().MustString())
		}),
	)
}

func TestExpect_Double(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().Equal(`Hello World`),
		Expect().Body().String().Equal(`Hello World`),
	)
}

func TestExpect_DeepFunc(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Custom(func(h1 Hit) {
				h1.MustDo(Expect().Custom(func(h2 Hit) {
					h2.MustDo(Expect().Custom(func(h3 Hit) {
						calledFunc = true
						h3.MustDo(Expect().Body().String().Equal("Hello Universe"))
					}))
				}))
			}),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
	require.True(t, calledFunc)
}
