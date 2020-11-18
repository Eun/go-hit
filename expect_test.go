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
		Expect().Custom(func(hit Hit) error {
			require.Equal(t, "Hello World", hit.Response().Body().MustString())
			return nil
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
			Expect().Custom(func(h1 Hit) error {
				h1.MustDo(Expect().Custom(func(h2 Hit) error {
					h2.MustDo(Expect().Custom(func(h3 Hit) error {
						calledFunc = true
						h3.MustDo(Expect().Body().String().Equal("Hello Universe"))
						return nil
					}))
					return nil
				}))
				return nil
			}),
		),
		PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
	)
	require.True(t, calledFunc)
}
