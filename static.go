package hit

import (
	"fmt"
	"net/http"
	"os"

	"io"

	"github.com/Eun/go-hit/errortrace"
)

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) ISend {
	if arg, ok := getLastArgument(data); ok {
		return &dummySend{newSend().Interface(arg)}
	}
	return newSend()
}

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) IExpect {
	if arg, ok := getLastArgument(data); ok {
		return &dummyExpect{newExpect().Interface(arg)}
	}
	return newExpect()
}

// Debug prints the current Request and Response to hit.Stdout()
func Debug() IStep {
	return Custom(BeforeExpectStep, func(hit Hit) {
		hit.Debug()
	})
}

// SetHTTPClient sets the client for the request
func SetHTTPClient(client *http.Client) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetHTTPClient(client)
	})
}

// SetStdout sets the output to the specified writer
func SetStdout(w io.Writer) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetStdout(w)
	})
}

// SetBaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, SetMethod call
func SetBaseURL(s string) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetBaseURL(s)
	})
}

// SetRequest creates a new Hit instance with an existing http request
func SetRequest(request *http.Request) IStep {
	return Custom(BeforeSendStep, func(hit Hit) {
		hit.SetRequest(request)
	})
}

// SetMethod creates a new Hit instance with the specified method and url
// SetMethod(t, "POST", "http://%s/%s", domain, path)
func SetMethod(method, url string, a ...interface{}) IStep {
	et := errortrace.Prepare()
	return Custom(BeforeSendStep, func(hit Hit) {
		request, err := http.NewRequest(method, makeURL(hit.BaseURL(), url, a...), nil)
		et.NoError(err, "unable to create request")
		hit.SetRequest(request)
	})
}

// Connect creates a new Hit instance with CONNECT as the http method, use the optional arguments to format the url
// Connect(t, "http://%s/%s", domain, path)
func Connect(url string, a ...interface{}) IStep {
	return SetMethod("CONNECT", url, a...)
}

// Delete creates a new Hit instance with DELETE as the http method, use the optional arguments to format the url
// Delete("http://%s/%s", domain, path)
func Delete(url string, a ...interface{}) IStep {
	return SetMethod("DELETE", url, a...)
}

// Get creates a new Hit instance with GET as the http method, use the optional arguments to format the url
// Get("http://%s/%s", domain, path)
func Get(url string, a ...interface{}) IStep {
	return SetMethod("GET", url, a...)
}

// Head creates a new Hit instance with HEAD as the http method, use the optional arguments to format the url
// Head("http://%s/%s", domain, path)
func Head(url string, a ...interface{}) IStep {
	return SetMethod("HEAD", url, a...)
}

// Post creates a new Hit instance with POST as the http method, use the optional arguments to format the url
// Post("http://%s/%s", domain, path)
func Post(url string, a ...interface{}) IStep {
	return SetMethod("POST", url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http method, use the optional arguments to format the url
// Options("http://%s/%s", domain, path)
func Options(url string, a ...interface{}) IStep {
	return SetMethod("OPTIONS", url, a...)
}

// Put creates a new Hit instance with PUT as the http method, use the optional arguments to format the url
// Put("http://%s/%s", domain, path)
func Put(url string, a ...interface{}) IStep {
	return SetMethod("PUT", url, a...)
}

// Trace creates a new Hit instance with TRACE as the http method, use the optional arguments to format the url
// Trace("http://%s/%s", domain, path)
func Trace(url string, a ...interface{}) IStep {
	return SetMethod("TRACE", url, a...)
}

// Test runs the specified steps and calls t.Error() if any error occurs during execution
func Test(t TestingT, steps ...IStep) {
	if err := Do(steps...); err != nil {
		t.Error(err.Error())
	}
}

// Do runs the specified steps and returns error if something was wrong
func Do(steps ...IStep) error {
	hit := &defaultInstance{
		client: http.DefaultClient,
		stdout: os.Stdout,
		steps:  steps,
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
