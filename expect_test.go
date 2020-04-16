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

func TestExpect_DeepFunc(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Custom(func(h1 Hit) {
				h1.MustDo(Expect().Custom(func(h2 Hit) {
					h2.MustDo(Expect().Custom(func(h3 Hit) {
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
