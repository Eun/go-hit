package extensibility_test

import (
	"testing"

	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/Eun/go-hit/examples/extensibility/myframework"
)

func TestExample(t *testing.T) {
	// simple echo server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		_, _ = io.Copy(writer, request.Body)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	Test(t,
		Post(server.URL),
		Send().User(User{10, "Joe ✨"}),
		Expect().User(User{10, "Joe ✨"}),
		CheckTheSparkles(),
	)
}
