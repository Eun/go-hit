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

	Test(t,
		Post(s.URL),
		Expect().Header().Equal(map[string]string{"X-Header": "Hello", "Content-Length": "0"}),
		Expect().Header().Equal(map[string]interface{}{"X-Header": "Hello", "Content-Length": "0"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header().Equal(map[string]string{"X-Header": "World", "Content-Length": "0"}),
		),
		PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
}

func TestExpectHeader_NotEqual(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header()["Date"] = nil
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header().NotEqual(map[string]string{"X-Header": "World", "Content-Length": "0"}),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header().NotEqual(map[string]string{"X-Header": "Hello", "Content-Length": "0"}),
		),
		PtrStr("should not be map[string]string{"), nil, nil, nil,
	)
}

func TestExpectHeaders_Contains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
		writer.Header()["Date"] = nil
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header().Contains("X-Header"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header().Contains("X-Header2"),
		),
		PtrStr("http.Header{"), nil, nil, nil, nil, nil, nil, PtrStr(`} does not contain "X-Header2"`),
	)
}

func TestExpectHeaders_NotContains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
		writer.Header()["Date"] = nil
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header().NotContains("X-Header2"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header().NotContains("X-Header"),
		),
		PtrStr("http.Header{"), nil, nil, nil, nil, nil, nil, PtrStr(`} should not contain "X-Header"`),
	)
}

func TestExpectHeaders_OneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().OneOf(map[string]string{"X-Header": "Hello"}),
	)
}

func TestExpectHeaders_NotOneOf(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Custom(BeforeExpectStep, func(hit Hit) {
			hit.Response().Header = map[string][]string{
				"X-Header": []string{"Hello"},
			}
		}),
		Expect().Header().NotOneOf(map[string]string{"X-Header": "World"}),
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
		Expect().Header().Empty(),
		Expect().Header().Len(0),
	)
}
