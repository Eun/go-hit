package hit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Eun/go-hit"
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

	Test(t,
		Post(s.URL),
		Expect().Headers().Get("X-Header").Equal("Hello"),
	)
}

func TestExpectHeaders_Empty(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{}
		}),
		Expect().Headers().Empty(),
		Expect().Headers().Len(0),
	)
}
