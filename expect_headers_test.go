package hit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eun/go-hit"
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
		hit.Post(t, s.URL).
			Expect().Headers().Equal(map[string]string{"X-Header": "Hello", "Content-Length": "0"}).
			Do()
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
		hit.Post(t, s.URL).
			Expect().Headers().Contains("X-Header").
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("http.Header{"), nil, nil, nil, nil, nil, nil, PtrStr(`} does not contain "X-Header2"`)), s.URL).
				Expect().Headers().Contains("X-Header2").
				Do()
		})
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
		hit.Post(t, s.URL).
			Expect().Headers().Get("X-Header").Equal("Hello").
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("Get can only be used if no header was already specified")), s.URL).
				Expect().Headers().Get("X-Header").Get("X-Header").Equal("Hello").
				Do()
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
		hit.Post(t, s.URL).
			Expect().Headers("X-Header").Len(5).
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr(`"Hello" should have 0 item(s), but has 5`)), s.URL).
				Expect().Headers("X-Header").Len(0).
				Do()
		})
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
		hit.Post(t, s.URL).
			Expect().Headers("X-Header2").Empty().
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr(`"Hello" should be empty, but has 5 item(s)`)), s.URL).
				Expect().Headers("X-Header1").Empty().
				Do()
		})
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
		hit.Post(t, s.URL).
			Expect().Headers("X-Header").Equal("Hello").
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("Not equal"), nil, nil, nil, nil, nil, nil), s.URL).
				Expect().Headers("X-Header").Equal("Bye").
				Do()
		})
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
		hit.Post(t, s.URL).
			Expect().Headers("X-Header").Contains("Hello").
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr(`"Hello" does not contain "Bye"`)), s.URL).
				Expect().Headers("X-Header").Contains("Bye").
				Do()
		})
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
		hit.Post(t, s.URL).
			Expect().Headers("X-Header").OneOf("Hello", "World").
			Do()
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("[]interface {}{"), PtrStr(`"Universe",`), PtrStr(`} does not contain "Hello"`)), s.URL).
				Expect().Headers("X-Header").OneOf("Universe").
				Do()
		})
	})

	t.Run("", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr("OneOf can only be used if a header was already specified")), s.URL).
				Expect().Headers().OneOf("Universe").
				Do()
		})
	})
}
