package hit

import (
	"net/http"
	"os"

	"io"

	"fmt"

	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"github.com/Eun/go-hit/internal/minitest"
)

var ett *errortrace.ErrorTraceTemplate

func init() {
	ett = errortrace.New(
		"testing",
		"runtime",
		errortrace.IgnoreFunc((*hitStep).exec),
		errortrace.IgnoreStruct((*defaultInstance).Do),
		errortrace.IgnoreFunc((*expect).Custom),
		errortrace.IgnoreFunc((*send).Custom),
		errortrace.IgnoreFunc(Do),
		errortrace.IgnoreFunc(Test),
		errortrace.IgnorePackage(minitest.Errorf),
	)
}

// Send sends the specified data as the body payload
//
// Examples:
//     MustDo(
//         Post("https://example.com"),
//         Send("Hello World"),
//     )
//
//     MustDo(
//         Post("https://example.com"),
//         Send().Body("Hello World")
//     )
func Send(data ...interface{}) ISend {
	return newSend(newClearPath("Send", data), data)
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Expect().Body().Contains("Hello World")
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Expect("Hello World"),
//     )
func Expect(data ...interface{}) IExpect {
	return newExpect(newClearPath("Expect", data), data)
}

// Debug prints the current Request and Response to hit.Stdout(), you can filter the output based on expressions
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Debug(),
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Debug("Response.Headers"),
//     )
func Debug(expression ...string) IStep {
	return newDebug(expression)
}

// HTTPClient sets the client for the request
//
// Example:
//     var client http.Client
//     MustDo(
//         Get("https://example.com"),
//         HTTPClient(&client),
//     )
func HTTPClient(client *http.Client) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.SetHTTPClient(client)
			return nil
		},
	}
}

// Stdout sets the output to the specified writer
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         Stdout(os.Stderr),
//         Debug(),
//     )
func Stdout(w io.Writer) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.SetStdout(w)
			return nil
		},
	}
}

// BaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace or Method
//
// Examples:
//     MustDo(
//         BaseURL("https://example.com")
//     )
//
//     MustDo(
//         BaseURL("https://%s/%s", "example.com", "index.html")
//     )
func BaseURL(url string, a ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.SetBaseURL(url, a...)
			return nil
		},
	}
}

// Request creates a new Hit instance with an existing http request
//
// Example:
//     request, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
//     MustDo(
//         Request(request),
//     )
func Request(request *http.Request) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.SetRequest(request)
			return nil
		},
	}
}

// Method creates a new Hit instance with the specified method and url
//
// Examples:
//     MustDo(
//         Method(http.MethodGet, "https://example.com"),
//     )
//
//     MustDo(
//         Method(http.MethodGet, "https://%s/%s", "example.com", "index.html"),
//     )

func Method(method, url string, a ...interface{}) IStep {
	return Method(method, url, a)
}

func makeMethodStep(method, url string, a ...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			request, err := http.NewRequest(method, internal.MakeURL(hit.BaseURL(), url, a...), nil)
			if err != nil {
				return err
			}
			// remove some standard headers
			request.Header.Set("User-Agent", "")
			hit.SetRequest(request)
			return nil
		},
	}
}

// Connect creates a new Hit instance with CONNECT as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Connect("https://example.com"),
//     )
//
//     MustDo(
//         Connect("https://%s/%s", "example.com", "index.html"),
//     )
func Connect(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodConnect, url, a...)
}

// Delete creates a new Hit instance with DELETE as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Delete("https://example.com"),
//     )
//
//     MustDo(
//         Delete("https://%s/%s", "example.com", "index.html"),
//     )
//
func Delete(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodDelete, url, a...)
}

// Get creates a new Hit instance with GET as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//     )
//
//     MustDo(
//         Get("https://%s/%s", "example.com", "index.html"),
//     )
//
func Get(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodGet, url, a...)
}

// Head creates a new Hit instance with HEAD as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Head("https://example.com"),
//     )
//
//     MustDo(
//         Head("https://%s/%s", "example.com", "index.html"),
//     )
//
func Head(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodHead, url, a...)
}

// Post creates a new Hit instance with POST as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Post("https://example.com"),
//     )
//
//     MustDo(
//         Post("https://%s/%s", "example.com", "index.html"),
//     )
//
func Post(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodPost, url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Options("https://example.com"),
//     )
//
//     MustDo(
//         Options("https://%s/%s", "example.com", "index.html"),
//     )
//
func Options(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodOptions, url, a...)
}

// Put creates a new Hit instance with PUT as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Put("https://example.com"),
//     )
//
//     MustDo(
//         Put("https://%s/%s", "example.com", "index.html"),
//     )
//
func Put(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodPut, url, a...)
}

// Trace creates a new Hit instance with TRACE as the http makeMethodStep, use the optional arguments to format the url
//
// Examples:
//     MustDo(
//         Trace("https://example.com"),
//     )
//
//     MustDo(
//         Trace("https://%s/%s", "example.com", "index.html"),
//     )
//
func Trace(url string, a ...interface{}) IStep {
	return makeMethodStep(http.MethodTrace, url, a...)
}

// Test runs the specified steps and calls t.Error() if any error occurs during execution
func Test(t TestingT, steps ...IStep) {
	if err := Do(steps...); err != nil {
		if _, ok := err.(errortrace.ErrorTraceError); !ok {
			os.Stderr.WriteString(ett.Format("", err.Error()).Error())
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
		state:  CombineStep,
	}
	if err := hit.runSteps(CombineStep); err != nil {
		return err
	}
	hit.state = CleanStep
	if err := hit.runSteps(CleanStep); err != nil {
		return err
	}
	hit.state = BeforeSendStep
	if err := hit.runSteps(BeforeSendStep); err != nil {
		return err
	}
	if hit.request == nil {
		return fmt.Errorf("unable to perform request: no request set, did you called Post(), Get(), ...?")
	}
	hit.state = SendStep
	if err := hit.runSteps(SendStep); err != nil {
		return err
	}
	hit.state = AfterSendStep
	if err := hit.runSteps(AfterSendStep); err != nil {
		return err
	}
	hit.request.Request.Body = hit.request.Body().Reader()
	res, err := hit.client.Do(hit.request.Request)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s", err.Error())
	}

	hit.response = newHTTPResponse(hit, res)
	hit.state = BeforeExpectStep
	if err := hit.runSteps(BeforeExpectStep); err != nil {
		return err
	}
	hit.state = ExpectStep
	if err := hit.runSteps(ExpectStep); err != nil {
		return err
	}
	hit.state = AfterExpectStep
	if err := hit.runSteps(AfterExpectStep); err != nil {
		return err
	}
	if hit.request.Request.Body != nil {
		if err := hit.request.Request.Body.Close(); err != nil {
			return err
		}
	}
	return err
}

// MustDo runs the specified steps and panics with the error if something was wrong
func MustDo(steps ...IStep) {
	if err := Do(steps...); err != nil {
		panic(err)
	}
}

// CombineSteps combines multiple steps to one
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         CombineSteps(
//            Expect().Status(http.StatusOK),
//            Expect().Body("Hello World"),
//         ),
//     )
func CombineSteps(steps ...IStep) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CombineStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.AddSteps(steps...)
			return nil
		},
	}
}

// Description sets a custom description for this test.
// The description will be printed in an error case
func Description(description string) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeSendStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			hit.SetDescription(description)
			return nil
		},
	}
}

// Clear can be used to remove previous steps.
//
// Examples:
//     Clear().Send("Hello World")          // will remove all Send("Hello World") steps
//     Clear().Send().Body("Hello World")   // will remove all Send().Body("Hello World") steps
//     Clear().Expect().Body()              // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
//     Clear().Expect().Body("Hello World") // will remove all Expect().Body("Hello World") steps
//
//     MustDo(
//         Post("https://example.com"),
//         Send("Hello Earth"),
//         Expect("Hello World"),
//         Clear().Expect(),
//         Expect("Hello Earth"),
//     )
func Clear() IClear {
	return newClear()
}
