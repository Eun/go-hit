package hit_test

import (
	"bytes"
	"testing"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	hit.Post(t, s.URL).
		Send().Custom(func(hit hit.Hit) {
		calledFunc = true
		hit.Request().Body().SetString("Hello World")
	}).
		Expect().Body().Equal(`Hello World`).
		Do()
	require.True(t, calledFunc)
}

func TestSend_Custom_InvalidParameter(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	hit.Post(t, s.URL).
		Send(func() {
			calledFunc = true
		}).
		Do()
	require.True(t, calledFunc)
}

func TestSend_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().JSON("Hello World").
			Expect().Body().Equal(`"Hello World"`).
			Do()
	})
	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().JSON([]string{"A", "B"}).
			Expect().Body().Equal(`["A","B"]`).
			Do()
	})

	t.Run("object", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().JSON(map[string]interface{}{"A": "1", "B": "2"}).
			Expect().Body().Equal(`{"A":"1","B":"2"}`).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().JSON(8).
			Expect().Body().Equal(`8`).
			Do()
	})

	t.Run("struct", func(t *testing.T) {
		var user = struct {
			Name string
			ID   int
		}{
			"Joe",
			10,
		}

		hit.Post(t, s.URL).
			Send().JSON(user).
			Expect().Body().Equal(`{"Name":"Joe","ID":10}`).
			Do()
	})
}

func TestSend(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Request)", func(t *testing.T) {
			calledFunc := false
			hit.Post(t, s.URL).
				Send(func(e hit.Hit) {
					calledFunc = true
					e.Request().Body().SetString("Hello World")
				}).
				Expect().Body().Equal(`Hello World`).
				Do()
			require.True(t, calledFunc)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			calledFunc := false
			hit.Post(t, s.URL).
				Send(func(e hit.Hit) {
					calledFunc = true
					e.Send(`Hello World`)
				}).
				Expect().Body().Equal(`Hello World`).
				Do()
			require.True(t, calledFunc)
		})
		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			hit.Post(t, s.URL).
				Send(func() {
					calledFunc = true
				}).
				Expect().Body().Equal(``).
				Do()
			require.True(t, calledFunc)
		})
	})

	t.Run("bytes", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send([]byte("Hello World")).
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send("Hello World").
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("reader", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send(bytes.NewBufferString("Hello World")).
			Expect().Body().Equal(`Hello World`).
			Do()
	})

	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send([]string{"A", "B"}).
			Expect().Body().Equal(`["A","B"]`).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send(8).
			Expect().Body().Equal(`8`).
			Do()
	})
}

func TestSend_Double(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	require.Panics(t, func() {
		e := hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL)
		require.Equal(t, e.State(), hit.Ready)
		e.Send().Custom(func(_ hit.Hit) {
			require.Equal(t, e.State(), hit.Working)
		})
		e.Do()
		require.Equal(t, e.State(), hit.Done)
		e.Do()
		require.Equal(t, e.State(), hit.Done)
	})
}

func TestSend_Clear(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		calledFunc := false

		hit.Post(t, s.URL).
			Send(func(hit.Hit) {
				calledFunc = true
			}).
			Send().Clear().
			Do()
		require.False(t, calledFunc)
	})
}

func TestSend_AfterDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL).
				Do().
				Send("Hello")
		})
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL).
				Do().
				Send().Custom(func(hit hit.Hit) {})
		})
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Get(NewPanicWithMessage(t, PtrStr("request already fired")), s.URL).
				Do().
				Send().Headers().Set("X-Header", "Hello")
		})
	})
}
