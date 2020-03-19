package hit_test

import (
	"net/http"
	"testing"

	"strings"

	"bytes"

	"io/ioutil"
	"net/url"

	"net/http/httptest"

	. "github.com/Eun/go-hit"
	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

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
	mux.HandleFunc("/foo/", func(writer http.ResponseWriter, request *http.Request) {
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
	Test(t,
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

func TestAddSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	var callOrder []int
	Test(t,
		Post(s.URL),
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 1)
			hit.AddSteps(
				Custom(BeforeSendStep, func(hit Hit) {
					callOrder = append(callOrder, 2)
				}),
			)
		}),
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 3)
		}),
	)

	require.Equal(t, []int{1, 3, 2}, callOrder)
}

func TestInsertSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	var callOrder []int
	Test(t,
		Post(s.URL),
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 1)
			hit.InsertSteps(
				Custom(BeforeSendStep, func(hit Hit) {
					callOrder = append(callOrder, 2)
				}),
			)
		}),
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 3)
		}),
	)

	require.Equal(t, []int{1, 2, 3}, callOrder)
}

func TestAddAndRemoveSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	var callOrder []int
	someStep := Custom(BeforeSendStep, func(hit Hit) {})
	Test(t,
		Post(s.URL),
		someStep,
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 1)
			hit.AddSteps(
				Custom(BeforeSendStep, func(hit Hit) {
					callOrder = append(callOrder, 2)
				}),
			)
			hit.RemoveSteps(someStep)
		}),
		Custom(BeforeSendStep, func(hit Hit) {
			callOrder = append(callOrder, 3)
		}),
	)

	require.Equal(t, []int{1, 3, 2}, callOrder)
}

func TestMustDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	require.Panics(t, func() {
		MustDo(
			Post(s.URL),
			Send("Hello Alice"),
			Expect("Hello Joe"),
		)
	})
}

func TestMethods(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Method", request.Method)
		writer.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t, Method(http.MethodGet, s.URL), Expect().Header("Method").Equal(http.MethodGet))
	// Test(t, Connect(s.URL), Expect().Header("Method").Equal(http.MethodConnect))
	Test(t, Delete(s.URL), Expect().Header("Method").Equal(http.MethodDelete))
	Test(t, Get(s.URL), Expect().Header("Method").Equal(http.MethodGet))
	Test(t, Head(s.URL), Expect().Header("Method").Equal(http.MethodHead))
	Test(t, Post(s.URL), Expect().Header("Method").Equal(http.MethodPost))
	Test(t, Options(s.URL), Expect().Header("Method").Equal(http.MethodOptions))
	Test(t, Put(s.URL), Expect().Header("Method").Equal(http.MethodPut))
	Test(t, Trace(s.URL), Expect().Header("Method").Equal(http.MethodTrace))
}
