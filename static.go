package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Eun/go-hit/internal/misc"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/errortrace"

	urlpkg "net/url"
)

//nolint:gochecknoglobals // ett is used to initialize the errortrace, we always want to ignore some packages, so using
// a global would help to avoid unnecessary initializations.
var ett *errortrace.Template

//nolint:gochecknoinits // ett is used to initialize the errortrace, we always want to ignore some packages, so using
// a global would help to avoid unnecessary initializations.
func init() {
	ett = errortrace.New(
		"testing",
		"runtime",
		errortrace.IgnorePackage(Do),
	)
}

// Send sends the specified data as the body payload.
//
// Examples:
//     MustDo(
//         Post("https://example.com"),
//         Send().Body().String("Hello World"),
//     )
//
//     MustDo(
//         Post("https://example.com"),
//         Send().Body().Int(8),
//     )
func Send() ISend {
	return newSend(newCallPath("Send", nil))
}

// Expect provides assertions for the response data, e.g. for body or headers.
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Expect().Body().String().Equal("Hello World"),
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Expect().Body().String().Contains("Hello World"),
//     )
func Expect() IExpect {
	return newExpect(newCallPath("Expect", nil))
}

// Debug prints the current Request and Response.
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Debug(),
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Debug().Request().Headers(),
//         Debug().Response().Headers("Content-Type"),
//     )
func Debug() IDebug {
	return newDebug(newCallPath("Debug", nil), nil)
}

// Fdebug prints the current Request and Response to the specified writer.
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//         Fdebug(os.Stderr),
//     )
//
//     MustDo(
//         Get("https://example.com"),
//         Debug().Request().Headers(),
//         Debug().Response().Headers("Content-Type"),
//     )
func Fdebug(w io.Writer) IDebug {
	return newDebug(newCallPath("Fdebug", nil), w)
}

// Store stores the current Request or Response.
//
// Examples:
//     var body string
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Body().String().In(&body),
//     )
//
//     var headers http.Header
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Headers().In(&headers),
//     )
//     var contentType string
//     MustDo(
//         Get("https://example.com"),
//         Store().Response().Headers("Content-Type").In(&contentType),
//     )
func Store() IStore {
	return newStore()
}

// HTTPClient sets the client for the request.
//
// Example:
//     client := http.DefaultClient
//     MustDo(
//         Get("https://example.com"),
//         HTTPClient(client),
//     )
func HTTPClient(client *http.Client) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: newCallPath("HTTPClient", nil),
		Exec: func(hit *hitImpl) error {
			return hit.SetHTTPClient(client)
		},
	}
}

// BaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace or Method.
//
// Examples:
//     MustDo(
//         BaseURL("https://example.com/"),
//         Get("index.html"),
//     )
//
//     MustDo(
//         BaseURL("https://%s.%s/", "example", "com"),
//         Get("index.html"),
//     )
func BaseURL(url string, a ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: newCallPath("BaseURL", nil),
		Exec: func(hit *hitImpl) error {
			if len(a) > 0 {
				hit.baseURL = fmt.Sprintf(url, a...)
				return nil
			}
			hit.baseURL = url
			return nil
		},
	}
}

// Request provides methods to set request parameters.
//
// Example:
//     request, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
//     MustDo(
//         Request().Set(request),
//     )
func Request() IRequest {
	return newRequest(newCallPath("Request", nil))
}

// Method sets the specified method and url.
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
	return makeMethodStep("Method", method, url, a...)
}

func makeMethodStep(fnName, method, url string, a ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: newCallPath(fnName, nil),
		Exec: func(hit *hitImpl) error {
			hit.request.Method = method
			u := misc.MakeURL(hit.baseURL, url, a...)
			if u == "" {
				hit.request.URL = new(urlpkg.URL)
				return nil
			}
			var err error
			hit.request.URL, err = urlpkg.Parse(u)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

// Connect sets the the method to CONNECT and the specified url.
func Connect(url string, a ...interface{}) IStep {
	return makeMethodStep("Connect", http.MethodConnect, url, a...)
}

// Delete sets the the method to DELETE and the specified url.
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
	return makeMethodStep("Delete", http.MethodDelete, url, a...)
}

// Get sets the the method to GET and the specified url.
//
// Examples:
//     MustDo(
//         Get("https://example.com"),
//     )
//
//     MustDo(
//         Get("https://%s/%s", "example.com", "index.html"),
//         Expect().Status().Equal(http.StatusOK),
//         Expect().Body().String().Equal("Hello World"),
//     )
//
func Get(url string, a ...interface{}) IStep {
	return makeMethodStep("Get", http.MethodGet, url, a...)
}

// Head sets the the method to HEAD and the specified url.
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
	return makeMethodStep("Head", http.MethodHead, url, a...)
}

// Post sets the the method to POST and the specified url.
//
// Examples:
//     MustDo(
//         Post("https://example.com"),
//     )
//
//     MustDo(
//         Post("https://%s.%s", "example", "com"),
//         Expect().Status().Equal(http.StatusOK),
//         Send().Body().String("Hello World"),
//     )
//
func Post(url string, a ...interface{}) IStep {
	return makeMethodStep("Post", http.MethodPost, url, a...)
}

// Options sets the the method to OPTIONS and the specified url.
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
	return makeMethodStep("Options", http.MethodOptions, url, a...)
}

// Put sets the the method to PUT and the specified url.
//
// Examples:
//     MustDo(
//         Put("https://example.com"),
//     )
//
//     MustDo(
//         Put("https://%s,%s", "example", "com"),
//     )
//
func Put(url string, a ...interface{}) IStep {
	return makeMethodStep("Put", http.MethodPut, url, a...)
}

// Trace sets the the method to TRACE and the specified url.
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
	return makeMethodStep("Trace", http.MethodTrace, url, a...)
}

// Test runs the specified steps and calls t.FailNow() if any error occurs during execution.
func Test(t TestingT, steps ...IStep) {
	if err := Do(steps...); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		t.FailNow()
	}
}

// do func that ensures we always return an *ErrorTrace.
func do(steps ...IStep) *Error {
	hit := &hitImpl{
		client: http.DefaultClient,
		steps:  steps,
		state:  combineStep,
	}
	if err := hit.runSteps(combineStep); err != nil {
		return err
	}
	hit.state = cleanStep
	if err := hit.runSteps(cleanStep); err != nil {
		return err
	}

	hit.request = newHTTPRequest(hit, nil)
	hit.request.Header = map[string][]string{
		// remove some standard headers
		"User-Agent": {""},
	}
	hit.state = requestCreateStep
	if err := hit.runSteps(requestCreateStep); err != nil {
		return err
	}

	// user did not specify a request with Request().Set()
	if hit.request.Method == "" || hit.request.URL == nil || hit.request.URL.Host == "" {
		return wrapError(hit, xerrors.New("unable to create a request: did you called Post(), Get(), ...?"))
	}

	if hit.request.URL.Scheme == "" {
		hit.request.URL.Scheme = "https"
	}

	hit.state = BeforeSendStep
	if err := hit.runSteps(BeforeSendStep); err != nil {
		return err
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
		return wrapError(hit, xerrors.Errorf("unable to perform request: %w", err))
	}
	hit.request.Request.ContentLength, err = hit.request.Body().Length()
	if err != nil {
		return wrapError(hit, xerrors.Errorf("unable to get body length: %w", err))
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
			return wrapError(hit, err)
		}
	}
	if hit.response.Response.Body != nil {
		if err := hit.response.Response.Body.Close(); err != nil {
			return wrapError(hit, err)
		}
	}
	return nil
}

// Do runs the specified steps and returns error if something was wrong.
func Do(steps ...IStep) error {
	if err := do(steps...); err != nil {
		return err
	}
	return nil
}

// MustDo runs the specified steps and panics with the error if something was wrong.
func MustDo(steps ...IStep) {
	if err := Do(steps...); err != nil {
		panic(err)
	}
}

// CombineSteps combines multiple steps to one.
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         CombineSteps(
//            Expect().Status().Equal(http.StatusOK),
//            Expect().Body().String().Equal("Hello World"),
//         ),
//     )
func CombineSteps(steps ...IStep) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     combineStep,
		CallPath: newCallPath("CombineSteps", nil),
		Exec: func(hit *hitImpl) error {
			hit.InsertSteps(steps...)
			return nil
		},
	}
}

// Description sets a custom description for this test.
// The description will be printed in an error case.
//
// Example:
//     MustDo(
//         Description("Check if example.com is available"),
//         Get("https://example.com"),
//     )
func Description(description string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: newCallPath("Description", nil),
		Exec: func(hit *hitImpl) error {
			hit.SetDescription(description)
			return nil
		},
	}
}

// Clear can be used to remove previous steps.
//
// Usage:
//     Clear().Send()                                        // will remove all steps chained to Send()
//                                                           // e.g Send().Body().String("Hello World")
//     Clear().Send().Body()                                 // will remove all steps chained to Send().Body()
//                                                           // e.g Send().Body().String("Hello World")
//     Clear().Send().Body().String("Hello World")           // will remove all Send().Body().String("Hello World")
//                                                           // steps
//     Clear().Expect().Body()                               // will remove all steps chained to Expect().Body()
//                                                           // e.g. Expect().Body().Equal("Hello World")
//     Clear().Expect().Body().String().Equal("Hello World") // will remove all
//                                                           // Expect().Body().String().Equal("Hello World") steps
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         Expect().Status().Equal(http.StatusNotFound),
//         Expect().Body().String().Contains("Not found!"),
//         Clear().Expect(),
//         Expect().Status().Equal(http.StatusOK),
//         Expect().Body().String().Contains("Hello World"),
//     )
func Clear() IClear {
	return newClear(newCallPath("Clear", nil))
}

// Custom can be used to run custom logic during various steps.
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         Custom(ExpectStep, func(hit Hit) error {
//             if hit.Response().Body().MustString() != "Hello World" {
//                 return errors.New("Expected Hello World")
//             }
//             return nil
//         }),
//     )
func Custom(when StepTime, exec Callback) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     when,
		CallPath: newCallPath("Custom", []interface{}{when, exec}),
		Exec: func(hit *hitImpl) error {
			return exec(hit)
		},
	}
}

// Return stops the current execution of hit, resulting in ignoring all future steps.
//
// Example:
//     MustDo(
//         Get("https://example.com"),
//         Return(),
//         Expect().Body().String().Equal("Hello World"), // will never be executed
//     )
func Return() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     cleanStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			currentStep := hit.CurrentStep()
			removeAllFollowing := false
			var stepsToRemove []IStep
			for _, step := range hit.Steps() {
				if step == currentStep {
					removeAllFollowing = true
					continue
				}
				if removeAllFollowing {
					stepsToRemove = append(stepsToRemove, step)
				}
			}
			hit.RemoveSteps(stepsToRemove...)
			return nil
		},
	}
}

// Context sets the context for the request.
//
// Example:
//     ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//     defer cancel()
//     MustDo(
//         Context(ctx),
//         Get("https://example.com"),
//         Return(),
//         Expect().Body().String().Equal("Hello World"), // will never be executed
//     )
func Context(ctx context.Context) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			hit.request.Request = hit.request.Request.WithContext(ctx)
			return nil
		},
	}
}

// JoinURL joins the specified parts to one url.
//
// Example:
//     JoinURL("https://example.com", "folder", "file") // will return "https://example.com/folder/file"
//     JoinURL("https://example.com", "index.html")     // will return "https://example.com/index.html"
//     JoinURL("https://", "example.com", "index.html") // will return "https://example.com/index.html"
//     JoinURL("example.com", "index.html") // will return "example.com/index.html"
//     MustDo(
//         Get(JoinURL("https://example.com", "index.html")),
//     )
func JoinURL(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	notEmptyParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		notEmptyParts = append(notEmptyParts, part)
	}
	if len(notEmptyParts) == 0 {
		return ""
	}

	u, err := urlpkg.Parse(strings.Join(notEmptyParts, "/"))
	if err != nil {
		return ""
	}

	// replace all "double slashes"
	for {
		old := u.Path
		u.Path = strings.ReplaceAll(u.Path, "//", "/")
		if old == u.Path {
			break
		}
	}

	return strings.TrimRight(strings.ReplaceAll(u.String(), ":///", "://"), "/")
}
