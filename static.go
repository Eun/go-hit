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
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) ISend {
	return newSend(newClearPath("Send", data), data)
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) IExpect {
	return newExpect(newClearPath("Expect", data), data)
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

// BaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, SetMethod makeMethodStep
// Examples:
//           BaseURL("http://example.com")
//           BaseURL("http://%s/%s", domain, path)
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

// Method creates a new Hit instance with the specified makeMethodStep and url
// Examples:
//           Method("POST", "http://example.com")
//           Method("POST", "http://%s/%s", domain, path)
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

func MustDo(steps ...IStep) {
	if err := Do(steps...); err != nil {
		panic(err)
	}
}

// CombineSteps combines multiple steps to one
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

func Clear() IClear {
	return newClear()
}
