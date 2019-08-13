package hit_test

//
// import (
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
//
// 	"crypto/md5"
// 	"encoding/hex"
//
// 	"strings"
//
// 	"io/ioutil"
//
// 	"github.com/Eun/go-hit"
// 	"github.com/stretchr/testify/require"
// )
//
// func Example_hash() {
// 	// a testing interface that panics on first error
// 	t := hit.PanicT{}
//
// 	// hashes the payload with md5 and puts the value into the Content-Signature header
// 	hashBody := func(hit hit.Hit) {
// 		hash := md5.Sum(hit.Request().Body().Bytes())
// 		hit.Request().Header.Set("Content-Signature", hex.EncodeToString(hash[:]))
// 	}
//
// 	// expectsInnerText
// 	expectInnerText := func(text string) func(hit hit.Hit) {
// 		return func(hit hit.Hit) {
// 			if !strings.Contains(hit.Response().Body().String(), text) {
// 				t.Errorf("expected %s", text)
// 				t.FailNow()
// 			}
// 		}
// 	}
//
// 	// create a server that hashes the body and checks if the payload hash is correct
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
// 		body, err := ioutil.ReadAll(request.Body)
// 		if err != nil {
// 			t.FailNow()
// 		}
// 		hash := md5.Sum(body)
//
// 		require.Equal(t, hex.EncodeToString(hash[:]), request.Header.Get("Content-Signature"))
// 		_, _ = io.WriteString(writer, `<html><body>Hello Client!</body></html>`)
// 	})
// 	s := httptest.NewServer(mux)
// 	defer s.Close()
// 	hit.Post(t, s.URL).
// 		Send().Body("Hello Server").
// 		Send().Custom(hashBody).
// 		Expect().Status(200).
// 		Expect().Custom(expectInnerText("Hello Client!")).
// 		Do()
// }
