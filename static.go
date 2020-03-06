package hit

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"io"

	"sync"

	"runtime"

	"strings"

	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

var contexts sync.Map

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) ISend {
	return newSend(NewCleanPath("Send", data), data)
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) IExpect {
	return newExpect(NewCleanPath("Expect", data), data)
}

// Debug prints the current Request and Response to hit.Stdout(), you can filter the output based on expressions
// Examples:
//           Debug()
//           Debug("Request")
func Debug(expression ...string) IStep {
	return newDebug(expression)
}

// HTTPClient sets the client for the request
func HTTPClient(client *http.Client) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.SetHTTPClient(client)
		},
	})
}

// Stdout sets the output to the specified writer
func Stdout(w io.Writer) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.SetStdout(w)
		},
	})
}

// BaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, SetMethod makeMethodStep
// Examples:
//           BaseURL("http://example.com")
//           BaseURL("http://%s/%s", domain, path)
func BaseURL(url string, a ...interface{}) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.SetBaseURL(url, a...)
		},
	})
}

// Request creates a new Hit instance with an existing http request
func Request(request *http.Request) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.SetRequest(request)
		},
	})
}

// Method creates a new Hit instance with the specified makeMethodStep and url
// Examples:
//           Method("POST", "http://example.com")
//           Method("POST", "http://%s/%s", domain, path)
func Method(method, url string, a ...interface{}) IStep {
	return Method(method, url, a)
}

func makeMethodStep(method, url string, a ...interface{}) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			request, err := http.NewRequest(method, internal.MakeURL(hit.BaseURL(), url, a...), nil)
			minitest.NoError(err, "unable to create request")
			hit.SetRequest(request)
		},
	})
}

// Connect creates a new Hit instance with CONNECT as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Connect("http://example.com")
//           Connect("http://%s/%s", domain, path)
func Connect(url string, a ...interface{}) IStep {
	return makeMethodStep("CONNECT", url, a...)
}

// Delete creates a new Hit instance with DELETE as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Delete("http://example.com")
//           Delete("http://%s/%s", domain, path)
func Delete(url string, a ...interface{}) IStep {
	return makeMethodStep("DELETE", url, a...)
}

// Get creates a new Hit instance with GET as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Get("http://example.com")
//           Get("http://%s/%s", domain, path)
func Get(url string, a ...interface{}) IStep {
	return makeMethodStep("GET", url, a...)
}

// Head creates a new Hit instance with HEAD as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Head("http://example.com")
//           Head("http://%s/%s", domain, path)
func Head(url string, a ...interface{}) IStep {
	return makeMethodStep("HEAD", url, a...)
}

// Post creates a new Hit instance with POST as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Post("http://example.com")
//           Post("http://%s/%s", domain, path)
func Post(url string, a ...interface{}) IStep {
	return makeMethodStep("POST", url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Options("http://example.com")
//           Options("http://%s/%s", domain, path)
func Options(url string, a ...interface{}) IStep {
	return makeMethodStep("OPTIONS", url, a...)
}

// Put creates a new Hit instance with PUT as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Put("http://example.com")
//           Put("http://%s/%s", domain, path)
func Put(url string, a ...interface{}) IStep {
	return makeMethodStep("PUT", url, a...)
}

// Trace creates a new Hit instance with TRACE as the http makeMethodStep, use the optional arguments to format the url
// Examples:
//           Trace("http://example.com")
//           Trace("http://%s/%s", domain, path)
func Trace(url string, a ...interface{}) IStep {
	return makeMethodStep("TRACE", url, a...)
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

	if err := hit.runSteps(CleanStep); err != nil {
		return err
	}

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

func MustDo(steps ...IStep) {
	if err := Do(steps...); err != nil {
		panic(err)
	}
}

// CombineSteps combines multiple steps to one
func CombineSteps(steps ...IStep) IStep {
	return custom(Step{
		When:      BeforeSendStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.AddSteps(steps...)
		},
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
	return custom(Step{
		When:      BeforeExpectStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			hit.SetDescription(description)
		},
	})
}

func Clear() IClear {
	return newClear(getContext())
}
