package hit_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestSendSpecificHeaders_Set(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Header("X-Headers").Set("World"),
		Expect().Body().Equal("World"),
	)
}

func TestSendSpecificHeaders_Clear(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Headers().Set("X-Headers", "World"),
		Send().Headers().Clear(),
		Expect().Body().Equal(""),
	)
}

func TestSendSpecificHeaders_DoubleSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, request.Header.Get("X-Headers"))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Headers().Set("X-Headers", "World"),
		Send().Headers().Set("X-Headers", "Universe"),
		Expect().Body().Equal("Universe"),
	)
}
