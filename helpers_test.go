package hit_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Eun/go-hit/errortrace"

	"github.com/Eun/go-hit"

	"github.com/stretchr/testify/require"

	"github.com/Eun/go-hit/internal/minitest"
)

func EchoServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header()["Date"] = nil
		for k, v := range request.Header {
			writer.Header()[k] = v
		}

		for k := range request.Trailer {
			writer.Header().Add("Trailer", k)
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = io.Copy(writer, request.Body)

		for k, v := range request.Trailer {
			writer.Header()[k] = v
		}
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

var whiteSpaceRegex = regexp.MustCompile(`\s+`)

func ExpectError(t *testing.T, err error, equalLines ...*string) {
	require.NotNil(t, err)

	hitErr, ok := err.(*hit.Error)
	if !ok {
		require.FailNow(t, "err is not *hit.Error")
	}
	etError, ok := hitErr.Unwrap().(*errortrace.ErrorTrace)
	if !ok {
		require.FailNow(t, "err is not *errortrace.ErrorTrace")
	}

	lines := strings.FieldsFunc(etError.ErrorText(), func(r rune) bool {
		return r == '\n'
	})

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(whiteSpaceRegex.ReplaceAllString(lines[i], " "))
	}
	for i := 0; i < len(equalLines); i++ {
		if equalLines[i] != nil {
			equalLines[i] = PtrStr(strings.TrimSpace(whiteSpaceRegex.ReplaceAllString(*equalLines[i], " ")))
		}
	}

	require.Equal(t, len(equalLines), len(lines), "expected: %s\nactual:   %s\n", minitest.PrintValue(equalLines), minitest.PrintValue(lines))

	for i := 0; i < len(lines); i++ {
		if equalLines[i] != nil {
			require.Equal(t,
				*equalLines[i],
				lines[i])
		}
	}
}

func PtrStr(s string) *string {
	return &s
}
