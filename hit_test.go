package hit_test

import (
	"net/http"
	"testing"

	"strings"

	"bytes"

	"io/ioutil"
	"net/url"

	"encoding/json"

	"net/http/httptest"

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

		Test(t,
			Post(s.URL),
			Stdout(buf),
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

		Test(t,
			Post(s.URL),
			Stdout(buf),
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

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Nil(t, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Debug("Request"),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("debug in custom", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.MustDo(Debug("Request"))
			}),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})
}

//
func Test_Stdout(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Custom(BeforeSendStep, func(hit Hit) {
			require.Equal(t, buf, hit.Stdout())
		}),
	)
}

func TestRequest(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	req, err := http.NewRequest("POST", s.URL, bytes.NewReader([]byte("Hello World")))
	require.NoError(t, err)
	Test(t,
		Request(req),
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
			Request(req),
			Expect().Body("Hello World"),
		)

		Test(t,
			Request(req),
			Expect().Body("Hello World"),
		)
	})
	t.Run("header/trailer", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)
		req.Header.Set("X-Header", "Foo")
		req.Trailer = map[string][]string{"X-Trailer": {"Bar"}}

		Test(t,
			Request(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"Foo"}, hit.Request().Header["X-Header"])
				require.Equal(t, []string{"Bar"}, hit.Request().Trailer["X-Trailer"])
			}),
		)

		Test(t,
			Request(req),
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
			Request(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().Form["a"])
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().PostForm["a"])
			}),
		)

		Test(t,
			Request(req),
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
			Request(req),
			Custom(AfterSendStep, func(hit Hit) {
				require.Equal(t, []string{"bar"}, hit.Request().MultipartForm.Value["foo"])
			}),
		)

		require.Len(t, req.MultipartForm.File["file1"], 1)

		Test(t,
			Request(req),
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
			Request(req),
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
			Request(req),
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
			Send().Header("Content-Type", "application/json"),
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
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		BaseURL("%s/foo", s.URL),
		Get("/"),
		Expect().Status(http.StatusNoContent),
	)
}

func TestFormatURL(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Get("%s/foo", s.URL),
		Expect().Status(http.StatusNoContent),
	)
}

func TestHTTPClient(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	client := &http.Client{}

	Test(t,
		Post(s.URL),
		HTTPClient(client),
		Custom(AfterSendStep, func(hit Hit) {
			require.Equal(t, client, hit.HTTPClient())
		}),
	)
}

func TestCombineSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			CombineSteps(
				Post(s.URL),
				Send("Hello"),
				Expect("Hello"),
			),
			Clear().Expect(),
			Expect("World"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "World"`), nil, nil, nil, nil, nil,
	)
}

func TestCombineSteps_DoubleExecution(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("send", func(t *testing.T) {
		calls := 0
		Test(
			t,
			Post(s.URL),
			CombineSteps(
				Send().Custom(func(hit Hit) {
					calls++
				}),
			),
		)
		require.Equal(t, 1, calls)
	})
	t.Run("expect", func(t *testing.T) {
		calls := 0
		Test(
			t,
			Post(s.URL),
			CombineSteps(
				Expect().Custom(func(hit Hit) {
					calls++
				}),
			),
		)
		require.Equal(t, 1, calls)
	})
	t.Run("other", func(t *testing.T) {
		calls := 0
		Test(
			t,
			Post(s.URL),
			CombineSteps(
				Custom(BeforeSendStep, func(hit Hit) {
					calls++
				}),
			),
		)
		require.Equal(t, 1, calls)
	})
}

func TestDescription(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	err := Do(
		Description("Test #1"),
		Custom(BeforeSendStep, func(hit Hit) {
			require.Equal(t, "Test #1", hit.Description())
			hit.SetDescription("Test #A")
			require.Equal(t, "Test #A", hit.Description())
		}),
		Post(s.URL),
		Send("Hello"),
		Expect("World"),
	)
	require.NotNil(t, err)
	require.True(t, strings.HasPrefix(vtclean.Clean(err.Error(), false), "Description:\tTest #A"))
}

func TestCustomError(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Custom(func(hit Hit) {
				panic("some error")
			}),
		),
		PtrStr("some error"),
	)
}

func TestDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Stop Execution", func(t *testing.T) {
		shouldNotRun := false
		ExpectError(t,
			Do(
				Post(s.URL),
				Send("Hello World"),
				Custom(ExpectStep, func(hit Hit) {
					hit.MustDo(
						Expect("Hello Universe"),
					)
					shouldNotRun = true
				}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
		)
		require.False(t, shouldNotRun)
	})

	t.Run("Expect in Send Step", func(t *testing.T) {
		shouldNotRun := false
		ExpectError(t,
			Do(
				Post(s.URL),
				Send("Hello World"),
				Custom(SendStep, func(hit Hit) {
					hit.MustDo(
						Expect("Hello Universe"),
					)
					shouldNotRun = true
				}),
			),
			PtrStr("unable to execute `Expect' during SendStep, can only be run during ExpectStep"),
		)
		require.False(t, shouldNotRun)
	})
}

func TestOutOfContext(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(
		t,
		Post(s.URL),
		Send("World"),
		Expect("World"),
		Send().Custom(func(hit Hit) {
			Send("Hello Universe") // this will never be run, because you need to wrap this with hit.Do()/MustDo()
		}),
		Expect().Custom(func(hit Hit) {
			Expect("Hello Universe") // this will never be run, because you need to wrap this with hit.Do()/MustDo()
		}),
	)
}
