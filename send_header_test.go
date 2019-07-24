package hit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/Eun/go-hit"
)

func TestSendHeaders_Set(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Headers().Set("X-Headers", "World").
			Expect().Body().Equal("World").
			Do()
	})
}

func TestSendHeaders_Clear(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Headers().Set("X-Headers", "World").
			Send().Headers().Clear().
			Expect().Body().Equal("").
			Do()
	})
}

func TestSendHeaders_DoubleSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Headers().Set("X-Headers", "World").
			Send().Headers().Set("X-Headers", "Universe").
			Expect().Body().Equal("Universe").
			Do()
	})
}

func TestSendHeaders_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Headers().Set("X-Headers", "World").
			Send().Headers().Delete("X-Headers").
			Expect().Body().Equal("").
			Do()
	})
}
