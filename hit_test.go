package hit_test

import (
	"bytes"
	"net/http"
	"testing"

	"encoding/json"

	"io/ioutil"
	"strings"

	"net/url"

	"github.com/Eun/go-hit"
	"github.com/Eun/go-hit/expr"
	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func Test_Debug(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("no json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		hit.Post(t, s.URL).SetStdout(buf).Send([]byte{0, 1, 2, 3}).Do().Debug()

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, "AAECAw==", expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		hit.Post(t, s.URL).SetStdout(buf).Send([]int{1, 2, 3}).Do().Debug()

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug without response", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		hit.Post(t, s.URL).SetStdout(buf).Send([]int{1, 2, 3}).Debug().Do()

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Nil(t, expr.MustGetValue(m, "Request.Body", expr.IgnoreNotFound, expr.IgnoreError))
		require.Nil(t, expr.MustGetValue(m, "Response", expr.IgnoreNotFound, expr.IgnoreError))
	})
}

func Test_Stdout(t *testing.T) {
	h := hit.Post(t, "")
	buf := bytes.NewBuffer(nil)
	h.SetStdout(buf)
	require.Equal(t, buf, h.Stdout())
}

func TestWithRequest(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	req, err := http.NewRequest("POST", s.URL, bytes.NewReader([]byte("Hello World")))
	require.NoError(t, err)
	hit.WithRequest(t, req).
		Do().
		Expect().Body("Hello World")
}

func TestCopy(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("expect/send", func(t *testing.T) {
		calledFunc := 0
		template := hit.Post(t, s.URL).
			Send().Custom(func(hit hit.Hit) {
			hit.Send().Headers().Set("Content-Type", "application/json")
			calledFunc++
		}).
			Expect().Custom(func(hit hit.Hit) {
			hit.Expect().Headers("Content-Type").Equal("application/json")
			calledFunc++
		})

		template.Copy().
			Send("Hello World").
			Expect("Hello World").
			Do()

		template.Copy().
			Send("Bye World").
			Expect("Bye World").
			Do()

		require.Equal(t, 4, calledFunc)
	})

	t.Run("request", func(t *testing.T) {
		t.Run("body", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
			require.NoError(t, err)
			template := hit.WithRequest(t, req)

			template.Copy().
				Expect("Hello World").
				Do()

			// test if body is still available
			template.Copy().
				Expect("Hello World").
				Do()

			// test if we can still consume original body
			template.Expect("Hello World").
				Do()
		})
		t.Run("header/trailer", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
			require.NoError(t, err)
			req.Header.Set("X-Header", "Foo")
			req.Trailer = map[string][]string{"X-Trailer": {"Bar"}}
			template := hit.WithRequest(t, req)

			require.Equal(t, []string{"Foo"}, template.Request().Header["X-Header"])
			require.Equal(t, []string{"Bar"}, template.Request().Trailer["X-Trailer"])

			require.Equal(t, []string{"Foo"}, template.Copy().Request().Header["X-Header"])
			require.Equal(t, []string{"Bar"}, template.Copy().Request().Trailer["X-Trailer"])
		})

		t.Run("PostForm/Form", func(t *testing.T) {
			u, err := url.Parse(s.URL)
			require.NoError(t, err)

			req := &http.Request{
				URL:    u,
				Method: "POST",
				Header: map[string][]string{
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
				Body: ioutil.NopCloser(strings.NewReader("a=1&a=2&a=banana")),
			}
			require.NoError(t, req.ParseForm())

			template := hit.WithRequest(t, req)

			require.Equal(t, []string{"1", "2", "banana"}, template.Request().Form["a"])
			require.Equal(t, []string{"1", "2", "banana"}, template.Request().PostForm["a"])
			require.Equal(t, []string{"1", "2", "banana"}, template.Copy().Request().Form["a"])
			require.Equal(t, []string{"1", "2", "banana"}, template.Copy().Request().PostForm["a"])
		})

		t.Run("MultipartForm", func(t *testing.T) {
			u, err := url.Parse(s.URL)
			require.NoError(t, err)

			req := &http.Request{
				URL:    u,
				Method: "POST",
				Header: map[string][]string{
					"Content-Type": {`multipart/form-data; boundary="foo123"`},
				},
				Body: ioutil.NopCloser(strings.NewReader("--foo123\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n--foo123\r\nContent-Disposition: form-data; name=\"file1\"; filename=\"file1\"\r\n\r\nbaz\r\n--foo123\r\n")),
			}
			require.NoError(t, req.ParseMultipartForm(10000))

			template := hit.WithRequest(t, req)

			require.Equal(t, []string{"bar"}, template.Request().MultipartForm.Value["foo"])
			require.Equal(t, []string{"bar"}, template.Copy().Request().MultipartForm.Value["foo"])

			require.Len(t, template.Request().MultipartForm.File["file1"], 1)
			require.Equal(t, "file1", template.Request().MultipartForm.File["file1"][0].Filename)
			f, err := template.Request().MultipartForm.File["file1"][0].Open()
			require.NoError(t, err)
			defer f.Close()

			buf, err := ioutil.ReadAll(f)
			require.NoError(t, err)
			require.Equal(t, "baz", string(buf))

			require.Len(t, template.Copy().Request().MultipartForm.File["file1"], 1)
			require.Equal(t, "file1", template.Copy().Request().MultipartForm.File["file1"][0].Filename)
			f, err = template.Copy().Request().MultipartForm.File["file1"][0].Open()
			require.NoError(t, err)
			defer f.Close()

			buf, err = ioutil.ReadAll(f)
			require.NoError(t, err)
			require.Equal(t, "baz", string(buf))
		})

		t.Run("TransferEncoding", func(t *testing.T) {
			u, err := url.Parse(s.URL)
			require.NoError(t, err)

			req := &http.Request{
				URL:              u,
				Method:           "POST",
				TransferEncoding: []string{"a", "b"},
			}
			template := hit.WithRequest(t, req)

			require.Equal(t, []string{"a", "b"}, template.Request().TransferEncoding)
			require.Equal(t, []string{"a", "b"}, template.Copy().Request().TransferEncoding)
		})
	})
}

func TestMutate(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	v := hit.Post(t, s.URL)

	v.Connect(s.URL)
	require.Equal(t, http.MethodConnect, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Delete(s.URL)
	require.Equal(t, http.MethodDelete, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Get(s.URL)
	require.Equal(t, http.MethodGet, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Head(s.URL)
	require.Equal(t, http.MethodHead, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Post(s.URL)
	require.Equal(t, http.MethodPost, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Options(s.URL)
	require.Equal(t, http.MethodOptions, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Put(s.URL)
	require.Equal(t, http.MethodPut, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Trace(s.URL)
	require.Equal(t, http.MethodTrace, v.Request().Method)
	require.Equal(t, s.URL, v.Request().URL.String())

	v.Do()

	require.Equal(t, hit.Done, v.State())
}
