package hit_test

import (
	"testing"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	hit.Post(t, s.URL).
		Send().Body("Hello World").
		Expect().Custom(func(hit hit.Hit) {
		require.Equal(t, "Hello World", hit.Response().Body().String())
	}).
		Do()
}

func TestExpect_Double(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	hit.Post(t, s.URL).
		Send().Body("Hello World").
		Expect().Body().Equal(`Hello World`).
		Expect().Body().Equal(`Hello World`).
		Do()
}

func TestExpect_Clear(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	hit.Post(t, s.URL).
		Send().Body("Hello World").
		Expect().Body().Equal(`Hello Universe`).
		Expect().Clear().
		Expect().Body().Equal(`Hello World`).
		Do()
}

func TestExpect_AfterDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	hit.Get(t, s.URL).
		Send("Hello").
		Do().
		Expect().Body().Equal("Hello")
}

func TestExpect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Response)", func(t *testing.T) {
			require.Panics(t, func() {
				hit.Post(NewPanicWithMessage(t, PtrStr("Not Equal")), s.URL).
					Send().Body("Hello World").
					Expect(func(hit hit.Hit) {
						require.Equal(hit.T(), "Hello Universe", hit.Response().Body().String())
					}).
					Do()
			})
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			require.Panics(t, func() {
				hit.Post(NewPanicWithMessage(t, PtrStr("Not Equal")), s.URL).
					Send().Body("Hello World").
					Expect(func(hit hit.Hit) {
						hit.Expect("Hello Universe")
					}).
					Do()
			})
		})

		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			hit.Post(t, s.URL).
				Send().Body("Hello World").
				Expect(func() {
					calledFunc = true
				}).
				Do()
			require.True(t, calledFunc)
		})
	})

	t.Run("body", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body("Hello World").
			Expect("Hello World").
			Do()
	})
}

func TestExpect_DeepFunc(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Before Do", func(t *testing.T) {
		calledFunc := false
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("Not Equal")), s.URL).
				Send().Body("Hello World").
				Expect(func(h1 hit.Hit) {
					h1.Expect(func(h2 hit.Hit) {
						h2.Expect(func(h3 hit.Hit) {
							calledFunc = true
							h3.Expect().Body().Equal("Hello Universe")
						})
					})
				}).
				Do()
		})
		require.True(t, calledFunc)
	})

	t.Run("After Do", func(t *testing.T) {
		calledFunc := false
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("Not Equal")), s.URL).
				Send().Body("Hello World").
				Do().
				Expect(func(h1 hit.Hit) {
					h1.Expect(func(h2 hit.Hit) {
						h2.Expect(func(h3 hit.Hit) {
							calledFunc = true
							h3.Expect().Body().Equal("Hello Universe")
						})
					})
				})
		})
		require.True(t, calledFunc)
	})
}
