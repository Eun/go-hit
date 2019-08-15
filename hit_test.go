package hit_test

import (
	"net/http"
	"testing"

	"fmt"
	"strings"

	"errors"

	"bytes"

	"io/ioutil"
	"net/url"

	"encoding/json"

	. "github.com/Eun/go-hit"
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
		Test(t,
			Post(s.URL),
			SetStdout(buf),
			Send("Hello World"),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		Test(t,
			Post(s.URL),
			SetStdout(buf),
			Send([]int{1, 2, 3}),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug without body", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		Test(t,
			Post(s.URL),
			SetStdout(buf),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Nil(t, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})
}

//
func Test_Stdout(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		SetStdout(buf),
		Custom(BeforeSendStep, func(hit Hit) {
			require.Equal(t, buf, hit.Stdout())
		}),
	)
}

func TestSetRequest(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	req, err := http.NewRequest("POST", s.URL, bytes.NewReader([]byte("Hello World")))
	require.NoError(t, err)
	Test(t,
		SetRequest(req),
		Expect().Body("Hello World"),
	)
}

func TestMultiUse(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)

		Test(t,
			SetRequest(req),
			Expect().Body("Hello World"),
		)

		Test(t,
			SetRequest(req),
			Expect().Body("Hello World"),
		)
	})
	t.Run("header/trailer", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)
		req.Header.Set("X-Header", "Foo")
		req.Trailer = map[string][]string{"X-Trailer": {"Bar"}}

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"Foo"}, hit.Request().Header["X-Header"])
				require.Equal(t, []string{"Bar"}, hit.Request().Trailer["X-Trailer"])
			}),
		)

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"Foo"}, hit.Request().Header["X-Header"])
				require.Equal(t, []string{"Bar"}, hit.Request().Trailer["X-Trailer"])
			}),
		)
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

		require.Equal(t, []string{"1", "2", "banana"}, req.Form["a"])
		require.Equal(t, []string{"1", "2", "banana"}, req.PostForm["a"])

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().Form["a"])
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().PostForm["a"])
			}),
		)

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().Form["a"])
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().PostForm["a"])
			}),
		)
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

		require.Equal(t, []string{"bar"}, req.MultipartForm.Value["foo"])

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"bar"}, hit.Request().MultipartForm.Value["foo"])
			}),
		)

		require.Len(t, req.MultipartForm.File["file1"], 1)

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, "file1", hit.Request().MultipartForm.File["file1"][0].Filename)
			}),
		)
		f, err := req.MultipartForm.File["file1"][0].Open()
		require.NoError(t, err)
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, "baz", string(buf))

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Len(t, hit.Request().MultipartForm.File["file1"], 1)
				require.Equal(t, "file1", hit.Request().MultipartForm.File["file1"][0].Filename)
				f, err = hit.Request().MultipartForm.File["file1"][0].Open()
				require.NoError(t, err)
				defer f.Close()
				buf, err = ioutil.ReadAll(f)
				require.NoError(t, err)
				require.Equal(t, "baz", string(buf))
			}),
		)
	})

	t.Run("TransferEncoding", func(t *testing.T) {
		u, err := url.Parse(s.URL)
		require.NoError(t, err)

		req := &http.Request{
			URL:              u,
			Method:           "POST",
			TransferEncoding: []string{"a", "b"},
		}

		require.Equal(t, []string{"a", "b"}, req.TransferEncoding)

		Test(t,
			SetRequest(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"a", "b"}, hit.Request().TransferEncoding)
			}),
		)
	})

	t.Run("with embedded request", func(t *testing.T) {
		s := EchoServer()
		defer s.Close()
		template := []IStep{
			Post(s.URL),
			Send().Header("Content-Type").Set("application/json"),
			Expect().Header("Content-Type").Equal("application/json"),
		}
		Test(t,
			append(template,
				Send().Body().JSON("Hello World"),
				Expect().Body().JSON("Hello World"),
			)...,
		)

		Test(t,
			append(template,
				Send().Body().JSON("Hello Universe"),
				Expect().Body().JSON("Hello Universe"),
			)...,
		)
	})
}

func TestBaseURL(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		SetBaseURL(s.URL),
		Get("/"),
		Expect().Status(http.StatusOK),
	)
}

func TestFormatURL(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Get("%s", s.URL),
		Expect().Status(http.StatusOK),
	)
}

type DummyTestingT struct {
	Args []interface{}
}

func (d *DummyTestingT) Error(args ...interface{}) {
	d.Args = args
}

func TestTestError(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	d := &DummyTestingT{}
	Test(d,
		Get(s.URL),
		Expect().Status(http.StatusNotFound),
	)

	var sb strings.Builder
	for i := 0; i < len(d.Args); i++ {
		fmt.Fprintln(&sb, d.Args[i])
	}

	ExpectError(t, errors.New(sb.String()), PtrStr("Expected status code to be 404 but was 200 instead"))
}

func TestSetHTTPClient(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	client := &http.Client{}

	Test(t,
		Post(s.URL),
		SetHTTPClient(client),
		Custom(AfterSendStep, func(hit Hit) {
			require.Equal(t, client, hit.HTTPClient())
		}),
	)
}
