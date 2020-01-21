package hit_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"regexp"

	"github.com/Eun/go-hit/internal/minitest"
	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func EchoServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		_, _ = io.Copy(writer, request.Body)
	})
	return httptest.NewServer(mux)
}

func PrintJSONServer(jsn interface{}) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if err := json.NewEncoder(writer).Encode(jsn); err != nil {
			panic(err)
		}
	})
	return httptest.NewServer(mux)
}

var errorRegex = regexp.MustCompile(`(?s)Error:\s*(.*)Error Trace:\s*`)

func ExpectError(t *testing.T, err error, equalLines ...*string) {
	require.NotNil(t, err)
	matches := errorRegex.FindStringSubmatch(vtclean.Clean(err.Error(), false))

	require.Len(t, matches, 2)

	lines := strings.FieldsFunc(matches[1], func(r rune) bool {
		return r == '\n'
	})

	require.Equal(t, len(equalLines), len(lines), "expected: %s\nactual:   %s\n", minitest.PrintValue(equalLines), minitest.PrintValue(lines))

	for i := 0; i < len(lines); i++ {
		if equalLines[i] != nil {
			require.Equal(t, *equalLines[i], strings.TrimSpace(lines[i]))
		}
	}
}

type instantErrorReader struct {
	Error error
}

func (e instantErrorReader) Read(p []byte) (int, error) {
	return 0, e.Error
}

func PtrStr(s string) *string {
	return &s
}
