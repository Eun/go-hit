package hash_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"crypto/md5"
	"encoding/hex"

	"strings"

	"io/ioutil"

	"fmt"

	"testing"

	. "github.com/Eun/go-hit"
)

func TestHash(t *testing.T) {
	// hashes the payload with md5 and puts the value into the Content-Signature header
	hashBody := func(hit Hit) {
		hash := md5.Sum(hit.Request().Body().MustBytes())
		hit.Request().Header.Set("Content-Signature", hex.EncodeToString(hash[:]))
	}

	// expectsInnerText
	expectInnerText := func(text string) func(hit Hit) {
		return func(hit Hit) {
			if !strings.Contains(hit.Response().Body().MustString(), text) {
				t.Error(fmt.Sprintf("expected %s", text))
			}
		}
	}

	// create a server that hashes the body and checks if the payload hash is correct
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Error()
		}
		hash := md5.Sum(body)

		if hex.EncodeToString(hash[:]) != request.Header.Get("Content-Signature") {
			t.Error(fmt.Sprintf("expected %s but got %s", hex.EncodeToString(hash[:]), request.Header.Get("Content-Signature")))
		}
		_, _ = io.WriteString(writer, `<html><body>Hello Client!</body></html>`)
	})
	s := httptest.NewServer(mux)
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello Server"),
		Send().Custom(hashBody),
		Expect().Status().Equal(200),
		Expect().Custom(expectInnerText("Hello Client!")),
	)
}
