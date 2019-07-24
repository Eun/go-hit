package hit_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

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

type PanicWithMessage struct {
	T     *testing.T
	Lines []*string
}

func NewPanicWithMessage(t *testing.T, lines ...*string) *PanicWithMessage {
	return &PanicWithMessage{
		T:     t,
		Lines: lines,
	}
}

func (p *PanicWithMessage) Errorf(format string, args ...interface{}) {
	panic("never be called")
}

func (p *PanicWithMessage) FailNow() {
	panic("never be called")
}

func (p *PanicWithMessage) ErrorMessage(msg, detail string) {
	var sb strings.Builder

	if msg != "" {
		sb.WriteString(msg)
	}

	if detail != "" {
		if sb.Len() > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(detail)
	}

	lines := strings.FieldsFunc(sb.String(), func(r rune) bool {
		return r == '\n'
	})

	require.Len(p.T, lines, len(p.Lines))

	for i := 0; i < len(lines); i++ {
		if p.Lines[i] != nil {
			require.Equal(p.T, *p.Lines[i], strings.TrimSpace(lines[i]))
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
