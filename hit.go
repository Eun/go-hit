// Package hit provides an http integration test framework.
//
// It is designed to be flexible as possible, but to keep a simple to use interface for developers.
// Example:
//
// package main
//
// import (
//     "net/http"
//     . "github.com/Eun/go-hit"
// )
//
// func main() {
//     MustDo(
//         Description("Post to httpbin.org"),
//         Get("https://httpbin.org/post"),
//         Expect().Status().Equal(http.StatusMethodNotAllowed),
//         Expect().Body().String().Contains("Method Not Allowed"),
//     )
// }
//
// Or use the `Test()` function:
// package main_test
// import (
//     "testing"
//     "net/http"
//     . "github.com/Eun/go-hit"
// )
//
// func TestHttpBin(t *testing.T) {
//     Test(t,
//         Description("Post to httpbin.org"),
//         Get("https://httpbin.org/post"),
//         Expect().Status().Equal(http.StatusMethodNotAllowed),
//         Expect().Body().String().Contains("Method Not Allowed"),
//     )
// }
package hit

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"
)

// Callback will be used for Custom() functions.
type Callback func(hit Hit) error

// Hit is the interface that will be passed in for Custom() steps.
type Hit interface {
	// Request returns the current request.
	Request() *HTTPRequest

	// SetRequest sets the request for the current instance.
	SetRequest(request *http.Request) error

	// Response returns the current Response.
	Response() *HTTPResponse

	// HTTPClient gets the current http.Client.
	HTTPClient() *http.Client

	// SetHTTPClient sets the http client for the request.
	SetHTTPClient(client *http.Client) error

	// BaseURL returns the current base url.
	BaseURL() string

	// CurrentStep returns the current working step.
	CurrentStep() IStep

	// Steps returns the current step list.
	Steps() []IStep

	// AddSteps appends the specified steps to the step list.
	AddSteps(steps ...IStep)

	// InsertSteps inserts the specified steps right after the current step.
	InsertSteps(steps ...IStep)

	// RemoveSteps removes the specified steps from the step list.
	RemoveSteps(steps ...IStep)

	// Do runs the specified steps in in this context.
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Custom(func(hit Hit) error {
	//             return hit.Do(
	//                 Expect().Status().Equal(http.StatusOK),
	//             )
	//         }),
	//     )
	Do(steps ...IStep) error

	// MustDo runs the specified steps in in this context and panics on failure.
	//
	// Example:
	//    MustDo(
	//         Get("https://example.com"),
	//         Expect().Custom(func(hit Hit) error {
	//             hit.MustDo(
	//                 Expect().Status().Equal(http.StatusOK),
	//             )
	//             return nil
	//         }),
	//    )
	MustDo(steps ...IStep)

	// Description gets the current description that will be printed in an error case.
	Description() string

	// SetDescription sets a custom description for this test.
	// The description will be printed in an error case
	SetDescription(string)

	// Context returns the current request context
	Context() context.Context
}

type hitImpl struct {
	steps       []IStep
	currentStep IStep
	request     *HTTPRequest
	response    *HTTPResponse
	client      *http.Client
	state       StepTime
	baseURL     string
	description string
}

func (hit *hitImpl) Request() *HTTPRequest {
	return hit.request
}

func (hit *hitImpl) SetRequest(request *http.Request) error {
	if request == nil {
		return xerrors.New("request cannot be nil")
	}
	hit.request = newHTTPRequest(hit, request)
	return nil
}

func (hit *hitImpl) Response() *HTTPResponse {
	return hit.response
}

func (hit *hitImpl) HTTPClient() *http.Client {
	return hit.client
}

func (hit *hitImpl) SetHTTPClient(client *http.Client) error {
	if client == nil {
		return xerrors.New("client cannot be nil")
	}
	hit.client = client
	return nil
}

func (hit *hitImpl) BaseURL() string {
	return hit.baseURL
}

func (hit *hitImpl) collectSteps(state StepTime) []IStep {
	var collectedSteps []IStep
	for i := 0; i < len(hit.steps); i++ {
		if hit.steps[i] == nil {
			continue
		}
		w := hit.steps[i].when()
		if w == state {
			collectedSteps = append(collectedSteps, hit.steps[i])
		}
	}
	return collectedSteps
}

func (hit *hitImpl) runSteps(state StepTime) *Error {
	// find all steps we want to run
	stepsToRun := hit.collectSteps(state)
	size := len(stepsToRun)

	// be optimistic and hope nobody changes the size
	// if not the slice will resize
	executedSteps := make([]IStep, 0, size)

nextStep:
	for i := 0; i < size; i++ {
		for _, step := range executedSteps {
			// step already executed
			if step == stepsToRun[i] {
				continue nextStep
			}
		}

		hit.currentStep = stepsToRun[i]
		if err := execStep(hit, stepsToRun[i]); err != nil {
			return err
		}
		executedSteps = append(executedSteps, stepsToRun[i])

		// maybe the steps got modified in some way
		// lets quickly get the steps we need to execute and compare them with the ones we currently execute
		newSteps := hit.collectSteps(state)
		if !stepsAreEqual(stepsToRun, newSteps) {
			// yep something changed in our scope
			// start over again
			i = -1
			size = len(newSteps)
			stepsToRun = newSteps
			continue nextStep
		}
	}
	return nil
}

func stepsAreEqual(a, b []IStep) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (hit *hitImpl) CurrentStep() IStep {
	return hit.currentStep
}

func (hit *hitImpl) Steps() []IStep {
	return hit.steps
}

func (hit *hitImpl) AddSteps(steps ...IStep) {
	hit.steps = append(hit.steps, steps...)
}

func (hit *hitImpl) InsertSteps(steps ...IStep) {
	for i, step := range hit.steps {
		if step != hit.currentStep {
			continue
		}
		hit.steps = append(hit.steps[:i+1], append(steps, hit.steps[i+1:]...)...)
		return
	}
}

func (hit *hitImpl) RemoveSteps(steps ...IStep) {
	for j := len(steps) - 1; j >= 0; j-- {
		for i := len(hit.steps) - 1; i >= 0; i-- {
			if hit.steps[i] == steps[j] {
				// remove the step from steps and hit.steps
				steps = append(steps[:j], steps[j+1:]...)
				hit.steps = append(hit.steps[:i], hit.steps[i+1:]...)
				hit.RemoveSteps(steps...)
				return
			}
		}
	}
}

func (hit *hitImpl) Do(steps ...IStep) error {
	for _, step := range steps {
		if step.when() != hit.state && step.when() != cleanStep {
			return xerrors.Errorf(
				"unable to execute %s during %s, can only be run during %s",
				step.callPath().CallString(true),
				hit.state.String(), step.when().String(),
			)
		}
		if err := execStep(hit, step); err != nil {
			return err
		}
	}
	return nil
}

func (hit *hitImpl) MustDo(steps ...IStep) {
	if err := hit.Do(steps...); err != nil {
		panic(err)
	}
}

func (hit *hitImpl) Description() string {
	return hit.description
}

func (hit *hitImpl) SetDescription(description string) {
	hit.description = description
}

func (hit *hitImpl) Context() context.Context {
	return hit.request.Context()
}
