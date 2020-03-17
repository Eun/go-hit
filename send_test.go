package hit_test

import (
	"testing"

	"io"
	"net/http"
	"net/http/httptest"

	"bytes"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Request)", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send(func(e Hit) {
					calledFunc = true
					e.Request().Body().SetString("Hello World")
				}),
				Expect().Body().Equal(`Hello World`),
			)
			require.True(t, calledFunc)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send(func(e Hit) {
					calledFunc = true
					e.MustDo(Send(`Hello World`))
				}),
				Expect().Body().Equal(`Hello World`),
			)
			require.True(t, calledFunc)
		})
		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send(func() {
					calledFunc = true
				}),
				Expect().Body().Equal(``),
			)
			require.True(t, calledFunc)
		})
	})

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send([]byte("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send("Hello World"),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("reader", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send(bytes.NewBufferString("Hello World")),
			Expect().Body().Equal(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send(8),
			Expect().Body().Equal(`8`),
		)
	})
}

func TestSend_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("inline", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) {
				calledFunc = true
				hit.Request().Body().SetString("Hello World")
			}),
			Expect().Body().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
	t.Run("MustDo", func(t *testing.T) {
		calledFunc := false
		Test(t,
			Post(s.URL),
			Send().Custom(func(hit Hit) {
				calledFunc = true
				hit.MustDo(Send("Hello World"))
			}),
			Expect().Body().Equal(`Hello World`),
		)
		require.True(t, calledFunc)
	})
}

func TestSend_Custom_InvalidParameter(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	Test(t,
		Post(s.URL),
		Send(func() {
			calledFunc = true
		}),
	)
	require.True(t, calledFunc)
}

func TestSend_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().JSON("Hello World"),
			Expect().Body().Equal(`"Hello World"`),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().JSON([]string{"A", "B"}),
			Expect().Body().Equal(`["A","B"]`),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().JSON(map[string]interface{}{"A": "1", "B": "2"}),
			Expect().Body().Equal(`{"A":"1","B":"2"}`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().JSON(8),
			Expect().Body().Equal(`8`),
		)
	})

	t.Run("struct", func(t *testing.T) {
		var user = struct {
			Name string
			ID   int
		}{
			"Joe",
			10,
		}

		Test(t,
			Post(s.URL),
			Send().JSON(user),
			Expect().Body().Equal(`{"Name":"Joe","ID":10}`),
		)
	})
}

func TestSendHeader(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Header("X-Headers", "World"),
		Expect().Body().Equal("World"),
	)
}

func TestSendHeader_DoubleSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Header("X-Headers", "World"),
		Send().Header("X-Headers", "Universe"),
		Expect().Body().Equal("Universe"),
	)
}
