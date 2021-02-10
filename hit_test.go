package hit_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
	"golang.org/x/xerrors"

	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

//nolint:interfacer //this signature is required
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRequest(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)
		Test(t,
			Request().Set(req),
			Send().Body().String("Hello Universe"),
			Expect().Body().String().Equal("Hello Universe"),
		)
	})

	t.Run("overwrite request during send", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			Custom(BeforeSendStep, func(hit Hit) error {
				r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.URL, nil)
				if err != nil {
					return err
				}
				if err := hit.SetRequest(r); err != nil {
					return err
				}
				return nil
			}),
			Send().Body().String("Hello Universe"),
			Expect().Body().String().Equal("Hello Universe"),
		)
	})

	t.Run("context with timeout", func(t *testing.T) {
		closeChan := make(chan struct{})
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			tmr := time.NewTimer(time.Minute)
			defer tmr.Stop()
			select {
			case <-tmr.C:
				writer.WriteHeader(http.StatusOK)
			case <-closeChan:
			}
		})
		s := httptest.NewServer(mux)
		defer s.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)
		ExpectError(t,
			Do(
				Request().Set(req),
				Send().Body().String("Hello Universe"),
				Expect().Body().String().Equal("Hello Universe"),
			),
			PtrStr(fmt.Sprintf(`unable to perform request: Post "%s": context deadline exceeded`, s.URL)),
		)
		closeChan <- struct{}{}
	})

	t.Run("context with timeout using regular POST method", func(t *testing.T) {
		closeChan := make(chan struct{})
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			tmr := time.NewTimer(time.Minute)
			defer tmr.Stop()
			select {
			case <-tmr.C:
				writer.WriteHeader(http.StatusOK)
			case <-closeChan:
			}
		})
		s := httptest.NewServer(mux)
		defer s.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		ExpectError(t,
			Do(
				Post(s.URL),
				Context(ctx),
				Send().Body().String("Hello Universe"),
				Expect().Body().String().Equal("Hello Universe"),
			),
			PtrStr(fmt.Sprintf(`unable to perform request: Post "%s": context deadline exceeded`, s.URL)),
		)
		closeChan <- struct{}{}
	})
}

func TestMultiUse(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("body", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)

		Test(t,
			Request().Set(req),
			Expect().Body().String().Equal("Hello World"),
		)

		Test(t,
			Request().Set(req),
			Expect().Body().String().Equal("Hello World"),
		)
	})
	t.Run("header/trailer", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.URL, bytes.NewReader([]byte("Hello World")))
		require.NoError(t, err)
		req.Header.Set("X-Header", "Foo")
		req.Trailer = map[string][]string{"X-Trailer": {"Bar"}}

		Test(t,
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"Foo"}, hit.Request().Header["X-Header"])
				require.Equal(t, []string{"Bar"}, hit.Request().Trailer["X-Trailer"])
				return nil
			}),
		)

		Test(t,
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"Foo"}, hit.Request().Header["X-Header"])
				require.Equal(t, []string{"Bar"}, hit.Request().Trailer["X-Trailer"])
				return nil
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
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().Form["a"])
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().PostForm["a"])
				return nil
			}),
		)

		Test(t,
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().Form["a"])
				require.Equal(t, []string{"1", "2", "banana"}, hit.Request().PostForm["a"])
				return nil
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
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"bar"}, hit.Request().MultipartForm.Value["foo"])
				return nil
			}),
		)

		require.Len(t, req.MultipartForm.File["file1"], 1)

		Test(t,
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, "file1", hit.Request().MultipartForm.File["file1"][0].Filename)
				return nil
			}),
		)
		f, err := req.MultipartForm.File["file1"][0].Open()
		require.NoError(t, err)
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, "baz", string(buf))

		Test(t,
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Len(t, hit.Request().MultipartForm.File["file1"], 1)
				require.Equal(t, "file1", hit.Request().MultipartForm.File["file1"][0].Filename)
				f, err = hit.Request().MultipartForm.File["file1"][0].Open()
				require.NoError(t, err)
				defer f.Close()
				buf, err = ioutil.ReadAll(f)
				require.NoError(t, err)
				require.Equal(t, "baz", string(buf))
				return nil
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
			Request().Set(req),
			Custom(AfterSendStep, func(hit Hit) error {
				require.Equal(t, []string{"a", "b"}, hit.Request().TransferEncoding)
				return nil
			}),
		)
	})

	t.Run("with embedded request", func(t *testing.T) {
		s := EchoServer()
		defer s.Close()
		template := []IStep{
			Post(s.URL),
			Send().Headers("Content-Type").Add("application/json"),
			Expect().Headers("Content-Type").Equal("application/json"),
		}
		Test(t,
			append(template,
				Send().Body().JSON("Hello World"),
				Expect().Body().JSON().Equal("Hello World"),
			)...,
		)

		Test(t,
			append(template,
				Send().Body().JSON("Hello Universe"),
				Expect().Body().JSON().Equal("Hello Universe"),
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
		Expect().Status().Equal(http.StatusNoContent),
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
		Expect().Status().Equal(http.StatusNoContent),
	)
}

func TestHTTPClient(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	client := &http.Client{}

	Test(t,
		Post(s.URL),
		HTTPClient(client),
		Custom(AfterSendStep, func(hit Hit) error {
			require.Equal(t, client, hit.HTTPClient())
			return nil
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
				Send().Body().String("Hello"),
				Expect().Body().String().Equal("Hello"),
			),
			Expect().Body().String().Equal("World"),
		),
		PtrStr("not equal"), PtrStr(`expected: "World"`), nil, nil, nil, nil, nil,
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
				Send().Custom(func(hit Hit) error {
					calls++
					return nil
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
				Expect().Custom(func(hit Hit) error {
					calls++
					return nil
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
				Custom(BeforeSendStep, func(hit Hit) error {
					calls++
					return nil
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
		Custom(BeforeSendStep, func(hit Hit) error {
			require.Equal(t, "Test #1", hit.Description())
			hit.SetDescription("Test #A")
			require.Equal(t, "Test #A", hit.Description())
			return nil
		}),
		Post(s.URL),
		Send().Body().String("Hello"),
		Expect().Body().String().Equal("World"),
	)
	require.NotNil(t, err)
	require.True(t, strings.HasPrefix(vtclean.Clean(err.Error(), false), "Description:\tTest #A"))
}

func TestCustomError(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("panic", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Custom(func(hit Hit) error {
					panic("some error")
				}),
			),
			PtrStr("some error"),
		)
	})

	t.Run("error", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Custom(func(hit Hit) error {
					return errors.New("some error")
				}),
			),
			PtrStr("some error"),
		)
	})

	t.Run("errors is", func(t *testing.T) {
		myError := errors.New("some error")
		err := Do(
			Post(s.URL),
			Send().Custom(func(hit Hit) error {
				return myError
			}),
		)
		require.True(t, xerrors.Is(err, myError))
	})

	t.Run("errors as", func(t *testing.T) {
		err := Do(
			Post(s.URL),
			Send().Custom(func(hit Hit) error {
				return &os.SyscallError{
					Syscall: "test",
					Err:     xerrors.New("go away"),
				}
			}),
		)
		var sysCallErr *os.SyscallError
		require.True(t, xerrors.As(err, &sysCallErr))
		require.Equal(t, sysCallErr.Syscall, "test")
		require.Equal(t, sysCallErr.Err.Error(), "go away")
	})
}

func TestDo(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Stop Execution", func(t *testing.T) {
		shouldNotRun := false
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("Hello World"),
				Custom(ExpectStep, func(hit Hit) error {
					hit.MustDo(
						Expect().Body().String().Equal("Hello Universe"),
					)
					shouldNotRun = true
					return nil
				}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
		)
		require.False(t, shouldNotRun)
	})

	t.Run("Expect in Send Step", func(t *testing.T) {
		shouldNotRun := false
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("Hello World"),
				Custom(SendStep, func(hit Hit) error {
					hit.MustDo(
						Expect().Body().String().Equal("Hello Universe"),
					)
					shouldNotRun = true
					return nil
				}),
			),
			PtrStr(`unable to execute Expect().Body().String().Equal("Hello Universe") during SendStep, can only be run during ExpectStep`),
		)
		require.False(t, shouldNotRun)
	})
}

func TestOutOfContext(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().String("World"),
		Expect().Body().String().Equal("World"),
		Send().Custom(func(hit Hit) error {
			Send().Body().String("Hello Universe") // this will never be run, because you need to wrap this with hit.Do()/MustDo()
			return nil
		}),
		Expect().Custom(func(hit Hit) error {
			Expect().Body().String().Equal("Hello Universe") // this will never be run, because you need to wrap this with hit.Do()/MustDo()
			return nil
		}),
	)
}

func TestAddSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	var callOrder []int
	Test(t,
		Post(s.URL),
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 1)
			hit.AddSteps(
				Custom(BeforeSendStep, func(hit Hit) error {
					callOrder = append(callOrder, 2)
					return nil
				}),
			)
			return nil
		}),
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 3)
			return nil
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
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 1)
			hit.InsertSteps(
				Custom(BeforeSendStep, func(hit Hit) error {
					callOrder = append(callOrder, 2)
					return nil
				}),
			)
			return nil
		}),
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 3)
			return nil
		}),
	)

	require.Equal(t, []int{1, 2, 3}, callOrder)
}

func TestAddAndRemoveSteps(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	var callOrder []int
	someStep := Custom(BeforeSendStep, func(hit Hit) error { return nil })
	Test(t,
		Post(s.URL),
		someStep,
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 1)
			hit.AddSteps(
				Custom(BeforeSendStep, func(hit Hit) error {
					callOrder = append(callOrder, 2)
					return nil
				}),
			)
			hit.RemoveSteps(someStep)
			return nil
		}),
		Custom(BeforeSendStep, func(hit Hit) error {
			callOrder = append(callOrder, 3)
			return nil
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
			Send().Body().String("Hello Alice"),
			Expect().Body().String().Equal("Hello Joe"),
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

	Test(t, Method(http.MethodGet, s.URL), Expect().Headers("Method").Equal(http.MethodGet))
	Test(t, Delete(s.URL), Expect().Headers("Method").Equal(http.MethodDelete))
	Test(t, Get(s.URL), Expect().Headers("Method").Equal(http.MethodGet))
	Test(t, Head(s.URL), Expect().Headers("Method").Equal(http.MethodHead))
	Test(t, Post(s.URL), Expect().Headers("Method").Equal(http.MethodPost))
	Test(t, Options(s.URL), Expect().Headers("Method").Equal(http.MethodOptions))
	Test(t, Put(s.URL), Expect().Headers("Method").Equal(http.MethodPut))
	Test(t, Trace(s.URL), Expect().Headers("Method").Equal(http.MethodTrace))
}

func TestReturn(t *testing.T) {
	executed := false
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Custom(func(hit Hit) error {
			executed = true
			return nil
		}),
		Expect().Body().String().Equal("Hello World"),
		Return(),
		Expect().Body().String().Equal("Hello Universe"),
	)
	require.True(t, executed)
}

func BenchmarkBigPayloads(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			io.WriteString(writer, err.Error())
			return
		}
		size, err := strconv.ParseInt(request.PostForm.Get("length"), 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			io.WriteString(writer, err.Error())
			return
		}
		writer.WriteHeader(http.StatusOK)
		io.Copy(writer, io.LimitReader(rand.Reader, size))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	steps := CombineSteps(
		Post(s.URL),
		Send().Headers("Content-Type").Add("application/x-www-form-urlencoded"),
		Expect().Status().Equal(http.StatusOK),
	)

	payloadSizes := []int{
		1024,
		1024 * 10,
		1024 * 100,
		1024 * 1000,
	}
	for _, size := range payloadSizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			Test(b,
				steps,
				Send().Body().FormValues("length").Add(size),
			)
		})
	}
}

func TestMissingRequest(t *testing.T) {
	ExpectError(t,
		Do(
			Send().Body().String("Hello World"),
		),
		PtrStr("unable to create a request: did you called Post(), Get(), ...?"),
	)
}

func TestNoUserAgent(t *testing.T) {
	// test if the default is no user agent set

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Empty(t, request.UserAgent())
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t, Get(s.URL))
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		Params   []string
		Expected string
	}{
		{
			[]string{
				"http://example.com",
				"index.html",
			},
			"http://example.com/index.html",
		},
		{
			[]string{
				"http://example.com/",
				"index.html",
			},
			"http://example.com/index.html",
		},
		{
			[]string{
				"http://example.com",
				"/index.html",
			},
			"http://example.com/index.html",
		},
		{
			[]string{
				"http://example.com",
				"",
			},
			"http://example.com",
		},
		{
			[]string{
				"http://example.com/",
				"",
			},
			"http://example.com",
		},
		{
			[]string{
				"",
				"index.html",
			},
			"index.html",
		},
		{
			[]string{
				"",
				"/index.html",
			},
			"/index.html",
		},
		{
			[]string{
				"/index.html",
			},
			"/index.html",
		},
		{
			[]string{
				"http://",
				"example.com",
				"index.html",
			},
			"http://example.com/index.html",
		},
		{
			[]string{
				"http://",
				"/example.com/",
				"/index.html/",
			},
			"http://example.com/index.html",
		},
		{
			[]string{
				"example.com",
				"index.html",
			},
			"example.com/index.html",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run("", func(t *testing.T) {
			if s := JoinURL(test.Params...); s != test.Expected {
				t.Errorf("expected `%s' got `%s'", test.Expected, s)
			}
		})
	}
}

func TestNilStep(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t, nil, Get(s.URL), nil)
}
