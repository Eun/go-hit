// +build -ignore

// You can use this file as an template to build your own framework. Just change / add the functions you need.
// See also examples/extensibility

package main

import (
	"io"
	"net/http"

	"github.com/Eun/go-hit"
)

// Send sends the specified data as the body payload
//
// Examples:
//     MustDo(
//         Post("https://example.com"),
//         Send().Body("Hello World"),
//     )
//
//     MustDo(
//         Post("https://example.com"),
//         Send().Body().Interface("Hello World")
//     )
func Send() hit.ISend {
	return hit.Send()
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Expect().Body("Hello World")
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Expect().Body().Contains("Hello World")
//     )
//
func Expect() hit.IExpect {
	return hit.Expect()
}

// Debug prints the current Request and Response to hit.Stdout(), omit the parameter to get more options
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Debug(),
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Debug().Request().Header(),
//         Debug().Response().Header("Content-Type"),
//     )
func Debug(expression ...string) hit.IDebug {
	return hit.Debug(expression...)
}

// Store stores the current Request or Response
//
// Examples:
//     var body string
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Body().In(&body),
//     )
//
//     var headers http.Header
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Header.In(&headers),
//     )
//     var contentType string
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Header("Content-Type").In(&contentType),
//     )
func Store() hit.IStore {
	return hit.Store()
}

// HTTPClient sets the client for the request
//
// Example:
//     var client http.Client
//     MustDo(
//         Get("https://example.com"),
//         HTTPClient(&client),
//     )
func HTTPClient(client *http.Client) hit.IStep {
	return hit.HTTPClient(client)
}

// Stdout sets the output to the specified writer
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         Stdout(os.Stderr),
//         Debug(),
//     )
func Stdout(w io.Writer) hit.IStep {
	return hit.Stdout(w)
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
func BaseURL(url string, a ...interface{}) hit.IStep {
	return hit.BaseURL(url, a...)
}

// Request creates a new Hit instance with an existing http request
//
// Example:
//     request, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
//     MustDo(
//         Request(request),
//     )
func Request(request *http.Request) hit.IStep {
	return hit.Request(request)
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
func Method(method, url string, a ...interface{}) hit.IStep {
	return hit.Method(method, url, a...)
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
func Connect(url string, a ...interface{}) hit.IStep {
	return hit.Connect(url, a...)
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
func Delete(url string, a ...interface{}) hit.IStep {
	return hit.Delete(url, a...)
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
func Get(url string, a ...interface{}) hit.IStep {
	return hit.Get(url, a...)
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
func Head(url string, a ...interface{}) hit.IStep {
	return hit.Head(url, a...)
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
func Post(url string, a ...interface{}) hit.IStep {
	return hit.Post(url, a...)
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
func Options(url string, a ...interface{}) hit.IStep {
	return hit.Options(url, a...)
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
func Put(url string, a ...interface{}) hit.IStep {
	return hit.Put(url, a...)
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

// MustDo runs the specified steps and panics with the error if something was wrong
func MustDo(steps ...hit.IStep) {
	hit.MustDo(steps...)
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
func CombineSteps(steps ...hit.IStep) hit.IStep {
	return hit.CombineSteps(steps...)
}

// Description sets a custom description for this test.
// The description will be printed in an error case
//
// Example:
//     MustDo(
//         Description("Check if example.com is available"),
//         Get("https://example.com"),
//     )
func Description(description string) hit.IStep {
	return hit.Description(description)
}

// Clear can be used to remove previous steps.
//
// Usage:
//     Clear().Send("Hello World")          // will remove all Send("Hello World") steps
//     Clear().Send().Body("Hello World")   // will remove all Send().Body("Hello World") steps
//     Clear().Expect().Body()              // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
//     Clear().Expect().Body("Hello World") // will remove all Expect().Body("Hello World") steps
//
// Example:
//     MustDo(
//         Post("https://example.com"),
//         Expect().Status(http.StatusOK),
//         Expect().Body().Contains("Welcome to example.com"),
//         Clear().Expect(),
//         Expect().Status(http.NotFound),
//         Expect().Body().Contains("Not found!"),
//     )
func Clear() hit.IClear {
	return hit.Clear()
}

// Custom can be used to run custom logic during various steps.
//
// Example:
//     MustDo(
//         Post("https://example.com"),
//         Custom(ExpectStep, func(hit Hit) {
//             if hit.Response().Body().String() != "Hello Earth" {
//                 panic("Expected Hello Earth")
//             }
//         }),
//     )
func Custom(when hit.StepTime, exec hit.Callback) hit.IStep {
	return hit.Custom(when, exec)
}
