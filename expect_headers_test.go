package hit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectHeader_Equal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header()["Date"] = nil
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers().Equal(map[string]string{"X-Header": "Hello", "Content-Length": "0"}),
		)
	})
}

func TestExpectHeaders_Contains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
		writer.Header()["Date"] = nil
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers().Contains("X-Header"),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers().Contains("X-Header2"),
			),
			PtrStr("http.Header{"), nil, nil, nil, nil, nil, nil, PtrStr(`} does not contain "X-Header2"`),
		)
	})
}

func TestExpectHeaders_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers().Get("X-Header").Equal("Hello"),
		)
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			Expect().Headers().Get("X-Header").Get("X-Header")
		})
	})
}

func TestExpectHeadersSpecificHeader_Len(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers("X-Header").Len(5),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers("X-Header").Len(0),
			),
			PtrStr(`"Hello" should have 0 item(s), but has 5`),
		)
	})
}

func TestExpectHeadersSpecificHeader_Empty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header1", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers("X-Header2").Empty(),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers("X-Header1").Empty(),
			),
			PtrStr(`"Hello" should be empty, but has 5 item(s)`),
		)
	})
}

func TestExpectSpecificHeader_Equal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers("X-Header").Equal("Hello"),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers("X-Header").Equal("Bye"),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
		)
	})
}

func TestExpectSpecificHeader_Contains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers("X-Header").Contains("Hello"),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers("X-Header").Contains("Bye"),
			),
			PtrStr(`"Hello" does not contain "Bye"`),
		)
	})
}

func TestExpectSpecificHeader_OneOf(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Headers("X-Header").OneOf("Hello", "World"),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers("X-Header").OneOf("Universe"),
			),
			PtrStr("[]interface {}{"), PtrStr(`"Universe",`), PtrStr(`} does not contain "Hello"`),
		)
	})

	t.Run("", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Headers().OneOf("Universe"),
			),
			PtrStr("OneOf can only be used if a header was already specified"),
		)
	})
}
