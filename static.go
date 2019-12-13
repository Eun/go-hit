package hit

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"io"

	"sync"

	"runtime"

	"encoding/json"
	"time"

	"strings"

	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal/minitest"
	"github.com/gookit/color"
	"github.com/tidwall/pretty"
)

var contexts sync.Map

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) ISend {
	hit := getContext()
	if arg, ok := getLastArgument(data); ok {
		return finalSend{newSend(hit).Interface(arg)}
	}
	return newSend(hit)
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) IExpect {
	hit := getContext()
	if arg, ok := getLastArgument(data); ok {
		return finalExpect{newExpect(hit).Interface(arg)}
	}
	return newExpect(hit)
}

// Debug prints the current Request and Response to hit.Stdout(), you can filter the output based on expressions
// Examples:
//           Debug()
//           Debug("Request")
func Debug(expression ...string) IStep {
	debug := func(hit Hit) {
		type M map[string]interface{}

		getBody := func(body *HTTPBody) interface{} {
			reader := body.JSON().body.Reader()
			// if there is a json reader
			if reader != nil {
				var container interface{}
				if err := json.NewDecoder(reader).Decode(&container); err == nil {
					return container
				}
			}
			s := body.String()
			if len(s) == 0 {
				return nil
			}
			return s
		}

		getHeader := func(header http.Header) map[string]interface{} {
			m := make(map[string]interface{})
			for key := range header {
				m[key] = header.Get(key)
			}
			return m
		}

		m := M{
			"Time": time.Now().String(),
		}

		if hit.Request() != nil {
			m["Request"] = M{
				"Header":           getHeader(hit.Request().Header),
				"Trailer":          getHeader(hit.Request().Trailer),
				"Method":           hit.Request().Method,
				"URL":              hit.Request().URL,
				"Proto":            hit.Request().Proto,
				"ProtoMajor":       hit.Request().ProtoMajor,
				"ProtoMinor":       hit.Request().ProtoMinor,
				"ContentLength":    hit.Request().ContentLength,
				"TransferEncoding": hit.Request().TransferEncoding,
				"Host":             hit.Request().Host,
				"Form":             hit.Request().Form,
				"PostForm":         hit.Request().PostForm,
				"MultipartForm":    hit.Request().MultipartForm,
				"RemoteAddr":       hit.Request().RemoteAddr,
				"RequestURI":       hit.Request().RequestURI,
				"Body":             getBody(hit.Request().Body()),
			}
		}

		if hit.Response() != nil {
			m["Response"] = M{
				"Header":           getHeader(hit.Response().Header),
				"Trailer":          getHeader(hit.Response().Trailer),
				"Proto":            hit.Response().Proto,
				"ProtoMajor":       hit.Response().ProtoMajor,
				"ProtoMinor":       hit.Response().ProtoMinor,
				"ContentLength":    hit.Response().ContentLength,
				"TransferEncoding": hit.Response().TransferEncoding,
				"Body":             getBody(hit.Response().body),
				"Status":           hit.Response().Status,
				"StatusCode":       hit.Response().StatusCode,
			}
		}

		var v interface{} = m
		if size := len(expression); size > 0 {
			v = expr.MustGetValue(m, expression[size-1], expr.IgnoreError, expr.IgnoreNotFound)
		}

		bytes, err := json.Marshal(v)
		minitest.NoError(err)
		bytes = pretty.Pretty(bytes)
		if color.IsSupportColor() {
			bytes = pretty.Color(bytes, nil)
		}
		_, _ = hit.Stdout().Write(bytes)
	}

	hit := getContext()
	if hit != nil {
		debug(hit)
		return nil
	}
	return Custom(BeforeExpectStep, debug)
}

// HTTPClient sets the client for the request
func HTTPClient(client *http.Client) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetHTTPClient(client)
	})
}

// Stdout sets the output to the specified writer
func Stdout(w io.Writer) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetStdout(w)
	})
}

// BaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, SetMethod method
// Examples:
//           BaseURL("http://example.com")
//           BaseURL("http://%s/%s", domain, path)
func BaseURL(url string, a ...interface{}) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetBaseURL(url, a...)
	})
}

// Request creates a new Hit instance with an existing http request
func Request(request *http.Request) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetRequest(request)
	})
}

// Method creates a new Hit instance with the specified method and url
// Examples:
//           Method("POST", "http://example.com")
//           Method("POST", "http://%s/%s", domain, path)
func Method(method, url string, a ...interface{}) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		request, err := http.NewRequest(method, makeURL(hit.BaseURL(), url, a...), nil)
		minitest.NoError(err, "unable to create request")
		hit.SetRequest(request)
	})
}

// Connect creates a new Hit instance with CONNECT as the http method, use the optional arguments to format the url
// Examples:
//           Connect("http://example.com")
//           Connect("http://%s/%s", domain, path)
func Connect(url string, a ...interface{}) IStep {
	return Method("CONNECT", url, a...)
}

// Delete creates a new Hit instance with DELETE as the http method, use the optional arguments to format the url
// Examples:
//           Delete("http://example.com")
//           Delete("http://%s/%s", domain, path)
func Delete(url string, a ...interface{}) IStep {
	return Method("DELETE", url, a...)
}

// Get creates a new Hit instance with GET as the http method, use the optional arguments to format the url
// Examples:
//           Get("http://example.com")
//           Get("http://%s/%s", domain, path)
func Get(url string, a ...interface{}) IStep {
	return Method("GET", url, a...)
}

// Head creates a new Hit instance with HEAD as the http method, use the optional arguments to format the url
// Examples:
//           Head("http://example.com")
//           Head("http://%s/%s", domain, path)
func Head(url string, a ...interface{}) IStep {
	return Method("HEAD", url, a...)
}

// Post creates a new Hit instance with POST as the http method, use the optional arguments to format the url
// Examples:
//           Post("http://example.com")
//           Post("http://%s/%s", domain, path)
func Post(url string, a ...interface{}) IStep {
	return Method("POST", url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http method, use the optional arguments to format the url
// Examples:
//           Options("http://example.com")
//           Options("http://%s/%s", domain, path)
func Options(url string, a ...interface{}) IStep {
	return Method("OPTIONS", url, a...)
}

// Put creates a new Hit instance with PUT as the http method, use the optional arguments to format the url
// Examples:
//           Put("http://example.com")
//           Put("http://%s/%s", domain, path)
func Put(url string, a ...interface{}) IStep {
	return Method("PUT", url, a...)
}

// Trace creates a new Hit instance with TRACE as the http method, use the optional arguments to format the url
// Examples:
//           Trace("http://example.com")
//           Trace("http://%s/%s", domain, path)
func Trace(url string, a ...interface{}) IStep {
	return Method("TRACE", url, a...)
}

// Test runs the specified steps and calls t.Error() if any error occurs during execution
func Test(t TestingT, steps ...IStep) {
	if err := Do(steps...); err != nil {
		if _, ok := err.(errortrace.ErrorTraceError); !ok {
			os.Stderr.WriteString(errortrace.Format("", err.Error()).Error())
			t.FailNow()
		}
		os.Stderr.WriteString(err.Error())
		t.FailNow()
	}
}

// Do runs the specified steps and returns error if something was wrong
func Do(steps ...IStep) error {
	hit := &defaultInstance{
		client: http.DefaultClient,
		stdout: os.Stdout,
		steps:  steps,
	}
	id := getContextID()
	if id == 0 {
		return errors.New("unable to get context id")
	}
	contexts.Store(id, hit)
	defer contexts.Delete(id)

	if err := hit.runSteps(BeforeSendStep); err != nil {
		return err
	}
	if hit.request == nil {
		return fmt.Errorf("unable to perform request: no request set, did you called Post(), Get(), ...?")
	}
	if err := hit.runSteps(SendStep); err != nil {
		return err
	}
	if err := hit.runSteps(AfterSendStep); err != nil {
		return err
	}
	hit.request.Request.Body = hit.request.Body().Reader()
	res, err := hit.client.Do(hit.request.Request)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s", err.Error())
	}

	hit.response = newHTTPResponse(hit, res)
	if err := hit.runSteps(BeforeExpectStep); err != nil {
		return err
	}
	if err := hit.runSteps(ExpectStep); err != nil {
		return err
	}
	if err := hit.runSteps(AfterExpectStep); err != nil {
		return err
	}
	if hit.request.Request.Body != nil {
		_ = hit.request.Request.Body.Close()
	}
	return nil
}

// CombineSteps combines multiple steps to one
func CombineSteps(steps ...IStep) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.AddSteps(steps...)
	})
}

func getContextID() uintptr {
	var pc [16]uintptr

	skip := 0
	for {
		n := runtime.Callers(skip, pc[:])
		frames := runtime.CallersFrames(pc[:n])
		for {
			frame, more := frames.Next()
			if strings.HasSuffix(frame.Function, "github.com/Eun/go-hit.Do") {
				return frame.Entry
			}
			if !more {
				return 0
			}
		}
		skip += n
	}
}
func getContext() Hit {
	id := getContextID()
	if id == 0 {
		return nil
	}
	v, ok := contexts.Load(id)
	if !ok {
		return nil
	}
	return v.(Hit)
}

// Description sets a custom description for this test.
// The description will be printed in an error case
func Description(description string) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetDescription(description)
	})
}
