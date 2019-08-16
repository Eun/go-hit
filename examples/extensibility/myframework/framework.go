package myframework

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Eun/go-hit"
)

type User struct {
	ID   int
	Name string
}

// Add a complete new function
func CheckTheSparkles() hit.IStep {
	fmt.Println("Lets see if we can find sparkles ✨✨✨...")
	return Expect().Custom(func(hit hit.Hit) {
		// use your custom code
		if !strings.Contains(hit.Response().Body().String(), "✨") {
			panic("No sparkles present, too sad :-(")
		}
		// or use our expect
		// hit.Expect().Body().Contains("✨")
	})
}

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) MySend {
	return MySend{
		ISend: hit.Send(data...),
	}
}

type MySend struct {
	hit.ISend
}

// User sets the body to the user
func (snd MySend) User(user User) hit.IStep {
	return hit.Custom(hit.SendStep, func(hit hit.Hit) {
		hit.Send().Headers().Set("Content-Type", "application/json")
		hit.Send().Body().JSON(user)
		// or hit.Request()
		// hit.Request().Header.Set("Content-Type", "application/json")
		// hit.Request().Body().JSON().Set(user)
	})
}

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) MyExpect {
	return MyExpect{
		IExpect: hit.Expect(data...),
	}
}

type MyExpect struct {
	hit.IExpect
}

// User expects the body to be equal to the specific user
func (exp MyExpect) User(user User) hit.IStep {
	// We have to return our Framework to make sure we can use our custom functions later on
	return hit.Custom(hit.ExpectStep, func(hit hit.Hit) {
		// use hit.IExpect
		hit.Expect().Header("Content-Type").Equal("application/json")
		hit.Expect().Body().JSON().Equal("", user)

		// or hit.Response()
		// if hit.Response().Header.Get("Content-Type") != "application/json" {
		// 	panic(fmt.Sprintf("%#v != %#v", hit.Response().Header.Get("Content-Type"), "application/json"))
		// }
		// if !reflect.DeepEqual(hit.Response().Body().JSON().GetAs(&User{}), &user) {
		// 	panic(fmt.Sprintf("%#v != %#v", hit.Response().Body().JSON().GetAs(&User{}), &user))
		// }
	})
}

// Implement the rest for convenience reasons
// you can copy the contents from the template_framework.go file

// Custom calls a custom Step on the specified execution time
func Custom(when hit.StepTime, exec hit.Callback) hit.IStep {
	return hit.Custom(when, exec)
}

// Debug prints the current Request and Response to hit.Stdout(), you can filter the output based on expressions
// Examples:
//           Debug()
//           Debug("Request")
func Debug(expression ...string) hit.IStep {
	return hit.Debug(expression...)
}

// SetHTTPClient sets the client for the request
func SetHTTPClient(client *http.Client) hit.IStep {
	return hit.SetHTTPClient(client)
}

// SetStdout sets the output to the specified writer
func SetStdout(w io.Writer) hit.IStep {
	return hit.SetStdout(w)
}

// SetBaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, SetMethod method
// Examples:
//           SetBaseURL("http://example.com")
//           SetBaseURL("http://%s/%s", domain, path)
func SetBaseURL(url string, a ...interface{}) hit.IStep {
	return hit.SetBaseURL(url, a...)
}

// SetRequest creates a new Hit instance with an existing http request
func SetRequest(request *http.Request) hit.IStep {
	return hit.SetRequest(request)
}

// SetMethod creates a new Hit instance with the specified method and url
// Examples:
//           SetMethod("POST", "http://example.com")
//           SetMethod("POST", "http://%s/%s", domain, path)
func SetMethod(method, url string, a ...interface{}) hit.IStep {
	return hit.SetMethod(method, url, a...)
}

// Connect creates a new Hit instance with CONNECT as the http method, use the optional arguments to format the url
// Examples:
//           Connect("http://example.com")
//           Connect("http://%s/%s", domain, path)
func Connect(url string, a ...interface{}) hit.IStep {
	return hit.Connect(url, a...)
}

// Delete creates a new Hit instance with DELETE as the http method, use the optional arguments to format the url
// Examples:
//           Delete("http://example.com")
//           Delete("http://%s/%s", domain, path)
func Delete(url string, a ...interface{}) hit.IStep {
	return hit.Delete(url, a...)
}

// Get creates a new Hit instance with GET as the http method, use the optional arguments to format the url
// Examples:
//           Get("http://example.com")
//           Get("http://%s/%s", domain, path)
func Get(url string, a ...interface{}) hit.IStep {
	return hit.Get(url, a...)
}

// Head creates a new Hit instance with HEAD as the http method, use the optional arguments to format the url
// Examples:
//           Head("http://example.com")
//           Head("http://%s/%s", domain, path)
func Head(url string, a ...interface{}) hit.IStep {
	return hit.Head(url, a...)
}

// Post creates a new Hit instance with POST as the http method, use the optional arguments to format the url
// Examples:
//           Post("http://example.com")
//           Post("http://%s/%s", domain, path)
func Post(url string, a ...interface{}) hit.IStep {
	return hit.Post(url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http method, use the optional arguments to format the url
// Examples:
//           Options("http://example.com")
//           Options("http://%s/%s", domain, path)
func Options(url string, a ...interface{}) hit.IStep {
	return hit.Options(url, a...)
}

// Put creates a new Hit instance with PUT as the http method, use the optional arguments to format the url
// Examples:
//           Put("http://example.com")
//           Put("http://%s/%s", domain, path)
func Put(url string, a ...interface{}) hit.IStep {
	return hit.Put(url, a...)
}

// Trace creates a new Hit instance with TRACE as the http method, use the optional arguments to format the url
// Examples:
//           Trace("http://example.com")
//           Trace("http://%s/%s", domain, path)
func Trace(url string, a ...interface{}) hit.IStep {
	return hit.Trace(url, a...)
}

// Test runs the specified steps and calls t.Error() if any error occurs during execution
func Test(t hit.TestingT, steps ...hit.IStep) {
	hit.Test(t, steps...)
}

// Do runs the specified steps and returns error if something was wrong
func Do(steps ...hit.IStep) error {
	return hit.Do(steps...)
}

// CombineSteps combines multiple steps to one
func CombineSteps(steps ...hit.IStep) hit.IStep {
	return hit.CombineSteps(steps...)
}
